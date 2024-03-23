package generic

import (
	"errors"
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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resultVal string
			var resultErr error

			defer func() {
				if r := recover(); r != nil {
					resultErr = r.(error)
				}
			}()

			resultVal, resultErr = cache.GetOrLoad(tc.key, tc.loadFunc)

			assert.Equal(t, tc.expectedResult, resultVal)
			assert.Equal(t, tc.expectedErr, resultErr)
		})
	}
}
