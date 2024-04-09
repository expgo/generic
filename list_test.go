package generic

import (
	"reflect"
	"testing"
)

func TestList_Add(t *testing.T) {
	type args struct {
		e int
	}
	tests := []struct {
		name string
		list *List[int]
		args args
		want *List[int]
	}{
		{
			name: "Adding to an empty list",
			list: &List[int]{},
			args: args{e: 1},
			want: &List[int]{items: []int{1}},
		},
		{
			name: "Adding to a non-empty list",
			list: &List[int]{items: []int{1, 2, 3}},
			args: args{e: 4},
			want: &List[int]{items: []int{1, 2, 3, 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.list.Add(tt.args.e)
			if !reflect.DeepEqual(tt.list, tt.want) {
				t.Errorf("Add() = %v, want %v", tt.list, tt.want)
			}
		})
	}
}

func TestList_Remove(t *testing.T) {

	type Test struct {
		name     string
		elem     int
		insert   []int
		remove   int
		expected []int
	}

	tests := []Test{
		{
			name:     "remove from middle",
			insert:   []int{1, 2, 3, 4, 5},
			remove:   3,
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "remove from beginning",
			insert:   []int{1, 2, 3, 4, 5},
			remove:   1,
			expected: []int{2, 3, 4, 5},
		},
		{
			name:     "remove from end",
			insert:   []int{1, 2, 3, 4, 5},
			remove:   5,
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "remove non-existent element",
			insert:   []int{1, 2, 3, 4, 5},
			remove:   6,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "remove from empty list",
			insert:   []int{},
			remove:   1,
			expected: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			list := &List[int]{}
			for _, n := range tc.insert {
				list.Add(n)
			}

			list.Remove(tc.remove)

			for i, item := range list.items {
				if item != tc.expected[i] {
					t.Fatalf("Expected %v, got %v", tc.expected, list.items)
				}
			}
		})
	}
}

func TestList_Contains(t *testing.T) {
	type testCase struct {
		name       string
		init       []string
		target     string
		wantResult bool
	}

	testCases := []testCase{
		{
			name:       "emptyList",
			init:       []string{},
			target:     "test",
			wantResult: false,
		},
		{
			name:       "singleElementPresent",
			init:       []string{"test"},
			target:     "test",
			wantResult: true,
		},
		{
			name:       "singleElementAbsent",
			init:       []string{"test"},
			target:     "absent",
			wantResult: false,
		},
		{
			name:       "multipleElementsPresent",
			init:       []string{"one", "two", "three", "one", "four"},
			target:     "one",
			wantResult: true,
		},
		{
			name:       "multipleElementsAbsent",
			init:       []string{"one", "two", "three", "one", "four"},
			target:     "absent",
			wantResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			list := List[string]{}
			for _, elem := range tc.init {
				list.Add(elem)
			}

			if got := list.Contains(tc.target); got != tc.wantResult {
				t.Errorf("List.Contains = %v; want %v", got, tc.wantResult)
			}
		})
	}
}

func TestList_RemoveAt(t *testing.T) {
	tests := []struct {
		name    string
		input   []int
		remove  int
		expect  []int
		success bool
	}{
		{
			name:    "Remove_One_Middle",
			input:   []int{1, 2, 3, 4, 5},
			remove:  2,
			expect:  []int{1, 2, 4, 5},
			success: true,
		},
		{
			name:    "Remove_One_First",
			input:   []int{1, 2, 3, 4, 5},
			remove:  0,
			expect:  []int{2, 3, 4, 5},
			success: true,
		},
		{
			name:    "Remove_One_Last",
			input:   []int{1, 2, 3, 4, 5},
			remove:  4,
			expect:  []int{1, 2, 3, 4},
			success: true,
		},
		{
			name:    "Remove_Empty",
			input:   []int{},
			remove:  0,
			expect:  []int{},
			success: false,
		},
		{
			name:    "Remove_Invalid",
			input:   []int{1, 2, 3, 4, 5},
			remove:  10,
			expect:  []int{1, 2, 3, 4, 5},
			success: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			list := List[int]{items: tc.input}
			result := list.RemoveAt(tc.remove)

			if result != tc.success {
				t.Errorf("Expected success %v but got %v", tc.success, result)
			}

			if len(list.items) != len(tc.expect) {
				t.Errorf("Expected length %d but got %d", len(tc.expect), len(list.items))
			}

			for i, item := range list.items {
				if item != tc.expect[i] {
					t.Errorf("Expected item %d but got %d", tc.expect[i], item)
				}
			}
		})
	}
}

func TestList_Clear(t *testing.T) {
	type fields struct {
		items []int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "EmptyList",
			fields: fields{
				items: []int{},
			},
		},
		{
			name: "SingleItemList",
			fields: fields{
				items: []int{1},
			},
		},
		{
			name: "MultiItemList",
			fields: fields{
				items: []int{1, 2, 3, 4, 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &List[int]{
				items: tt.fields.items,
			}
			l.Clear()

			if size := l.Size(); size != 0 {
				t.Errorf("List.Clear() = %v, want %v", size, 0)
			}
		})
	}
}
