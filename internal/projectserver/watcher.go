package projectserver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// worktreeWatcher is intentionally polling-based. It has no platform-specific
// dependency and observes both normal files and Git metadata changes.
type worktreeWatcher struct {
	root        string
	mu          sync.RWMutex
	revision    uint64
	fingerprint string
	changed     chan struct{}
}

func newWorktreeWatcher(root string) *worktreeWatcher {
	return &worktreeWatcher{root: root, revision: 1, changed: make(chan struct{})}
}
func (w *worktreeWatcher) Revision() uint64 { w.mu.RLock(); defer w.mu.RUnlock(); return w.revision }
func (w *worktreeWatcher) Run(ctx context.Context) {
	w.update()
	ticker := time.NewTicker(750 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.update()
		}
	}
}
func (w *worktreeWatcher) Wait(ctx context.Context, since uint64) uint64 {
	w.mu.RLock()
	if w.revision != since {
		r := w.revision
		w.mu.RUnlock()
		return r
	}
	changed := w.changed
	w.mu.RUnlock()
	select {
	case <-ctx.Done():
	case <-changed:
	}
	return w.Revision()
}
func (w *worktreeWatcher) update() {
	fingerprint := worktreeFingerprint(w.root)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.fingerprint == "" {
		w.fingerprint = fingerprint
		return
	}
	if fingerprint == w.fingerprint {
		return
	}
	w.fingerprint = fingerprint
	w.revision++
	close(w.changed)
	w.changed = make(chan struct{})
}
func worktreeFingerprint(root string) string {
	items := []string{}
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if rel == ".git/objects" || rel == ".git/refs" {
				return filepath.SkipDir
			}
			return nil
		}
		info, e := d.Info()
		if e == nil {
			items = append(items, rel+"|"+info.ModTime().UTC().Format(time.RFC3339Nano)+"|"+strconv.FormatInt(info.Size(), 10))
		}
		return nil
	})
	sort.Strings(items)
	sum := sha256.Sum256([]byte(strings.Join(items, "\n")))
	return hex.EncodeToString(sum[:])
}
