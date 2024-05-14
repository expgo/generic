package stream

import (
	"errors"
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

func TestMustMap(t *testing.T) {
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
			result := MustMap(testCase.input, testCase.mapFunc)
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
			if got := ToAny(tt.input); !reflect.DeepEqual(got, tt.want) {
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
			name:  "Empty Elements",
			elems: []int{},
		},
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

			if len(shuffled) == 0 {
				return
			}

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

func TestLimit(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		n    int
		want []int
	}{
		{
			name: "LimitGreaterThanLength",
			s:    []int{1, 2, 3, 4, 5},
			n:    7,
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "LimitLessThanLength",
			s:    []int{1, 2, 3, 4, 5},
			n:    3,
			want: []int{1, 2, 3},
		},
		{
			name: "LimitZero",
			s:    []int{1, 2, 3, 4, 5},
			n:    0,
			want: []int{},
		},
		{
			name: "LimitEqualToLength",
			s:    []int{1, 2, 3, 4, 5},
			n:    5,
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "EmptySlice",
			s:    []int{},
			n:    3,
			want: []int{},
		},
		{
			name: "NegativeLimit",
			s:    []int{1, 2, 3, 4, 5},
			n:    -3,
			want: []int{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Limit(test.s, test.n); !equal(got, test.want) {
				t.Errorf("Limit(%v, %v) = %v, want %v", test.s, test.n, got, test.want)
			}
		})
	}
}

// A helper function to compare slices.
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestSkip(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{
			name:     "Skip_Zero_Items",
			slice:    []int{1, 2, 3, 4, 5},
			n:        0,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Skip_Some_Items",
			slice:    []int{1, 2, 3, 4, 5},
			n:        3,
			expected: []int{4, 5},
		},
		{
			name:     "Skip_All_Items",
			slice:    []int{1, 2, 3, 4, 5},
			n:        5,
			expected: []int{},
		},
		{
			name:     "Skip_More_Than_Length_Items",
			slice:    []int{1, 2, 3, 4, 5},
			n:        7,
			expected: []int{},
		},
		{
			name:     "Skip_Negative_Items",
			slice:    []int{1, 2, 3, 4, 5},
			n:        -5,
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Skip(tt.slice, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v\n", tt.expected, result)
			}
		})
	}
}

func TestAllMatch(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		matchElem int
		want      bool
	}{
		{
			name:      "all elements match",
			input:     []int{1, 1, 1},
			matchElem: 1,
			want:      true,
		},
		{
			name:      "not all elements match",
			input:     []int{1, 2, 3},
			matchElem: 1,
			want:      false,
		},
		{
			name:      "empty slice",
			input:     []int{},
			matchElem: 1,
			want:      true,
		},
		{
			name:      "single element slice, match",
			input:     []int{3},
			matchElem: 3,
			want:      true,
		},
		{
			name:      "single element slice, no match",
			input:     []int{2},
			matchElem: 3,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllMatch(tt.input, tt.matchElem); got != tt.want {
				t.Errorf("AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAllMatchFunc(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		matchFunc func(int) bool
		want      bool
	}{
		{
			"All elements match",
			[]int{2, 4, 6, 8, 10},
			func(n int) bool { return n%2 == 0 },
			true,
		},
		{
			"Not all elements match ",
			[]int{2, 3, 6, 8, 10},
			func(n int) bool { return n%2 == 0 },
			false,
		},
		{
			"Empty slice",
			[]int{},
			func(n int) bool { return n%2 == 0 },
			true,
		},
		{
			"Single match",
			[]int{2},
			func(n int) bool { return n%2 == 0 },
			true,
		},
		{
			"Single mismatch",
			[]int{3},
			func(n int) bool { return n%2 == 0 },
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllMatchFunc(tt.input, tt.matchFunc); got != tt.want {
				t.Fatalf("AllMatchFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyMatch(t *testing.T) {
	// table-driven test cases
	tests := []struct {
		name   string
		input  []int
		target int
		want   bool
	}{
		{
			name:   "Non-Empty Slice, Target Exists",
			input:  []int{1, 2, 3, 4, 5},
			target: 3,
			want:   true,
		},
		{
			name:   "Non-Empty Slice, Target Does Not Exist",
			input:  []int{1, 2, 3, 4, 5},
			target: 6,
			want:   false,
		},
		{
			name:   "Empty Slice",
			input:  []int{},
			target: 1,
			want:   false,
		},
		{
			name:   "Slice With Duplicates, Target Exists",
			input:  []int{1, 2, 2, 3, 3},
			target: 2,
			want:   true,
		},
		{
			name:   "Slice With Duplicates, Target Does Not Exist",
			input:  []int{1, 2, 2, 3, 3},
			target: 4,
			want:   false,
		},
	}

	// running test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := AnyMatch(tc.input, tc.target)
			if got != tc.want {
				t.Errorf("Expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestAnyMatchFunc(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	containsN := func(n string) func(string) bool {
		return func(s string) bool { return s == n }
	}

	tests := []struct {
		name      string
		slice     interface{}
		matchFunc interface{}
		want      bool
	}{
		{
			name:      "WithIntegersAndEvenMatchFunction",
			slice:     []int{1, 3, 5, 7, 9, 2},
			matchFunc: isEven,
			want:      true,
		},
		{
			name:      "WithIntegersAndNoEvenMatchFunction",
			slice:     []int{1, 3, 5, 7, 9},
			matchFunc: isEven,
			want:      false,
		},
		{
			name:      "WithStringSliceAndValidValue",
			slice:     []string{"Hello", "World", "Goland"},
			matchFunc: containsN("Goland"),
			want:      true,
		},
		{
			name:      "WithStringSliceAndInvalidValue",
			slice:     []string{"Hello", "World", "Goland"},
			matchFunc: containsN("Test"),
			want:      false,
		},
		{
			name:      "WithEmptySlice",
			slice:     []int{},
			matchFunc: isEven,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			switch s := tt.slice.(type) {
			case []int:
				got = AnyMatchFunc(s, tt.matchFunc.(func(int) bool))
			case []string:
				got = AnyMatchFunc(s, tt.matchFunc.(func(string) bool))
			}
			if got != tt.want {
				t.Errorf("AnyMatchFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		transform func(int) (int, error)
		expected  []int
		expectErr bool
	}{
		{
			name:  "double values",
			input: []int{1, 2, 3},
			transform: func(n int) (int, error) {
				return n * 2, nil
			},
			expected:  []int{2, 4, 6},
			expectErr: false,
		},
		{
			name:  "error scenario",
			input: []int{1, 2, 3},
			transform: func(n int) (int, error) {
				return 0, errors.New("transformation error")
			},
			expected:  nil,
			expectErr: true,
		},
		{
			name:  "empty slice",
			input: []int{},
			transform: func(n int) (int, error) {
				return n * 2, nil
			},
			expected:  []int{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map(tt.input, tt.transform)
			if (err != nil) != tt.expectErr {
				t.Errorf("Map() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("Map() got = %v, want %v", got, tt.expected)
			}
			for i, val := range got {
				if val != tt.expected[i] {
					t.Errorf("Map() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	var getKeyFunc = func(i int) int {
		return i % 2
	}

	tests := []struct {
		name   string
		s      []int
		getKey func(int) int
		want   map[int][]int
	}{
		{
			name:   "Empty slice",
			s:      []int{},
			getKey: getKeyFunc,
			want:   make(map[int][]int),
		},
		{
			name:   "Slice with single element",
			s:      []int{5},
			getKey: getKeyFunc,
			want:   map[int][]int{1: {5}},
		},
		{
			name:   "Slice with multiple elements",
			s:      []int{1, 2, 3, 4, 5},
			getKey: getKeyFunc,
			want:   map[int][]int{0: {2, 4}, 1: {1, 3, 5}},
		},
		{
			name:   "Slice with duplicate elements",
			s:      []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4},
			getKey: getKeyFunc,
			want:   map[int][]int{0: {2, 2, 4, 4, 4, 4}, 1: {1, 3, 3, 3}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := GroupBy(test.s, test.getKey)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GroupBy() = %v, want %v", got, test.want)
			}
		})
	}
}
