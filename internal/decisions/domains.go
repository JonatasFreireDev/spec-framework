package decisions

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultDomainPaths are the product-relative roots for decision records.
// The index may override or extend these with a decisionDomains object.
var DefaultDomainPaths = map[string]string{
	"product": "knowledge/decisions/", "cross-cutting": "knowledge/decisions/",
	"design": "design/decisions/", "engineering": "engineering/decisions/",
}

func DomainPaths(index map[string]any) map[string]string {
	out := map[string]string{}
	for domain, path := range DefaultDomainPaths {
		out[domain] = normalizeRoot(path)
	}
	if raw, ok := index["decisionDomains"].(map[string]any); ok {
		for domain, value := range raw {
			if path, ok := value.(string); ok && strings.TrimSpace(domain) != "" && strings.TrimSpace(path) != "" {
				out[strings.TrimSpace(domain)] = normalizeRoot(path)
			}
		}
	}
	return out
}

func DomainForPath(path string, domainPaths map[string]string) string {
	path = filepath.ToSlash(strings.TrimSpace(path))
	for domain, root := range domainPaths {
		if strings.HasPrefix(path, root) {
			return domain
		}
	}
	return ""
}

func ValidatePath(domain, path string, domainPaths map[string]string) error {
	root, ok := domainPaths[strings.TrimSpace(domain)]
	path = filepath.ToSlash(strings.TrimSpace(path))
	if !ok {
		return fmt.Errorf("unsupported decision domain %q", domain)
	}
	if !strings.HasPrefix(path, root) {
		return fmt.Errorf("decision domain %q requires path under %s", domain, root)
	}
	return nil
}

func normalizeRoot(path string) string {
	path = strings.TrimPrefix(filepath.ToSlash(strings.TrimSpace(path)), "./")
	return strings.TrimSuffix(path, "/") + "/"
}
