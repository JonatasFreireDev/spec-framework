package dispatch

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReviewEnvelopePinsParentDiffAndCannotRun(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		t.Fatal(err)
	}
	parent := Envelope{Version: 1, ID: "DISPATCH-1", WorkspaceID: work, TaskID: "TK-1", Role: "code-runner", Agent: "runner", Status: "returned", DiffHash: "abc", InputHash: "input"}
	if err := write(filepath.Join(dir(root, work), parent.ID+".json"), parent); err != nil {
		t.Fatal(err)
	}
	review, err := AssignReview(root, work, parent.ID, "qa", "qa-1")
	if err != nil {
		t.Fatal(err)
	}
	if review.DiffHash != "abc" || review.ParentID != parent.ID || len(review.WriteScope) != 0 {
		t.Fatalf("review=%+v", review)
	}
	if _, err := Run(root, work, review.ID, false, "echo", nil); err == nil {
		t.Fatal("review run accepted")
	}
}

func TestReconcileReportsOrphanReview(t *testing.T) {
	root := t.TempDir()
	work := "WORK-001"
	if err := os.MkdirAll(dir(root, work), 0755); err != nil {
		t.Fatal(err)
	}
	if err := write(filepath.Join(dir(root, work), "DISPATCH-2.json"), Envelope{ID: "DISPATCH-2", WorkspaceID: work, Role: "code-review", ParentID: "missing", DiffHash: "abc", Status: "assigned"}); err != nil {
		t.Fatal(err)
	}
	items, err := Reconcile(root, work)
	if err != nil || len(items) != 1 || items[0].Kind != "orphaned-review" {
		t.Fatalf("items=%+v err=%v", items, err)
	}
}
