package projectserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

type projectStatus struct {
	Documents []document `json:"documents"`
	Metrics   metrics    `json:"metrics"`
}

type metrics struct {
	Total    int `json:"total"`
	Approved int `json:"approved"`
	Pending  int `json:"pending"`
	Rejected int `json:"rejected"`
}

type document struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Path    string `json:"path"`
	Folder  string `json:"folder"`
	Status  string `json:"status"`
	Updated string `json:"updated"`
	Content string `json:"content"`
}

type transitionRequest struct {
	ArtifactID string `json:"artifactId"`
	Status     string `json:"status"`
	ApprovedBy string `json:"approvedBy"`
	Notes      string `json:"notes"`
	Confirmed  bool   `json:"confirmed"`
}

type batchApprovalRequest struct {
	ArtifactIDs []string `json:"artifactIds"`
	ApprovedBy  string   `json:"approvedBy"`
	Notes       string   `json:"notes"`
	Confirmed   bool     `json:"confirmed"`
}

func readStatus(root string) (projectStatus, error) {
	registry, err := workflow.LoadRegistry(root)
	if err != nil {
		return projectStatus{}, fmt.Errorf("read artifact registry: %w", err)
	}
	result := projectStatus{Documents: make([]document, 0, len(registry.Artifacts))}
	for _, artifact := range registry.Artifacts {
		path := filepath.Join(root, filepath.FromSlash(artifact.Path))
		content, err := os.ReadFile(path)
		if err != nil {
			return projectStatus{}, fmt.Errorf("read %s: %w", artifact.Path, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return projectStatus{}, err
		}
		result.Documents = append(result.Documents, document{ID: artifact.ID, Title: titleFor(artifact, content), Path: filepath.ToSlash(artifact.Path), Folder: folderFor(artifact.Path), Status: artifact.Status, Updated: info.ModTime().UTC().Format(time.RFC3339), Content: string(content)})
		result.Metrics.Total++
		switch artifact.Status {
		case "approved":
			result.Metrics.Approved++
		case "rejected":
			result.Metrics.Rejected++
		default:
			result.Metrics.Pending++
		}
	}
	sort.Slice(result.Documents, func(i, j int) bool { return result.Documents[i].Path < result.Documents[j].Path })
	return result, nil
}

func findArtifact(root, id string) (workflow.Artifact, error) {
	registry, err := workflow.LoadRegistry(root)
	if err != nil {
		return workflow.Artifact{}, err
	}
	for _, artifact := range registry.Artifacts {
		if artifact.ID == id {
			return artifact, nil
		}
	}
	return workflow.Artifact{}, errors.New("artifact not found")
}

func titleFor(artifact workflow.Artifact, content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return artifact.ID
}

func folderFor(path string) string {
	folder := filepath.ToSlash(filepath.Dir(path))
	if folder == "." {
		return "Raiz"
	}
	return folder
}

func decodeJSON(request *http.Request, target any) error {
	defer request.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("invalid request body")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeAPIError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func methodNotAllowed(w http.ResponseWriter) {
	writeAPIError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}
