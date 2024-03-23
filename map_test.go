package generic

import (
	"github.com/stretchr/testify/assert"
	"sync"
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

func TestMapLoad(t *testing.T) {
	type KeyType string
	type ValueType int

	tests := []struct {
		name    string
		key     KeyType
		value   ValueType
		loadKey KeyType
		want    ValueType
		exist   bool
	}{
		{"Existing Key", "key1", 123, "key1", 123, true},
		{"Non Existing Key", "key1", 123, "key2", 0, false},
		{"Zero Value Key", "key1", 0, "key1", 0, true},
		{"Empty Key", "", 123, "", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[KeyType, ValueType]{}
			m.Store(tt.key, tt.value)

			got, exist := m.Load(tt.loadKey)
			if exist != tt.exist {
				t.Errorf("Load() exist = %v, want = %v", exist, tt.exist)
			}
			if got != tt.want {
				t.Errorf("Load() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestMapLoadOrStore(t *testing.T) {
	var testCases = []struct {
		name     string
		key      int
		value    string
		expected string
		found    bool
	}{
		{
			name:     "New key",
			key:      1,
			value:    "hello",
			expected: "hello",
			found:    false,
		},
		{
			name:     "Existing key",
			key:      1,
			value:    "world",
			expected: "hello",
			found:    true,
		},
	}

	mapType := &Map[int, string]{innerMap: sync.Map{}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, found := mapType.LoadOrStore(tc.key, tc.value)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.found, found)
		})
	}
}

func TestMapLoadAndDelete(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name             string
		m                *Map[string, int]
		args             args
		wantV            int
		wantGot          bool
		initialStoreData map[string]int
	}{
		{
			name:             "KeyExists",
			m:                &Map[string, int]{},
			args:             args{k: "test1"},
			wantV:            1,
			wantGot:          true,
			initialStoreData: map[string]int{"test1": 1, "test2": 2},
		},
		{
			name:             "KeyNotExists",
			m:                &Map[string, int]{},
			args:             args{k: "test3"},
			wantV:            0,
			wantGot:          false,
			initialStoreData: map[string]int{"test1": 1, "test2": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.initialStoreData {
				tt.m.Store(k, v)
			}
			gotV, gotGot := tt.m.LoadAndDelete(tt.args.k)
			if gotV != tt.wantV {
				t.Errorf("Map.LoadAndDelete() got Value = %v, want %v", gotV, tt.wantV)
			}
			if gotGot != tt.wantGot {
				t.Errorf("Map.LoadAndDelete() got Got = %v, want %v", gotGot, tt.wantGot)
			}

			// Confirm it's deleted
			_, got := tt.m.Load(tt.args.k)
			if got {
				t.Error("Map.LoadAndDelete() failed to delete key")
			}
		})
	}
}

func TestMapSwap(t *testing.T) {
	type test struct {
		name      string
		key       int
		value     string
		toSwap    string
		expectedV string
		expectedB bool
	}

	tests := []test{
		{"Key exists", 1, "Value1", "ValueSwapped", "Value1", true},
		{"Key not exists", 2, "", "ValueSwapped", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map[int, string]{}
			if tt.value != "" {
				m.Store(tt.key, tt.value)
			}
			v, b := m.Swap(tt.key, tt.toSwap)

			assert.Equal(t, tt.expectedV, v)
			assert.Equal(t, tt.expectedB, b)
		})
	}
}
