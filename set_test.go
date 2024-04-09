package generic

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSet_Add(t *testing.T) {
	type test struct {
		name       string
		addElement int
		wantLoaded bool
	}

	tests := []test{
		{
			name:       "Adding new element",
			addElement: 5,
			wantLoaded: true,
		},
		{
			name:       "Re-Adding existing element",
			addElement: 5,
			wantLoaded: false,
		},
	}
	set := &Set[int]{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loaded := set.Add(tt.addElement)

			if loaded != tt.wantLoaded {
				t.Errorf("Add() loaded = %v, wantLoaded %v", loaded, tt.wantLoaded)
			}

			contains := set.Contains(tt.addElement)
			if !contains {
				t.Errorf("Add() element was not correctly added to the set.")
			}
		})
	}
}

func TestSet_Contains(t *testing.T) {
	type fields struct {
		elems Set[int]
	}
	tests := []struct {
		name   string
		fields fields
		input  int
		want   bool
	}{
		{
			name: "ValueExistsInSet",
			fields: fields{
				elems: func() Set[int] {
					var mt Set[int]
					mt.Add(1)
					return mt
				}(),
			},
			input: 1,
			want:  true,
		},
		{
			name: "ValueDoesNotExistsInSet",
			fields: fields{
				elems: func() Set[int] {
					var mt Set[int]
					mt.Add(1)
					mt.Add(2)
					return mt
				}(),
			},
			input: 3,
			want:  false,
		},
		{
			name: "ValueAddedThenRemoved",
			fields: fields{
				elems: func() Set[int] {
					var mt Set[int]
					mt.Add(1)
					mt.Remove(1)
					return mt
				}(),
			},
			input: 1,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.elems.Contains(tt.input); got != tt.want {
				t.Errorf("Set.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	// initialize test cases
	testCases := []struct {
		name     string
		init     []int
		remove   []int
		expected []int
	}{
		{"Remove single", []int{1, 2, 3}, []int{2}, []int{1, 3}},
		{"Remove all", []int{1, 2, 3}, []int{1, 2, 3}, []int{}},
		{"Remove empty", []int{}, []int{2}, []int{}},
		{"Remove non-existing", []int{1, 2, 3}, []int{0}, []int{1, 2, 3}},
		{"Remove duplicate", []int{1, 1, 2, 3}, []int{1}, []int{2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create set and populate with initial data
			mySet := Set[int]{}
			for _, elem := range tc.init {
				mySet.Add(elem)
			}

			// remove elements
			for _, elem := range tc.remove {
				mySet.Remove(elem)
			}

			// check remaining elements against expected
			for _, elem := range tc.expected {
				if !mySet.Contains(elem) {
					t.Errorf("Expected set to contain all of %v after Remove. set: %v", tc.expected, mySet.ToStream())
				}
			}

			if mySet.Size() != len(tc.expected) {
				t.Errorf("Expected set size to be %v, but got %v", len(tc.expected), mySet.Size())
			}
		})
	}
}

func TestSet_Size(t *testing.T) {
	type testCase struct {
		name string
		set  Set[int]
		want int
	}

	tests := []testCase{
		{
			name: "empty set",
			set:  Set[int]{},
			want: 0,
		},
		{
			name: "single item set",
			set: func() Set[int] {
				var s Set[int]
				s.Add(1)
				return s
			}(),
			want: 1,
		},
		{
			name: "multiple distinct items in set",
			set: func() Set[int] {
				var s Set[int]
				s.Add(1)
				s.Add(2)
				s.Add(3)
				return s
			}(),
			want: 3,
		},
		{
			name: "multiple same items in set",
			set: func() Set[int] {
				var s Set[int]
				s.Add(1)
				s.Add(1)
				s.Add(1)
				return s
			}(),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Clear(t *testing.T) {
	type testCase struct {
		name    string
		initial []int
	}
	cases := []testCase{
		{"Empty Set", []int{}},
		{"Single element Set", []int{1}},
		{"Multiple element Set", []int{1, 2, 3, 4, 5}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			set := Set[int]{}
			for _, i := range c.initial {
				set.Add(i)
			}
			set.Clear()
			assert.Equal(t, true, set.IsEmpty())
		})
	}
}

// Assuming IsEmpty method exists in the set struct which returns true if set is empty.
func (s *Set[T]) IsEmpty() bool {
	isEmpty := true
	s.elemMap.Range(func(key, value interface{}) bool {
		isEmpty = false
		return false // stop iteration as soon as we find an element
	})
	return isEmpty
}

func TestSetToStream(t *testing.T) {
	cases := []struct {
		input    []int
		expected []int
	}{
		{[]int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{-1, -2, -3, -4, -5}, []int{-1, -2, -3, -4, -5}},
		{[]int{}, nil},
		{[]int{0, 0, 0, 0, 0}, []int{0}},
	}

	for _, tt := range cases {
		// Arrange
		set := &Set[int]{}

		for _, v := range tt.input {
			set.Add(v)
		}

		// Act
		stream := set.ToStream()

		// Assert
		result, err := stream.ToSlice()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Failed: got %v, want %v", result, tt.expected)
		}
	}
}
