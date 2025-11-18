package js

import "testing"

func TestMemberExpression(t *testing.T) {
	var m map[string]any
	err := Unmarshal([]byte("{a: window.location, b: obj.prop.nested}"), &m)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if m["a"] != nil {
		t.Fatalf("expected a (window.location) to be nil, got %#v", m["a"])
	}
	if m["b"] != nil {
		t.Fatalf("expected b (obj.prop.nested) to be nil, got %#v", m["b"])
	}
}

func TestMemberExpressionInArray(t *testing.T) {
	var a []any
	err := Unmarshal([]byte("[window.location, document.body, 42]"), &a)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if len(a) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(a))
	}
	if a[0] != nil {
		t.Fatalf("expected a[0] to be nil, got %#v", a[0])
	}
	if a[1] != nil {
		t.Fatalf("expected a[1] to be nil, got %#v", a[1])
	}
	if a[2] != float64(42) {
		t.Fatalf("expected a[2] to be 42, got %#v", a[2])
	}
}

func TestMemberExpressionMixed(t *testing.T) {
	var m map[string]any
	err := Unmarshal([]byte("{a: true, b: window.x, c: false, d: obj.y, e: null, f: someVar}"), &m)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if m["a"] != true {
		t.Fatalf("expected a to be true, got %#v", m["a"])
	}
	if m["b"] != nil {
		t.Fatalf("expected b (window.x) to be nil, got %#v", m["b"])
	}
	if m["c"] != false {
		t.Fatalf("expected c to be false, got %#v", m["c"])
	}
	if m["d"] != nil {
		t.Fatalf("expected d (obj.y) to be nil, got %#v", m["d"])
	}
	if m["e"] != nil {
		t.Fatalf("expected e (null) to be nil, got %#v", m["e"])
	}
	if m["f"] != nil {
		t.Fatalf("expected f (someVar) to be nil, got %#v", m["f"])
	}
}
