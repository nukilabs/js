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

func TestUnmarshal_NumberToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		wantErr  bool
	}{
		{"zero to false", `{"flag": 0}`, false, false},
		{"one to true", `{"flag": 1}`, true, false},
		{"direct zero", `0`, false, false},
		{"direct one", `1`, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasPrefix(tt.input, "{") {
				// Test with map
				var result struct {
					Flag bool `js:"flag"`
				}
				err := js.Unmarshal([]byte(tt.input), &result)
				if (err != nil) != tt.wantErr {
					t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !tt.wantErr && result.Flag != tt.expected {
					t.Errorf("expected Flag=%v, got %v", tt.expected, result.Flag)
				}
			} else {
				// Test direct bool
				var b bool
				err := js.Unmarshal([]byte(tt.input), &b)
				if (err != nil) != tt.wantErr {
					t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !tt.wantErr && b != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, b)
				}
			}
		})
	}

	// Test that other numbers should fail
	t.Run("invalid number to bool", func(t *testing.T) {
		var b bool
		err := js.Unmarshal([]byte(`2`), &b)
		if err == nil {
			t.Error("expected error when unmarshaling 2 into bool, got nil")
		}
	})

	// Test with array of bools
	t.Run("array of numbers to bools", func(t *testing.T) {
		var result struct {
			Flags []bool `js:"flags"`
		}
		err := js.Unmarshal([]byte(`{"flags": [0, 1, 1, 0]}`), &result)
		if err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		expected := []bool{false, true, true, false}
		if !reflect.DeepEqual(result.Flags, expected) {
			t.Errorf("expected %v, got %v", expected, result.Flags)
		}
	})

	// Test with interface{}
	t.Run("number to interface", func(t *testing.T) {
		var result struct {
			Value interface{} `js:"value"`
		}
		// When unmarshaling to interface{}, numbers stay as numbers
		err := js.Unmarshal([]byte(`{"value": 1}`), &result)
		if err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		// Should be float64(1), not true
		if _, ok := result.Value.(float64); !ok {
			t.Errorf("expected float64 for interface{}, got %T", result.Value)
		}
	})

	// Test with mixed object
	t.Run("mixed object with number-to-bool", func(t *testing.T) {
		var result struct {
			Active bool   `js:"active"`
			Status bool   `js:"status"`
			Name   string `js:"name"`
			Count  int    `js:"count"`
		}
		err := js.Unmarshal([]byte(`{active: 1, status: 0, name: "test", count: 42}`), &result)
		if err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if result.Active != true {
			t.Errorf("expected Active=true, got %v", result.Active)
		}
		if result.Status != false {
			t.Errorf("expected Status=false, got %v", result.Status)
		}
		if result.Name != "test" {
			t.Errorf("expected Name='test', got %v", result.Name)
		}
		if result.Count != 42 {
			t.Errorf("expected Count=42, got %v", result.Count)
		}
	})
}
