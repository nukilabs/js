package js_test

import (
	"reflect"
	"strings"
	"testing"

	js "github.com/nukilabs/js"
)

func TestUnmarshal_ObjectKeysVariants(t *testing.T) {
	cases := []string{
		`{"a":1,"b":"x"}`,
		`{a:1,b:"x"}`,
		`{'a':1,'b':"x"}`,
		`{a_b$:2, $x:3}`,
	}
	for _, tc := range cases {
		var m map[string]any
		if err := js.Unmarshal([]byte(tc), &m); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", tc, err)
		}
	}
}

func TestDecoder_Token_ObjectKeysVariants(t *testing.T) {
	cases := []string{
		`{a:1}`,
		`{'a':1}`,
		`{"a":1}`,
	}
	for _, tc := range cases {
		dec := js.NewDecoder(strings.NewReader(tc))
		// Expect Delim('{')
		tok, err := dec.Token()
		if err != nil {
			t.Fatalf("Token start err: %v", err)
		}
		if _, ok := tok.(js.Delim); !ok {
			t.Fatalf("expected Delim start, got %T", tok)
		}
		// Key
		key, err := dec.Token()
		if err != nil {
			t.Fatalf("Token key err: %v", err)
		}
		if key.(string) != "a" {
			t.Fatalf("expected key 'a', got %v", key)
		}
		// Value
		val, err := dec.Token()
		if err != nil {
			t.Fatalf("Token val err: %v", err)
		}
		if !reflect.DeepEqual(val, float64(1)) && !reflect.DeepEqual(val, js.Number("1")) {
			// default is float64(1)
			t.Fatalf("expected value 1, got %#v", val)
		}
	}
}

func TestUnmarshal_SingleQuotedStringValuesAllowed(t *testing.T) {
	var v any
	err := js.Unmarshal([]byte(`{'x':'y'}`), &v)
	if err != nil {
		t.Fatalf("unexpected error for single-quoted string value: %v", err)
	}
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %#v", v)
	}
	if got := m["x"]; got != "y" {
		t.Fatalf("expected m['x']='y', got %#v", got)
	}
}
