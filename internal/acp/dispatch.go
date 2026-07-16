package acp

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/workflow"
)

// Request describes one explicitly enabled ACP-compatible invocation. The
// caller must have already enforced task readiness, lease, and write scope.
type Request struct {
	Enabled     bool
	ProductRoot string
	GraphPath   string
	TaskID      string
	TaskPath    string
	WorkDir     string
	Command     string
	Args        []string
}

type Result struct {
	Output string
	Exit   int
}

func Dispatch(request Request) (Result, error) {
	if !request.Enabled {
		return Result{}, errors.New("ACP runtime dispatch is disabled")
	}
	if strings.TrimSpace(request.ProductRoot) == "" || strings.TrimSpace(request.GraphPath) == "" || strings.TrimSpace(request.TaskID) == "" || strings.TrimSpace(request.TaskPath) == "" || strings.TrimSpace(request.WorkDir) == "" || strings.TrimSpace(request.Command) == "" {
		return Result{}, errors.New("ACP dispatch requires product root, graph path, task id, task path, workdir, and command")
	}
	readiness, err := workflow.CheckTaskReadiness(request.ProductRoot, request.GraphPath, request.TaskID)
	if err != nil {
		return Result{}, err
	}
	if !readiness.Ready {
		return Result{}, errors.New("ACP dispatch requires a ready approved task with an available lease")
	}
	prompt, err := os.ReadFile(request.TaskPath)
	if err != nil {
		return Result{}, err
	}
	command := exec.Command(request.Command, request.Args...)
	command.Dir = request.WorkDir
	command.Stdin = strings.NewReader(string(prompt))
	output, err := command.CombinedOutput()
	result := Result{Output: string(output)}
	if exit, ok := err.(*exec.ExitError); ok {
		result.Exit = exit.ExitCode()
	}
	return result, err
}
