package eval

import (
	"testing"
)

func TestFlattenContext(t *testing.T) {
	ctx := [][]interface{}{
		{"Title 1", []interface{}{"Sentence 1. ", "Sentence 2."}},
		{"Title 2", []interface{}{"Sentence 3."}},
	}

	expected := "Document [Title 1]: Sentence 1. Sentence 2.\nDocument [Title 2]: Sentence 3.\n"
	actual := FlattenContext(ctx)

	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestFlattenContext_Empty(t *testing.T) {
	ctx := [][]interface{}{}
	expected := ""
	actual := FlattenContext(ctx)

	if actual != expected {
		t.Errorf("Expected empty string, got: %s", actual)
	}
}
