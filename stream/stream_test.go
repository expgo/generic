package stream

import (
	"reflect"
	"testing"
)

func TestFilter(t *testing.T) {
	type test struct {
		name     string
		slice    []int
		filter   func(int) bool
		expected []int
	}

	tests := []test{
		{
			name:  "filter_out_odds",
			slice: []int{1, 2, 3, 4, 5},
			filter: func(n int) bool {
				return n%2 == 0
			},
			expected: []int{2, 4},
		},
		{
			name:  "filter_out_zeros",
			slice: []int{0, 1, 0, 3, 0, 5},
			filter: func(n int) bool {
				return n != 0
			},
			expected: []int{1, 3, 5},
		},
		{
			name:  "filter_out_all",
			slice: []int{1, 2, 3, 4, 5},
			filter: func(n int) bool {
				return n > 5
			},
			expected: []int{},
		},
		{
			name:  "empty_slice",
			slice: []int{},
			filter: func(n int) bool {
				return true
			},
			expected: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := Filter(tc.slice, tc.filter)
			if len(res) != len(tc.expected) {
				t.Fatalf("expected length: %v, got: %v", len(tc.expected), len(res))
			}
			for i, v := range res {
				if v != tc.expected[i] {
					t.Fatalf("expected item %v to be %v, got: %v", i, tc.expected[i], v)
				}
			}
		})
	}
}

func TestMap(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		input    []int
		mapFunc  func(int) int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			mapFunc:  func(x int) int { return x * 2 },
			expected: nil,
		},
		{
			name:     "single element",
			input:    []int{1},
			mapFunc:  func(x int) int { return x * 2 },
			expected: []int{2},
		},
		{
			name:     "two elements",
			input:    []int{1, 2},
			mapFunc:  func(x int) int { return x * 2 },
			expected: []int{2, 4},
		},
		{
			name:     "negative numbers",
			input:    []int{-1, -2},
			mapFunc:  func(x int) int { return x * 2 },
			expected: []int{-2, -4},
		},
		{
			name:     "map function that adds",
			input:    []int{1, 2},
			mapFunc:  func(x int) int { return x + 1 },
			expected: []int{2, 3},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Call the function and check the output
			result := Map(testCase.input, testCase.mapFunc)
			if !reflect.DeepEqual(result, testCase.expected) {
				t.Errorf("Failed test '%s': got %v, expected %v", testCase.name, result, testCase.expected)
			}
		})
	}
}

func TestMapToAny(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []interface{}
	}{
		{
			name:  "EmptySlice",
			input: []int{},
			want:  nil,
		},
		{
			name:  "IntSlice",
			input: []int{1, 2, 3},
			want:  []interface{}{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapToAny(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToAny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShuffle(t *testing.T) {
	// Function for generating sample data
	generateData := func(n int) []int {
		data := make([]int, n)
		for i := 0; i < n; i++ {
			data[i] = i
		}
		return data
	}

	testCases := []struct {
		name   string
		elems  []int
		expErr bool
	}{
		{
			name:  "Three Elements",
			elems: []int{1, 2, 3},
		},
		{
			name:  "Four Elements",
			elems: []int{1, 2, 3, 4},
		},
		{
			name:  "Multiple Elements",
			elems: generateData(1000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			shuffled := Shuffle(tc.elems)

			// Check if shuffle returns a new stream object
			if reflect.DeepEqual(tc.elems, shuffled) {
				// try once again
				shuffled = Shuffle(tc.elems)
				if reflect.DeepEqual(tc.elems, shuffled) {
					t.Errorf("Shuffle() must return new stream object")
				}
			}

			// Check the number of elements is same in the original and shuffled stream
			if got, want := len(tc.elems), len(shuffled); got != want {
				t.Errorf("len(shuffled) got %v, want %v", got, want)
			}

			// Check at least one element is in a different position
			var found bool
			for i, v := range tc.elems {
				if v != shuffled[i] {
					found = true
					break
				}
			}

			if !found {
				t.Error("Shuffle() should alter the order of the elements")
			}
		})
	}
}
