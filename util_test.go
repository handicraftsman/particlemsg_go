package particlemsg

import "testing"

func TestIsSupersetOf(t *testing.T) {
	ptrn := map[string]interface{}{
		"Foo": "Bar",
		"Baz": "Quux",
		"AMap": map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
		},
	}
	what := map[string]interface{}{
		"Foo":  "Bar",
		"Baz":  "Quux",
		"Asdf": "kek",
		"AMap": map[string]interface{}{
			"key1": "val1",
			"key2": "val2",
			"key3": "val3",
		},
	}
	if !isSupersetOf(what, ptrn) {
		t.Fail()
	}
}
