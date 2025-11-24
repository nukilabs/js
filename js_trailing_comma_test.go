package js_test

import (
	"testing"

	js "github.com/nukilabs/js"
)

func TestUnmarshal_TrailingCommaInObject(t *testing.T) {
	input := `{a: 1, b: 2,}`

	var result map[string]interface{}
	err := js.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal failed with trailing comma: %v", err)
	}

	if result["a"] != float64(1) {
		t.Errorf("expected a=1, got %v", result["a"])
	}

	if result["b"] != float64(2) {
		t.Errorf("expected b=2, got %v", result["b"])
	}
}

func TestUnmarshal_TrailingCommaInArray(t *testing.T) {
	input := `[1, 2, 3,]`

	var result []interface{}
	err := js.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal failed with trailing comma in array: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 items, got %d", len(result))
	}

	if result[0] != float64(1) || result[1] != float64(2) || result[2] != float64(3) {
		t.Errorf("unexpected array values: %v", result)
	}
}

func TestUnmarshal_TrailingCommaInNestedStructures(t *testing.T) {
	input := `{
		items: [1, 2,],
		info: {name: 'test', value: 42,},
	}`

	var result map[string]interface{}
	err := js.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal failed with nested trailing commas: %v", err)
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		t.Fatalf("expected items to be array, got %T", result["items"])
	}

	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}

	info, ok := result["info"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected info to be object, got %T", result["info"])
	}

	if info["name"] != "test" {
		t.Errorf("expected name='test', got %v", info["name"])
	}

	if info["value"] != float64(42) {
		t.Errorf("expected value=42, got %v", info["value"])
	}
}
