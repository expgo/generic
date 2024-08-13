package gmap

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		m   *Map[string, int]
		key string
	}
	tests := []struct {
		name   string
		args   args
		wantV  int
		wantOk bool
		setup  func(args)
	}{
		{
			name: "ExistingKey",
			args: args{
				m:   NewMap[string, int](),
				key: "one",
			},
			wantV:  1,
			wantOk: true,
			setup: func(a args) {
				Store(a.m, "one", 1)
			},
		},
		{
			name: "NonExistingKey",
			args: args{
				m:   NewMap[string, int](),
				key: "two",
			},
			wantV:  0,
			wantOk: false,
			setup:  func(a args) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			tt.setup(tt.args)
			gotV, gotOk := Load(tt.args.m, tt.args.key)
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("Load() gotV = %v, want %v", gotV, tt.wantV)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestLoadAndDelete(t *testing.T) {
	t.Run("successful operations", func(t *testing.T) {
		m := NewMap[int, string]()
		Store(m, 1, "value1")
		Store(m, 2, "value2")

		tests := []struct {
			name string
			key  int
			want string
		}{
			{name: "Delete existing item", key: 1, want: "value1"},
			{name: "Delete another existing item", key: 2, want: "value2"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, ok := LoadAndDelete(m, tt.key)
				if !ok || !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LoadAndDelete() = %v, want %v", got, tt.want)
				}
				_, exists := Load(m, tt.key)
				if exists {
					t.Errorf("LoadAndDelete() failed to delete key = %v", tt.key)
				}
			})
		}
	})

	t.Run("unsuccessful operations", func(t *testing.T) {
		m := NewMap[int, string]()
		Store(m, 1, "value1")

		tests := []struct {
			name string
			key  int
		}{
			{name: "Delete non-existing item", key: 2},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, ok := LoadAndDelete(m, tt.key)
				if ok {
					t.Errorf("LoadAndDelete() = true, want false for key = %v", tt.key)
				}
			})
		}
	})
}

func TestLoadOrStore(t *testing.T) {
	type key struct {
		id   int
		name string
	}

	tests := []struct {
		name  string
		input struct {
			key   key
			value string
		}
		expect struct {
			value  string
			loaded bool
		}
	}{
		{
			name: "key not previously present",
			input: struct {
				key   key
				value string
			}{key: key{id: 1, name: "one"}, value: "valueOne"},
			expect: struct {
				value  string
				loaded bool
			}{value: "valueOne", loaded: false},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMap[key, string]()

			// Call the function first time
			value, loaded := LoadOrStore(m, tc.input.key, tc.input.value)

			// Assert first call
			assert.Equal(t, tc.expect.value, value)
			assert.Equal(t, tc.expect.loaded, loaded)

			// Call the function second time
			value2, loaded2 := LoadOrStore(m, tc.input.key, tc.input.value)

			// Assert second call
			assert.Equal(t, tc.input.value, value2)
			assert.True(t, loaded2)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		input    map[int]int // Modify the key and value  type according to your actual requirements
		key      int         // Modify the key type according to your actual requirements
		expected map[int]int // Modify the key and value type according to your actual requirements
	}{
		{
			name:     "delete existing element",
			input:    map[int]int{1: 1, 2: 2, 3: 3},
			key:      2,
			expected: map[int]int{1: 1, 3: 3},
		},
		{
			name:     "delete non-existing element",
			input:    map[int]int{1: 1, 2: 2, 3: 3},
			key:      4,
			expected: map[int]int{1: 1, 2: 2, 3: 3},
		},
		{
			name:     "delete from empty map",
			input:    map[int]int{},
			key:      1,
			expected: map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMap[int, int]()
			for key, value := range tt.input {
				Store(m, key, value)
			}

			Delete(m, tt.key)

			// Check if elements match expected

			for key := range tt.expected {
				if _, present := Load(m, key); !present {
					t.Errorf("Key %v not found in the map, but expected to be present", key)
				}
			}

			// Check if there are no extra elements

			Range(m, func(key int, value int) bool { // Modify the key and value type according to your actual requirements
				if _, present := tt.expected[key]; !present {
					t.Errorf("Key %v found in the map, but not expected to be present", key)
				}

				// Stop iteration
				return false
			})
		})
	}
}

func TestSwap(t *testing.T) {
	var tests = []struct {
		name     string
		key      string
		value    string
		previous string
		loaded   bool
	}{
		{name: "new", key: "k0", value: "v0", previous: "", loaded: false},
		{name: "existing", key: "k1", value: "v2", previous: "v1", loaded: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetMap := NewMap[string, string]()
			Store(targetMap, "k1", "v1")

			previous, loaded := Swap(targetMap, tt.key, tt.value)

			if previous != tt.previous {
				t.Errorf("got previous %q, want %q", previous, tt.previous)
			}

			if loaded != tt.loaded {
				t.Errorf("got loaded %t, want %t", loaded, tt.loaded)
			}

			actual, loaded := Load(targetMap, tt.key)
			if tt.value != "" && loaded {
				if actual != tt.value {
					t.Errorf("got actual %q, want %q", actual, tt.value)
				}
			} else if loaded {
				t.Errorf("loaded %t, want false", loaded)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		setupMap func() *Map[int, int]
		want     int
	}{
		{
			name: "empty map",
			setupMap: func() *Map[int, int] {
				return NewMap[int, int]()
			},
			want: 0,
		},
		{
			name: "map with one item",
			setupMap: func() *Map[int, int] {
				m := NewMap[int, int]()
				Store(m, 1, 1)
				return m
			},
			want: 1,
		},
		{
			name: "map with multiple items",
			setupMap: func() *Map[int, int] {
				m := NewMap[int, int]()
				Store(m, 1, 1)
				Store(m, 2, 2)
				Store(m, 3, 3)
				return m
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupMap()
			got := Size(m)
			assert.Equal(t, tt.want, got)
		})
	}
}
