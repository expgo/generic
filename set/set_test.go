package set

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name       string
		set        []int
		element    int
		want       []int
		wantResult bool
	}{
		{
			name:       "Add to empty set",
			set:        []int{},
			element:    1,
			want:       []int{1},
			wantResult: true,
		},
		{
			name:       "Add new element to set",
			set:        []int{1, 2, 3},
			element:    4,
			want:       []int{1, 2, 3, 4},
			wantResult: true,
		},
		{
			name:       "Add existing element to set",
			set:        []int{1, 2, 3},
			element:    2,
			want:       []int{1, 2, 3},
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if ret, gotResult := Add(tt.set, tt.element); gotResult != tt.wantResult {
				t.Errorf("Add() = %v, wantResult %v", gotResult, tt.wantResult)
			} else {
				assert.Equal(t, tt.want, ret)
			}
		})
	}
}

func TestAddFunc(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		e    int
		want []int
		ok   bool
	}{
		{
			name: "EmptySlice",
			s:    []int{},
			e:    1,
			want: []int{1},
			ok:   true,
		},
		{
			name: "ElementExists",
			s:    []int{1, 2, 3},
			e:    2,
			want: []int{1, 2, 3},
			ok:   false,
		},
		{
			name: "ElementNotExists",
			s:    []int{1, 2, 3},
			e:    4,
			want: []int{1, 2, 3, 4},
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := AddFunc(tt.s, tt.e, func(v int) bool {
				return v == tt.e
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddFunc() got = %v, want %v", got, tt.want)
			}
			if ok != tt.ok {
				t.Errorf("AddFunc() ok = %v, want %v", ok, tt.ok)
			}
		})
	}
}
