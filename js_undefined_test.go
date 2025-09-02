package js

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUnmarshal_UndefinedToInterface(t *testing.T) {
	var v any = 123
	if err := Unmarshal([]byte("undefined"), &v); err != nil {
		t.Fatalf("Unmarshal undefined: %v", err)
	}
	if v != nil {
		t.Fatalf("expected nil, got %#v", v)
	}
}

func TestUnmarshal_UndefinedInObject(t *testing.T) {
	var m map[string]any
	data := []byte("{a:undefined, 'b':1, \"c\":'x'}")
	if err := Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal object: %v", err)
	}
	if _, ok := m["a"]; !ok {
		t.Fatalf("missing key a")
	}
	if m["a"] != nil {
		t.Fatalf("expected a to be nil, got %#v", m["a"])
	}
	if got, want := m["b"], float64(1); !reflect.DeepEqual(got, want) {
		t.Fatalf("b: got %#v, want %#v", got, want)
	}
	if got, want := m["c"], "x"; !reflect.DeepEqual(got, want) {
		t.Fatalf("c: got %#v, want %#v", got, want)
	}
}

func TestUnmarshal_UndefinedToPointerAndPrimitives(t *testing.T) {
	type S struct {
		P *int
	}
	var s S
	if err := Unmarshal([]byte("{P:undefined}"), &s); err != nil {
		t.Fatalf("Unmarshal struct: %v", err)
	}
	if s.P != nil {
		t.Fatalf("expected P to be nil, got %v", *s.P)
	}

	// Primitives should be left unchanged when top-level undefined
	iv := 7
	if err := Unmarshal([]byte("undefined"), &iv); err != nil {
		t.Fatalf("Unmarshal undefined into int: %v", err)
	}
	if iv != 7 {
		t.Fatalf("int changed: got %d, want 7", iv)
	}

	sv := "keep"
	if err := Unmarshal([]byte("undefined"), &sv); err != nil {
		t.Fatalf("Unmarshal undefined into string: %v", err)
	}
	if sv != "keep" {
		t.Fatalf("string changed: got %q, want %q", sv, "keep")
	}
}

func TestUnmarshal_SingleQuotedStringValue(t *testing.T) {
	var s string
	if err := Unmarshal([]byte("'hello'"), &s); err != nil {
		t.Fatalf("Unmarshal single-quoted string: %v", err)
	}
	if s != "hello" {
		t.Fatalf("got %q, want %q", s, "hello")
	}
}

func TestUnmarshal_ArrayWithUndefined(t *testing.T) {
	var a []any
	if err := Unmarshal([]byte("[undefined, null, 1]"), &a); err != nil {
		t.Fatalf("Unmarshal array: %v", err)
	}
	want := []any{nil, nil, float64(1)}
	if !reflect.DeepEqual(a, want) {
		t.Fatalf("got %#v, want %#v", a, want)
	}
}

func TestDecoder_Token_WithUndefined(t *testing.T) {
	dec := NewDecoder(bytes.NewReader([]byte("{a:undefined}")))
	// {
	tok, err := dec.Token()
	if err != nil || tok.(Delim) != '{' {
		t.Fatalf("first token: %v, %v", tok, err)
	}
	// key "a"
	tok, err = dec.Token()
	if err != nil {
		t.Fatalf("key token err: %v", err)
	}
	if ks, ok := tok.(string); !ok || ks != "a" {
		t.Fatalf("key token: got %#v", tok)
	}
	// value nil (undefined treated as null)
	tok, err = dec.Token()
	if err != nil {
		t.Fatalf("value token err: %v", err)
	}
	if tok != nil {
		t.Fatalf("value token: got %#v, want nil", tok)
	}
	// }
	tok, err = dec.Token()
	if err != nil || tok.(Delim) != '}' {
		t.Fatalf("end token: %v, %v", tok, err)
	}
}
