package js

import (
	"testing"
)

func TestUnmarshal_UnknownIdentifiers(t *testing.T) {
	var m map[string]any
	err := Unmarshal([]byte("{a: unknownValue, b: myVar}"), &m)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if m["a"] != nil {
		t.Fatalf("expected a to be nil, got %#v", m["a"])
	}
	if m["b"] != nil {
		t.Fatalf("expected b to be nil, got %#v", m["b"])
	}
}

func TestUnmarshal_MixedIdentifiers(t *testing.T) {
	var m map[string]any
	err := Unmarshal([]byte("{a: true, b: false, c: null, d: undefined, e: unknown, f: 42}"), &m)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if m["a"] != true {
		t.Fatalf("expected a (true) to be true, got %#v", m["a"])
	}
	if m["b"] != false {
		t.Fatalf("expected b (false) to be false, got %#v", m["b"])
	}
	if m["c"] != nil {
		t.Fatalf("expected c (null) to be nil, got %#v", m["c"])
	}
	if m["d"] != nil {
		t.Fatalf("expected d (undefined) to be nil, got %#v", m["d"])
	}
	if m["e"] != nil {
		t.Fatalf("expected e (unknown) to be nil, got %#v", m["e"])
	}
	if m["f"] != float64(42) {
		t.Fatalf("expected f to be 42, got %#v", m["f"])
	}
}

func TestUnmarshal_UnknownInArray(t *testing.T) {
	var a []any
	err := Unmarshal([]byte("[unknown, 1, null, true]"), &a)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if len(a) != 4 {
		t.Fatalf("expected 4 elements, got %d", len(a))
	}
	if a[0] != nil {
		t.Fatalf("expected a[0] to be nil, got %#v", a[0])
	}
	if a[1] != float64(1) {
		t.Fatalf("expected a[1] to be 1, got %#v", a[1])
	}
	if a[2] != nil {
		t.Fatalf("expected a[2] to be nil, got %#v", a[2])
	}
	if a[3] != true {
		t.Fatalf("expected a[3] to be true, got %#v", a[3])
	}
}
