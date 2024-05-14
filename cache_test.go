package generic

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrLoad(t *testing.T) {

	loadFunc := func(k string) (string, error) {
		return "value for " + k, nil
	}

	cache := &Cache[string, string]{}

	testCases := []struct {
		name           string
		key            string
		loadFunc       func(k string) (string, error)
		expectedErr    error
		expectedResult string
	}{
		{
			name:           "Existing Key",
			key:            "testKey",
			loadFunc:       loadFunc,
			expectedErr:    nil,
			expectedResult: "value for testKey",
		},
		{
			name: "Load function returns error",
			key:  "anyKey",
			loadFunc: func(k string) (string, error) {
				return "", errors.New("load function error")
			},
			expectedErr:    errors.New("load function error"),
			expectedResult: "",
		},
		{
			name:           "load func is nil",
			key:            "testKey",
			loadFunc:       nil,
			expectedErr:    errors.New("load function must not be nil"),
			expectedResult: "value for testKey",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resultVal string
			var resultErr error

			defer func() {
				if r := recover(); r != nil {
					resultErr = r.(error)
				}
				assert.Equal(t, tc.expectedErr, resultErr)
			}()

			resultVal, resultErr = cache.GetOrLoad(tc.key, tc.loadFunc)

			assert.Equal(t, tc.expectedResult, resultVal)
			assert.Equal(t, tc.expectedErr, resultErr)
		})
	}
}

func TestCache_Evict(t *testing.T) {
	type test struct {
		cache   *Cache[string, int]
		key     string
		want    bool
		comment string
	}

	tests := []test{
		{
			cache:   &Cache[string, int]{innerMap: sync.Map{}},
			key:     "key1",
			want:    false,
			comment: "empty cache",
		},
		{
			cache: &Cache[string, int]{innerMap: func() sync.Map {
				m := sync.Map{}
				m.Store("key1", 1)
				return m
			}()},
			key:     "key1",
			want:    true,
			comment: "key exists in cache",
		},
		{
			cache: &Cache[string, int]{innerMap: func() sync.Map {
				m := sync.Map{}
				m.Store("key1", 1)
				return m
			}()},
			key:     "key2",
			want:    false,
			comment: "key does not exist in cache",
		},
	}
	for _, tt := range tests {
		t.Run(tt.comment, func(t *testing.T) {
			if got := tt.cache.Evict(tt.key); got != tt.want {
				t.Errorf("Cache.Evict() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Clear(t *testing.T) {
	type test struct {
		name  string
		load  func(k int) (int, error)
		keys  []int
		evict []int
	}
	tests := []test{
		{
			name: "empty cache",
			load: func(k int) (int, error) {
				return k * 2, nil
			},
			keys: nil,
		},
		{
			name: "single item cache",
			load: func(k int) (int, error) {
				return k * 2, nil
			},
			keys: []int{1},
		},
		{
			name: "multiple item cache",
			load: func(k int) (int, error) {
				return k * 2, nil
			},
			keys: []int{1, 2, 3, 4, 5},
		},
		{
			name: "clear does not affect other caches",
			load: func(k int) (int, error) {
				return k * 2, nil
			},
			keys:  []int{1, 2, 3, 4, 5},
			evict: []int{1, 3, 5},
		},
	}

	for _, testInstance := range tests {
		t.Run(testInstance.name, func(t *testing.T) {
			cache := Cache[int, int]{}
			// store initial data
			for _, k := range testInstance.keys {
				cache.GetOrLoad(k, testInstance.load)
			}
			cache.Clear()
			for _, k := range testInstance.keys {
				if _, err := cache.GetOrLoad(k, testInstance.load); err != nil {
					t.Errorf("key %d doesn't exist in the cache but it should", k)
				}
			}
			// evict some keys from the cache
			for _, k := range testInstance.evict {
				cache.Evict(k)
			}
			// verify the evicted keys
			for _, k := range testInstance.evict {
				if _, err := cache.GetOrLoad(k, func(int) (int, error) {
					return 0, errors.New("not found")
				}); err == nil {
					t.Errorf("key %d exists in the cache but it shouldn't", k)
				}
			}
		})
	}

}
