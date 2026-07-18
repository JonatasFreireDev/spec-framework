package validator

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPreSpecificationBaselineBlocksRealSpecification(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{".product", "knowledge/assessments", "engineering", "design/system", "domains/live"} {
		if err := os.MkdirAll(filepath.Join(root, path), 0755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(path, data string) {
		if err := os.WriteFile(filepath.Join(root, path), []byte(data), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write(".product/framework.json", `{"baseline_policy":{"pre_specification":"required"}}`)
	write("knowledge/assessments/product-landscape.md", "# Landscape\n")
	write("engineering/engineering-system.md", "# Engineering\n")
	write("design/system/design-system.md", "# Design\n")
	write("domains/live/specification.md", "# Specification\n")
	write(".product/artifacts.json", `{"artifacts":[{"id":"SPEC-1","type":"specification","status":"draft","path":"domains/live/specification.md"},{"id":"LANDSCAPE","type":"product-landscape","status":"draft","path":"knowledge/assessments/product-landscape.md"},{"id":"ENG","type":"engineering-system","status":"draft","path":"engineering/engineering-system.md"},{"id":"DS","type":"design-system","status":"draft","path":"design/system/design-system.md"}]}`)
	result, err := Validate(context.Background(), root, ".")
	if err != nil {
		t.Fatal(err)
	}
	found := 0
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Check == "pre-specification-baseline" {
			found++
		}
	}
	if found != 3 {
		t.Fatalf("baseline diagnostics=%d, want 3: %#v", found, result.Diagnostics)
	}
}
