package utils

import (
	"github.com/google/goexpect"
	"reflect"
	"testing"
)

// An equal number of send/receive is easier to program. Test this first.
// We don't need to check the error return value since if there is an error, the batcher will be wrong.
func TestSimpleBatcher_EqualSndRcv(t *testing.T) {
	batcher, _ := SimpleBatcher([][]string{
		{"a", "b"},
		{"c", "d"},
	})

	expected := []expect.Batcher{
		&expect.BExp{R: "a"},
		&expect.BSnd{S: "b"},
		&expect.BExp{R: "c"},
		&expect.BSnd{S: "d"},
	}

	if !reflect.DeepEqual(batcher, expected) {
		t.Fatalf("SimpleBatcher did not do what we expected.\nFound: %s\nExpect:%s", batcher, expected)
	}
}

// Having an uneven number is a little more tricky.
func TestSimpleBatcher_NormalSndRcv(t *testing.T) {
	batcher, _ := SimpleBatcher([][]string{
		{"a", "b"},
		{"c", "d"},
		{"e"}, // Normally we will have just one trailing send.
	})

	expected := []expect.Batcher{
		&expect.BExp{R: "a"},
		&expect.BSnd{S: "b"},
		&expect.BExp{R: "c"},
		&expect.BSnd{S: "d"},
		&expect.BExp{R: "e"},
	}

	if !reflect.DeepEqual(batcher, expected) {
		t.Fatalf("SimpleBatcher did not do what we expected.\nFound: %s\nExpect:%s", batcher, expected)
	}
}

// Having an uneven number is a little more tricky.
// We don't need to check the batcher, since we only need to know if we have an error.
func TestSimpleBatcher_BadInput(t *testing.T) {
	_, err := SimpleBatcher([][]string{
		{"a", "b", "wrong"},
		{"c", "d"},
		{""},
	})

	if err == nil {
		t.Fatalf("Expected an error, but did not receive an error.")
	}
}
