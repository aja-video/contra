package utils

import (
	"github.com/google/goexpect"
	"reflect"
	"testing"
)

// An equal number of send/receive is easier to program. Test this first.
// We don't need to check the error return value since if there is an error, the batcher will be wrong.
func TestSimpleBatcherEqualSndRcv(t *testing.T) {
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
func TestSimpleBatcherNormalSndRcv(t *testing.T) {
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
func TestSimpleBatcherBadInput(t *testing.T) {
	_, err := SimpleBatcher([][]string{
		{"a", "b", "wrong"},
		{"c", "d"},
		{""},
	})

	if err == nil {
		t.Fatalf("Expected an error, but did not receive an error.")
	}
}

// I can't figure out a way to make DeepEqual compare the two Batcher slices equally.
// In the below test, sameA and sameB are exactly the same, yet they don't equal because
// of something to do with the pointers being used. Commenting out for now.
// Perhaps we could loop through the batcher and check all the cases.
//func TestVariableBatcherBasic(t *testing.T) {
//	//batcher, _ := VariableBatcher([][]string{
//	//	{"option:", "8\n"},
//	//	{".*root", "cat /conf/config.xml\n"},
//	//	{"$", "cat /conf/config.xml\n"},
//	//	{"</pfsense>"},
//	//})
//
//	sameA := []expect.Batcher{
//		&expect.BCas{C: []expect.Caser{
//			&expect.BCase{R: `option:`, S: "8\n", Rt: 1, T: expect.Continue(expect.NewStatus(7, "something failed"))},
//			&expect.BCase{R: `.*root`, S: "cat /conf/config.xml\n", Rt: 1, T: expect.Continue(expect.NewStatus(7, "something failed"))},
//			&expect.BCase{R: `$`, S: "cat /conf/config.xml\n", Rt: 1, T: expect.Continue(expect.NewStatus(7, "something failed"))},
//			&expect.BCase{R: `</pfsense>`, T: expect.OK()},
//		}},
//	}
//
//	sameB := []expect.Batcher{
//		&expect.BCas{C: []expect.Caser{
//			&expect.BCase{R: `option:`, S: "8\n", Rt: 1, T: expect.Continue(expect.NewStatus(7, "something failed"))},
//			&expect.BCase{R: `.*root`, S: "cat /conf/config.xml\n", Rt: 1, T: expect.Continue(expect.NewStatus(7, "something failed"))},
//			&expect.BCase{R: `$`, S: "cat /conf/config.xml\n", Rt: 1, T: expect.Continue(expect.NewStatus(7, "something failed"))},
//			&expect.BCase{R: `</pfsense>`, T: expect.OK()},
//		}},
//	}
//
//	if !reflect.DeepEqual(sameA, sameB) {
//		t.Fatalf("VariableBatcher did not do what we expected.\nFound: %s\nExpect:%s", sameA, sameB)
//	}
//}
