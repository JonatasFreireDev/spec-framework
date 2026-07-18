package projectserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

func TestStartServesLocalPageAndStopRemovesDescriptor(t *testing.T) {
	root := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	running, err := Start(ctx, Config{ProductRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.Get(running.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status=%s", response.Status)
	}
	if _, err := ReadDescriptor(root); err != nil {
		t.Fatalf("descriptor: %v", err)
	}
	if _, err := Healthy(root); err != nil {
		t.Fatalf("health: %v", err)
	}
	if err := Stop(root); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-running.Done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop")
	}
	deadline := time.Now().Add(2 * time.Second)
	for {
		_, err := os.Stat(filepath.Join(root, ".product", descriptorName))
		if os.IsNotExist(err) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("descriptor still exists: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStatusAndRejectionEndpointsUseProductData(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	artifact := workflow.Artifact{ID: "TASK-1", Type: "task", Status: "draft", Path: "tasks/one.md"}
	data, err := json.Marshal(workflow.Registry{Artifacts: []workflow.Artifact{artifact}})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "tasks", "one.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("status: draft\n# Tarefa de exemplo\n\nConteúdo."), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	running, err := Start(ctx, Config{ProductRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = Stop(root); <-running.Done }()
	response, err := http.Get(running.URL + "/api/status")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var status projectStatus
	if err := json.NewDecoder(response.Body).Decode(&status); err != nil {
		t.Fatal(err)
	}
	if len(status.Documents) != 1 || status.Documents[0].Title != "Tarefa de exemplo" || status.Metrics.Pending != 1 {
		t.Fatalf("status=%+v", status)
	}
	planRequest := bytes.NewBufferString(`{"artifactIds":["TASK-1"]}`)
	response, err = http.Post(running.URL+"/api/batch-approval-plan", "application/json", planRequest)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("batch plan status=%s", response.Status)
	}
	var plan workflow.BatchPlan
	if err := json.NewDecoder(response.Body).Decode(&plan); err != nil {
		t.Fatal(err)
	}
	if len(plan.ToApprove) != 1 || plan.ToApprove[0].ID != "TASK-1" {
		t.Fatalf("batch plan=%+v", plan)
	}
	body := bytes.NewBufferString(`{"artifactId":"TASK-1","status":"rejected","approvedBy":"Product Owner","notes":"Falta definir os critérios.","confirmed":true}`)
	response, err = http.Post(running.URL+"/api/transition", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("rejection status=%s", response.Status)
	}
	registry, err := workflow.LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	if registry.Artifacts[0].Status != "rejected" {
		t.Fatalf("artifact status=%s", registry.Artifacts[0].Status)
	}
}

func TestProjectViewIncludesTypeRelationsWorktreeAndChanges(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := workflow.Registry{Artifacts: []workflow.Artifact{
		{ID: "FEATURE-1", Type: "feature", Status: "approved", Path: "domains/payments/feature.md"},
		{ID: "GRAPH-1", Type: "execution-graph", Status: "draft", Path: "domains/payments/execution-graph.json", ParentIDs: []string{"FEATURE-1"}},
		{ID: "LANDSCAPE-1", Type: "product-landscape", Status: "approved", Path: "knowledge/assessments/product-landscape.md"},
	}}
	data, _ := json.Marshal(registry)
	if err := os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "domains", "payments"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "knowledge", "assessments"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "domains", "payments", "feature.md"), []byte("# Pagamentos\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "domains", "payments", "execution-graph.json"), []byte(`{"nodes":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "knowledge", "assessments", "product-landscape.md"), []byte("# Product Landscape\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "notes.md"), []byte("não registrado"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	running, err := Start(ctx, Config{ProductRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = Stop(root); <-running.Done }()
	response, err := http.Get(running.URL + "/api/project-view")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status=%s", response.Status)
	}
	var view projectView
	if err := json.NewDecoder(response.Body).Decode(&view); err != nil {
		t.Fatal(err)
	}
	byID := map[string]artifactView{}
	for _, artifact := range view.Artifacts {
		byID[artifact.ID] = artifact
	}
	if len(view.Artifacts) != 3 || byID["GRAPH-1"].View.Renderer != "graph" || byID["LANDSCAPE-1"].View.Renderer != "markdown" {
		t.Fatalf("artifacts=%+v", view.Artifacts)
	}
	foundLandscapeType := false
	for _, config := range view.Types {
		if config.Type == "product-landscape" && config.Renderer == "markdown" {
			foundLandscapeType = true
			break
		}
	}
	if !foundLandscapeType {
		t.Fatalf("type catalog does not include product-landscape: %+v", view.Types)
	}
	if len(byID["FEATURE-1"].Children) != 1 || byID["FEATURE-1"].Children[0] != "GRAPH-1" {
		t.Fatalf("relations=%+v", byID["FEATURE-1"])
	}
	if view.Metrics.Untracked == 0 {
		t.Fatalf("untracked metrics=%+v", view.Metrics)
	}
	changes, err := http.Get(running.URL + "/api/project-view/changes?since=0&wait=0")
	if err != nil {
		t.Fatal(err)
	}
	defer changes.Body.Close()
	var update struct {
		Changed bool `json:"changed"`
	}
	if err := json.NewDecoder(changes.Body).Decode(&update); err != nil {
		t.Fatal(err)
	}
	if !update.Changed {
		t.Fatal("expected initial revision to differ from zero")
	}
}

func TestProjectViewSeparatesMarkdownFrontmatterFromContent(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product"), 0o755); err != nil {
		t.Fatal(err)
	}
	registryData, _ := json.Marshal(workflow.Registry{Artifacts: []workflow.Artifact{{ID: "SPEC-1", Type: "specification", Status: "draft", Path: "docs/spec.md"}}})
	if err := os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), registryData, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nstatus: draft\nmaturity: mockup\ndelivery:\n  level: L1\n---\n# Especificação\n\nConteúdo visível.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "spec.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	view, err := readProjectView(root, 1)
	if err != nil {
		t.Fatal(err)
	}
	artifact := view.Artifacts[0]
	if strings.Contains(artifact.Content, "maturity:") || !strings.Contains(artifact.Content, "Conteúdo visível") {
		t.Fatalf("content=%q", artifact.Content)
	}
	if artifact.Frontmatter["maturity"] != "mockup" || artifact.Frontmatter["delivery.level"] != "L1" {
		t.Fatalf("frontmatter=%+v", artifact.Frontmatter)
	}
}

func TestDashboardProvidesTheApprovedDocumentationFlows(t *testing.T) {
	for _, expected := range []string{
		"Buscar documentação",
		"Revisão de documentação",
		"Solicitar ajustes",
		"Recusar",
		"Recusar documentos",
		"Revisão em lote",
		"Mais filtros",
		"Abrir →",
		"/api/project-view",
		"/api/batch-approve",
	} {
		if !strings.Contains(dashboardHTML, expected) {
			t.Fatalf("dashboard is missing %q", expected)
		}
	}
}

func TestBatchRejectEndpointRequestsChangesForEverySelectedArtifact(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".product", "history"), 0o755); err != nil {
		t.Fatal(err)
	}
	registry := workflow.Registry{Artifacts: []workflow.Artifact{
		{ID: "DOC-1", Type: "specification", Status: "draft", Path: "docs/one.md"},
		{ID: "DOC-2", Type: "specification", Status: "draft", Path: "docs/two.md"},
	}}
	data, _ := json.Marshal(registry)
	if err := os.WriteFile(filepath.Join(root, ".product", "artifacts.json"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"docs/one.md", "docs/two.md"} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("status: draft\n# Documento\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	running, err := Start(ctx, Config{ProductRoot: root})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = Stop(root); <-running.Done }()
	response, err := http.Post(running.URL+"/api/batch-reject", "application/json", bytes.NewBufferString(`{"artifactIds":["DOC-1","DOC-2"],"approvedBy":"Product Owner","notes":"Ajustar os critérios.","confirmed":true}`))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status=%s", response.Status)
	}
	updated, err := workflow.LoadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, artifact := range updated.Artifacts {
		if artifact.Status != "rejected" {
			t.Fatalf("artifact %s status=%s", artifact.ID, artifact.Status)
		}
	}
}
