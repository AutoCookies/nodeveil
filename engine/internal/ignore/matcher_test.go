package ignore

import "testing"

func TestMatcher(t *testing.T) {
	m := New([]string{"*.tmp"})
	cases := map[string]bool{
		".git/config":          true,
		"node_modules/a.js":    true,
		"src/file.tmp":         true,
		"src/Thumbs.db":        true,
		"src/main.go":          false,
		"nested/build/out.txt": true,
	}
	for in, want := range cases {
		if got := m.Ignore(in); got != want {
			t.Fatalf("%s => %v want %v", in, got, want)
		}
	}
}
