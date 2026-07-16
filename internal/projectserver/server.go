package projectserver

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

const descriptorName = "server.json"

// Config constrains the dashboard runtime to a local product directory.
type Config struct {
	ProductRoot string
	Port        int
}

// Descriptor is the minimal local control record used by server status and stop.
// The token makes the stop endpoint unavailable to unrelated local processes.
type Descriptor struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type Running struct {
	URL  string
	Done <-chan error
}

func Start(ctx context.Context, config Config) (Running, error) {
	root, err := filepath.Abs(config.ProductRoot)
	if err != nil {
		return Running{}, err
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return Running{}, fmt.Errorf("product root is not a directory: %s", root)
	}
	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprint(config.Port)))
	if err != nil {
		return Running{}, err
	}
	token, err := newToken()
	if err != nil {
		_ = listener.Close()
		return Running{}, err
	}
	url := "http://" + listener.Addr().String()
	if err := writeDescriptor(root, Descriptor{URL: url, Token: token}); err != nil {
		_ = listener.Close()
		return Running{}, err
	}

	shutdown := make(chan struct{})
	var shutdownOnce sync.Once
	requestShutdown := func() { shutdownOnce.Do(func() { close(shutdown) }) }
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, dashboardHTML)
	})
	mux.HandleFunc("/__spec-framework/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	})
	mux.HandleFunc("/__spec-framework/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("X-Spec-Server-Token") != token {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		requestShutdown()
	})
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		status, err := readStatus(root)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, status)
	})
	mux.HandleFunc("/api/transition", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var request transitionRequest
		if err := decodeJSON(r, &request); err != nil {
			writeAPIError(w, http.StatusBadRequest, err)
			return
		}
		if !request.Confirmed {
			writeAPIError(w, http.StatusBadRequest, errors.New("confirmation is required"))
			return
		}
		if strings.TrimSpace(request.ApprovedBy) == "" {
			writeAPIError(w, http.StatusBadRequest, errors.New("approver identity is required"))
			return
		}
		artifact, err := findArtifact(root, request.ArtifactID)
		if err != nil {
			writeAPIError(w, http.StatusNotFound, err)
			return
		}
		record, err := workflow.Approve(root, filepath.Join(root, filepath.FromSlash(artifact.Path)), request.Status, request.ApprovedBy, request.Notes)
		if err != nil {
			writeAPIError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeJSON(w, http.StatusOK, record)
	})
	mux.HandleFunc("/api/batch-approve", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var request batchApprovalRequest
		if err := decodeJSON(r, &request); err != nil {
			writeAPIError(w, http.StatusBadRequest, err)
			return
		}
		if !request.Confirmed {
			writeAPIError(w, http.StatusBadRequest, errors.New("confirmation is required"))
			return
		}
		if strings.TrimSpace(request.ApprovedBy) == "" {
			writeAPIError(w, http.StatusBadRequest, errors.New("approver identity is required"))
			return
		}
		plan, err := workflow.BuildBatchApprovalPlan(root, workflow.BatchScope{IDs: request.ArtifactIDs}, "approved")
		if err != nil {
			writeAPIError(w, http.StatusUnprocessableEntity, err)
			return
		}
		records, err := workflow.ApproveBatch(root, plan, request.ApprovedBy, request.Notes)
		if err != nil {
			writeAPIError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeJSON(w, http.StatusOK, records)
	})
	mux.HandleFunc("/api/batch-approval-plan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		var request batchApprovalRequest
		if err := decodeJSON(r, &request); err != nil {
			writeAPIError(w, http.StatusBadRequest, err)
			return
		}
		plan, err := workflow.BuildBatchApprovalPlan(root, workflow.BatchScope{IDs: request.ArtifactIDs}, "approved")
		if err != nil {
			writeAPIError(w, http.StatusUnprocessableEntity, err)
			return
		}
		writeJSON(w, http.StatusOK, plan)
	})

	httpServer := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	done := make(chan error, 1)
	go func() {
		err := httpServer.Serve(listener)
		if !errors.Is(err, http.ErrServerClosed) {
			done <- err
			return
		}
		done <- nil
	}()
	go func() {
		select {
		case <-ctx.Done():
		case <-shutdown:
		}
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
		_ = os.Remove(descriptorPath(root))
	}()
	return Running{URL: url, Done: done}, nil
}

func Stop(productRoot string) error {
	descriptor, err := ReadDescriptor(productRoot)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, descriptor.URL+"/__spec-framework/stop", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Spec-Server-Token", descriptor.Token)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("local server is not reachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("local server refused stop request: %s", resp.Status)
	}
	return nil
}

func Healthy(productRoot string) (Descriptor, error) {
	descriptor, err := ReadDescriptor(productRoot)
	if err != nil {
		return Descriptor{}, err
	}
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Get(descriptor.URL + "/__spec-framework/health")
	if err != nil {
		return Descriptor{}, fmt.Errorf("local server is not reachable: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Descriptor{}, fmt.Errorf("local server health check failed: %s", response.Status)
	}
	return descriptor, nil
}

func ReadDescriptor(productRoot string) (Descriptor, error) {
	data, err := os.ReadFile(descriptorPath(productRoot))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Descriptor{}, errors.New("no local project server is recorded for this product")
		}
		return Descriptor{}, err
	}
	var descriptor Descriptor
	if err := json.Unmarshal(data, &descriptor); err != nil {
		return Descriptor{}, fmt.Errorf("read local server descriptor: %w", err)
	}
	if !strings.HasPrefix(descriptor.URL, "http://127.0.0.1:") || descriptor.Token == "" {
		return Descriptor{}, errors.New("local server descriptor is invalid")
	}
	return descriptor, nil
}

func descriptorPath(root string) string { return filepath.Join(root, ".product", descriptorName) }

func writeDescriptor(root string, descriptor Descriptor) error {
	dir := filepath.Dir(descriptorPath(root))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(descriptor)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".server-*.tmp")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, descriptorPath(root))
}

func newToken() (string, error) {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

//go:embed dashboard.html
var dashboardHTML string
