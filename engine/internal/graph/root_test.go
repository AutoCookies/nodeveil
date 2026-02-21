package graph

import "testing"

func TestNewRoot(t *testing.T) {
	t.Run("rejects empty", func(t *testing.T) {
		if _, err := NewRoot("   "); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("accepts valid", func(t *testing.T) {
		root, err := NewRoot("/tmp")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if root.Path != "/tmp" {
			t.Fatalf("unexpected path %q", root.Path)
		}
	})
}
