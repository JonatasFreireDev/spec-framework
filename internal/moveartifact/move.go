package moveartifact

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/JonatasFreireDev/spec-framework/internal/fsx"
)

var markdownLink = regexp.MustCompile(`\[([^\]\n]+)\]\(([^)\n]+)\)`)

type Rewrite struct{ Path, Kind, Content string }
type Plan struct {
	Root, From, To, OldRel, NewRel string
	Rewrites                       []Rewrite
	Mentions                       []string
}

func Build(root, from, to string) (Plan, error) {
	root, _ = filepath.Abs(root)
	fromAbs := filepath.Join(root, filepath.Clean(from))
	toAbs := filepath.Join(root, filepath.Clean(to))
	if !fsx.Inside(root, fromAbs) || !fsx.Inside(root, toAbs) {
		return Plan{}, errors.New("both --from and --to must stay inside the repository")
	}
	if _, err := os.Stat(fromAbs); err != nil {
		return Plan{}, fmt.Errorf("source does not exist: %s", from)
	}
	if _, err := os.Stat(toAbs); err == nil {
		return Plan{}, fmt.Errorf("target already exists: %s", to)
	}
	oldRel, _ := filepath.Rel(root, fromAbs)
	newRel, _ := filepath.Rel(root, toAbs)
	oldRel, newRel = filepath.ToSlash(oldRel), filepath.ToSlash(newRel)
	p := Plan{Root: root, From: fromAbs, To: toAbs, OldRel: oldRel, NewRel: newRel}
	var files []string
	err := filepath.WalkDir(root, func(item string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() && (strings.HasSuffix(item, ".md") || strings.HasSuffix(item, ".json")) {
			files = append(files, item)
		}
		return nil
	})
	if err != nil {
		return Plan{}, err
	}
	sort.Strings(files)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return Plan{}, err
		}
		text := strings.TrimPrefix(string(data), "\ufeff")
		current := file
		if file == fromAbs || strings.HasPrefix(file, fromAbs+string(filepath.Separator)) {
			current = filepath.Join(toAbs, strings.TrimPrefix(file, fromAbs))
		}
		if strings.HasSuffix(file, ".md") {
			next, changed := rewriteMarkdown(text, file, fromAbs, toAbs)
			if changed {
				p.Rewrites = append(p.Rewrites, Rewrite{current, "markdown-links", normalize(next)})
			}
			for i, line := range strings.Split(text, "\n") {
				if strings.Contains(line, oldRel) && !markdownLink.MatchString(line) {
					rel, _ := filepath.Rel(root, current)
					p.Mentions = append(p.Mentions, fmt.Sprintf("%s:%d: %s", filepath.ToSlash(rel), i+1, strings.TrimSpace(line)))
				}
			}
		} else if next, changed := rewriteJSON(text, oldRel, newRel); changed {
			p.Rewrites = append(p.Rewrites, Rewrite{current, "json-paths", next})
		}
	}
	sort.Slice(p.Rewrites, func(i, j int) bool { return p.Rewrites[i].Path < p.Rewrites[j].Path })
	sort.Strings(p.Mentions)
	return p, nil
}

func Apply(p Plan) error {
	if err := os.MkdirAll(filepath.Dir(p.To), 0o755); err != nil {
		return err
	}
	if err := os.Rename(p.From, p.To); err != nil {
		return err
	}
	originals := map[string][]byte{}
	for _, rewrite := range p.Rewrites {
		data, err := os.ReadFile(rewrite.Path)
		if err != nil {
			rollback(p, originals)
			return err
		}
		originals[rewrite.Path] = data
		if err := os.WriteFile(rewrite.Path, []byte(rewrite.Content), 0o644); err != nil {
			rollback(p, originals)
			return err
		}
	}
	return nil
}

func rollback(p Plan, originals map[string][]byte) {
	for file, data := range originals {
		_ = os.WriteFile(file, data, 0o644)
	}
	_ = os.Rename(p.To, p.From)
}

func rewriteMarkdown(text, file, oldAbs, newAbs string) (string, bool) {
	changed := false
	next := markdownLink.ReplaceAllStringFunc(text, func(full string) string {
		parts := markdownLink.FindStringSubmatch(full)
		target := strings.TrimSpace(parts[2])
		if target == "" || strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
			return full
		}
		anchor := ""
		if i := strings.Index(target, "#"); i >= 0 {
			anchor, target = target[i:], target[:i]
		}
		decoded, err := url.PathUnescape(strings.Trim(target, "<>"))
		if err != nil {
			decoded = target
		}
		resolved := filepath.Clean(filepath.Join(filepath.Dir(file), filepath.FromSlash(decoded)))
		if resolved != oldAbs && !strings.HasPrefix(resolved, oldAbs+string(filepath.Separator)) {
			return full
		}
		moved := filepath.Join(newAbs, strings.TrimPrefix(resolved, oldAbs))
		rel, _ := filepath.Rel(filepath.Dir(file), moved)
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, ".") {
			rel = "./" + rel
		}
		changed = true
		return fmt.Sprintf("[%s](%s%s)", parts[1], rel, anchor)
	})
	return next, changed
}

func rewriteJSON(text, oldRel, newRel string) (string, bool) {
	var value any
	if json.Unmarshal([]byte(text), &value) != nil {
		return text, false
	}
	changed := rewriteValue(&value, oldRel, newRel)
	if !changed {
		return text, false
	}
	data, _ := json.MarshalIndent(value, "", "  ")
	return string(data) + "\n", true
}

func rewriteValue(ptr *any, oldRel, newRel string) bool {
	switch value := (*ptr).(type) {
	case string:
		if value == oldRel || strings.HasPrefix(value, oldRel+"/") {
			*ptr = newRel + strings.TrimPrefix(value, oldRel)
			return true
		}
	case []any:
		changed := false
		for i := range value {
			changed = rewriteValue(&value[i], oldRel, newRel) || changed
		}
		return changed
	case map[string]any:
		changed := false
		for key, item := range value {
			changed = rewriteValue(&item, oldRel, newRel) || changed
			value[key] = item
		}
		return changed
	}
	return false
}

func normalize(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
}
