package diff

import (
	"reflect"
	"testing"
)

func TestExtractRemovedLinesContent(t *testing.T) {
	diffLines := []string{
		"-senha=123",
		"-api=456",
		"@@ -1,2 +1,2 @@",
		"+senha=789",
	}
	want := []string{"senha=123", "api=456"}
	got := ExtractRemovedLinesContent(diffLines)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNormalizeContent(t *testing.T) {
	in := "valor \\ No newline at end of file"
	want := "valor"
	got := NormalizeContent(in)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIsReformattedLine(t *testing.T) {
	removed := []string{"senha=123", "api=456"}
	if !IsReformattedLine("senha=123", removed) {
		t.Error("should detect reformatted line")
	}
	if IsReformattedLine("nova=789", removed) {
		t.Error("should not detect as reformatted")
	}
}

func TestIsContextLine(t *testing.T) {
	if !IsContextLine("linha normal") {
		t.Error("should be context line")
	}
	if IsContextLine("@@ -1,2 +1,2 @@") {
		t.Error("should not be context line")
	}
	if IsContextLine("+++") {
		t.Error("should not be context line")
	}
}
