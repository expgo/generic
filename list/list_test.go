package list

import (
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		list     []int
		element  int
		expected bool
	}{
		{"empty list", []int{}, 1, false},
		{"list with one element matching", []int{1}, 1, true},
		{"list with one element non-matching", []int{2}, 1, false},
		{"list with multiple elements non-matching", []int{2, 3, 4, 5}, 1, false},
		{"list with multiple elements matching", []int{1, 2, 3, 4, 5}, 1, true},
		{"list with multiple same elements matching", []int{1, 1, 1, 1}, 1, true},
		{"list with multiple same elements non-matching", []int{2, 2, 2, 2}, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.list, tt.element)
			if got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContainsFunc(t *testing.T) {
	tests := []struct {
		name      string
		s         []int
		e         int
		matchFunc func(int) bool
		want      bool
	}{
		{
			name: "Contains",
			s:    []int{1, 2, 3, 4, 5},
			e:    3,
			matchFunc: func(e int) bool {
				return e == 3
			},
			want: true,
		},
		{
			name: "Does Not Contain",
			s:    []int{1, 2, 3, 4, 5},
			e:    6,
			matchFunc: func(e int) bool {
				return e == 6
			},
			want: false,
		},
		{
			name: "Empty Slice",
			s:    []int{},
			e:    1,
			matchFunc: func(e int) bool {
				return e == 1
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsFunc(tt.s, tt.e, tt.matchFunc); got != tt.want {
				t.Errorf("ContainsFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name           string
		s              []int
		e              int
		want           []int
		wantDeleteFlag bool
	}{
		{"delete existing", []int{1, 2, 3}, 3, []int{1, 2}, true},
		{"delete non-existing", []int{1, 2, 3}, 4, []int{1, 2, 3}, false},
		{"delete empty-slice", []int{}, 1, []int{}, false},
		{"delete single-element-slice", []int{1}, 1, []int{}, true},
		{"delete first element", []int{1, 2, 3}, 1, []int{2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotDeleteFlag := Delete(tt.s, tt.e)
			if !reflect.DeepEqual(got, tt.want) || gotDeleteFlag != tt.wantDeleteFlag {
				t.Fatalf("Delete(%v, %v): got (%v, %v), want (%v, %v)",
					tt.s, tt.e, got, gotDeleteFlag, tt.want, tt.wantDeleteFlag)
			}
		})
	}
}

func TestDeleteFunc(t *testing.T) {
	type test struct {
		name     string
		input    []int
		element  int
		expected []int
		found    bool
		match    func(int) bool
	}

	tests := []test{
		{
			name:     "EmptySlice",
			input:    []int{},
			element:  1,
			expected: []int{},
			match:    func(i int) bool { return i == 1 },
			found:    false,
		},
		{
			name:     "SingleMatch",
			input:    []int{1, 2, 3, 4, 5},
			element:  3,
			expected: []int{1, 2, 4, 5},
			match:    func(i int) bool { return i == 3 },
			found:    true,
		},
		{
			name:     "NoMatch",
			input:    []int{1, 2, 4, 5},
			element:  3,
			expected: []int{1, 2, 4, 5},
			match:    func(i int) bool { return i == 3 },
			found:    false,
		},
		{
			name:     "MultiMatch",
			input:    []int{1, 2, 3, 3, 4, 5},
			element:  3,
			expected: []int{1, 2, 3, 4, 5},
			match:    func(i int) bool { return i == 3 },
			found:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, found := DeleteFunc(tc.input, tc.element, tc.match)
			if !compareSlices(result, tc.expected) || found != tc.found {
				t.Errorf("DeleteFunc(%v, %d) = %v, want %v", tc.input, tc.element, result, tc.expected)
			}
		})
	}
}

func compareSlices(s1, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}
	return true
}
