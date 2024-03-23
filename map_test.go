package generic

import (
	"testing"
)

func TestFromMap(t *testing.T) {
	tests := map[string]struct {
		input map[int]string
	}{
		"Empty": {
			input: map[int]string{},
		},
		"SingleElement": {
			input: map[int]string{1: "one"},
		},
		"MultipleElements": {
			input: map[int]string{1: "one", 2: "two", 3: "three"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := FromMap(tc.input)

			if got.Size() != len(tc.input) {
				t.Errorf("expected map length %d, but got %d", len(tc.input), got.Size())
			}

			for k, v := range tc.input {
				if val, ok := got.Load(k); !ok || val != v {
					t.Errorf("expected key %d to have value %s, but got %v", k, v, val)
				}
			}
		})
	}
}
