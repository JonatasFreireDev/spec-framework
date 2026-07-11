package fsx

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FS interface {
	ReadFile(string) ([]byte, error)
	WriteFile(string, []byte, fs.FileMode) error
	MkdirAll(string, fs.FileMode) error
	Rename(string, string) error
	RemoveAll(string) error
	Stat(string) (fs.FileInfo, error)
}

type OS struct{}

func (OS) ReadFile(name string) ([]byte, error) { return os.ReadFile(name) }
func (OS) WriteFile(name string, data []byte, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(name), 0o755); err != nil {
		return err
	}
	return os.WriteFile(name, data, mode)
}
func (OS) MkdirAll(path string, mode fs.FileMode) error { return os.MkdirAll(path, mode) }
func (OS) Rename(oldpath, newpath string) error         { return os.Rename(oldpath, newpath) }
func (OS) RemoveAll(path string) error                  { return os.RemoveAll(path) }
func (OS) Stat(name string) (fs.FileInfo, error)        { return os.Stat(name) }

func Inside(root, candidate string) bool {
	r, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	c, err := filepath.Abs(candidate)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(r, c)
	return err == nil && rel != ".." && !filepath.IsAbs(rel) && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
