package ignore

import (
	"path/filepath"
	"strings"
)

type Matcher struct {
	patterns []string
	prefixes []string
}

func New(extra []string) *Matcher {
	defaults := []string{".git", "node_modules", "dist", "build", ".DS_Store", "Thumbs.db"}
	all := append(defaults, extra...)
	m := &Matcher{}
	for _, p := range all {
		if strings.ContainsAny(p, "*?[") {
			m.patterns = append(m.patterns, p)
		} else {
			m.prefixes = append(m.prefixes, p)
		}
	}
	return m
}

func (m *Matcher) Ignore(relPath string) bool {
	relPath = filepath.ToSlash(relPath)
	parts := strings.Split(relPath, "/")
	for _, part := range parts {
		for _, p := range m.prefixes {
			if part == p || strings.HasPrefix(part, p+"-") {
				return true
			}
		}
	}
	for _, g := range m.patterns {
		ok, _ := filepath.Match(g, filepath.Base(relPath))
		if ok {
			return true
		}
	}
	return false
}
