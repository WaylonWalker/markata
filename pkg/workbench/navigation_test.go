package workbench

import "testing"

func TestWorkspaceRelativePath(t *testing.T) {
	tests := []struct {
		name, root, path, want string
	}{
		{"absolute under root", "/site", "/site/pages/blog/a.md", "pages/blog/a.md"},
		{"already relative", "/site", "pages/blog/a.md", "pages/blog/a.md"},
		{"outside root keeps path", "/site", "/elsewhere/a.md", "/elsewhere/a.md"},
		{"no root", "", "pages/a.md", "pages/a.md"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := workspaceRelativePath(tt.root, tt.path); got != tt.want {
				t.Errorf("workspaceRelativePath(%q, %q) = %q, want %q", tt.root, tt.path, got, tt.want)
			}
		})
	}
}
