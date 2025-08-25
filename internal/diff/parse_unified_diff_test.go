package diff

import (
	"testing"
)

func TestParseUnifiedDiffAndGetAddedLines(t *testing.T) {
	diffStr := "--- a/fakefile.txt\n+++ b/fakefile.txt\n@@ -1,2 +1,3 @@\n linha1\n+nova=secret\n linha2\n"
	file := "fakefile.txt"
	fd := ParseUnifiedDiff(diffStr, file)
	added := fd.GetAddedLines()
	if len(added) != 1 {
		t.Fatalf("expected 1 added line, got %d", len(added))
	}
	if added[0].Content != "nova=secret" {
		t.Errorf("expected content 'nova=secret', got %q", added[0].Content)
	}
	if added[0].Number != 2 {
		t.Errorf("expected line number 2, got %d", added[0].Number)
	}
}
