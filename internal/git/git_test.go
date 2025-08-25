package git

import (
	"testing"
)

func TestParseHunkStartLine(t *testing.T) {
	line := "@@ -1,2 +10,5 @@"
	got := parseHunkStartLine(line)
	want := 10
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}

	line2 := "@@ -1,2 +3 @@"
	got2 := parseHunkStartLine(line2)
	want2 := 3
	if got2 != want2 {
		t.Errorf("got %d, want %d", got2, want2)
	}
}

func TestNormalizeContent(t *testing.T) {
	in := "valor \\ No newline at end of file"
	want := "valor"
	got := normalizeContent(in)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
