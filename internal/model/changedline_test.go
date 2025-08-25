package model

import (
	"testing"
)

func TestChangedLineStruct(t *testing.T) {
	cl := ChangedLine{Number: 42, Content: "senha=123"}
	if cl.Number != 42 {
		t.Errorf("expected Number 42, got %d", cl.Number)
	}
	if cl.Content != "senha=123" {
		t.Errorf("expected Content 'senha=123', got %q", cl.Content)
	}
}
