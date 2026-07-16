package acp

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

// Request describes one explicitly enabled, acknowledged ACP-compatible run.
// ACP has no authority to approve, commit, push, merge, or resolve reviews.
type Request struct {
	Enabled      bool
	Acknowledged bool
	ProductRoot  string
	WorkspaceID  string
	GraphPath    string
	TaskID       string
	TaskPath     string
	WorkDir      string
	Agent        string
	Command      string
	Args         []string
}

type Result struct {
	Output     string
	Exit       int
	Transcript string
}

func Dispatch(request Request) (Result, error) {
	if !request.Enabled {
		return Result{}, errors.New("ACP runtime dispatch is disabled")
	}
	if !request.Acknowledged {
		return Result{}, errors.New("ACP dispatch requires explicit per-run acknowledgement")
	}
	if strings.TrimSpace(request.ProductRoot) == "" || strings.TrimSpace(request.WorkspaceID) == "" || strings.TrimSpace(request.GraphPath) == "" || strings.TrimSpace(request.TaskID) == "" || strings.TrimSpace(request.TaskPath) == "" || strings.TrimSpace(request.WorkDir) == "" || strings.TrimSpace(request.Agent) == "" || strings.TrimSpace(request.Command) == "" {
		return Result{}, errors.New("ACP dispatch requires product root, workspace, graph, task, workdir, agent, and command")
	}
	if strings.EqualFold(filepath.Base(request.Command), "git") {
		return Result{}, errors.New("ACP runtime cannot invoke git delivery commands")
	}
	readiness, err := workflow.CheckTaskReadiness(request.ProductRoot, request.GraphPath, request.TaskID)
	if err != nil {
		return Result{}, err
	}
	if !readiness.Ready {
		return Result{}, errors.New("ACP dispatch requires a ready approved task with an available lease")
	}
	if _, err := workflow.ClaimLease(request.ProductRoot, request.GraphPath, request.TaskID, request.Agent, 30*time.Minute); err != nil {
		return Result{}, err
	}
	defer workflow.ReleaseLease(request.ProductRoot, request.TaskID, request.Agent)
	prompt, err := os.ReadFile(request.TaskPath)
	if err != nil {
		return Result{}, err
	}
	command := exec.Command(request.Command, request.Args...)
	command.Dir = request.WorkDir
	command.Stdin = strings.NewReader(string(prompt))
	started := time.Now().UTC()
	output, runErr := command.CombinedOutput()
	result := Result{Output: string(output)}
	if exit, ok := runErr.(*exec.ExitError); ok {
		result.Exit = exit.ExitCode()
	}
	dir := filepath.Join(request.ProductRoot, ".product", "workspaces", request.WorkspaceID, "acp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return result, err
	}
	transcript := map[string]any{"version": 1, "workspace_id": request.WorkspaceID, "task_id": request.TaskID, "agent": request.Agent, "command": request.Command, "args": request.Args, "started_at": started.Format(time.RFC3339), "exit": result.Exit, "output_sha256": outputHash(output)}
	data, _ := json.MarshalIndent(transcript, "", "  ")
	result.Transcript = filepath.Join(dir, "ACP-"+request.TaskID+"-"+time.Now().UTC().Format("20060102T150405.000000000Z")+".json")
	if err := os.WriteFile(result.Transcript, append(data, '\n'), 0644); err != nil {
		return result, err
	}
	return result, runErr
}

func outputHash(value []byte) string { sum := sha256.Sum256(value); return hex.EncodeToString(sum[:]) }
