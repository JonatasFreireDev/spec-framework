package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JonatasFreireDev/spec-framework/internal/designsystem"
)

func TestDesignSystemCLIInspectAndImpact(t *testing.T) {
	root := t.TempDir()
	if _, err := designsystem.Init(root, "generate"); err != nil {
		t.Fatal(err)
	}
	var out, errout bytes.Buffer
	code := New("test").Run([]string{"design-system", "inspect", "--product-root", root}, &out, &errout)
	if code != 0 || !strings.Contains(out.String(), "DSYS-001") {
		t.Fatalf("inspect=%d out=%s err=%s", code, out.String(), errout.String())
	}
	out.Reset()
	errout.Reset()
	code = New("test").Run([]string{"impact", "--product-root", filepath.Clean(root), "--design-system", "DSYS-001"}, &out, &errout)
	if code != 0 || !strings.Contains(out.String(), "DSYS-001@0.1.0") {
		t.Fatalf("impact=%d out=%s err=%s", code, out.String(), errout.String())
	}
}
