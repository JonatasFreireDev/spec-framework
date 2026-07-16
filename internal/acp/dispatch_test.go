package acp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDispatchFailsClosedWhenDisabled(t *testing.T) {
	path := filepath.Join(t.TempDir(), "task.md")
	if err := os.WriteFile(path, []byte("# task"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Dispatch(Request{TaskPath: path, WorkDir: filepath.Dir(path), Command: "cmd"}); err == nil {
		t.Fatal("disabled dispatch accepted")
	}
}
