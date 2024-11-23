package activewindowreceiver

import (
	"reflect"
	"testing"
)

func TestParseWindowName(t *testing.T) {
	testCases := []struct {
		name           string
		windowName     string
		expectedValues []string
	}{
		{
			name:       "single title",
			windowName: "ChatAAT — Mozilla Firefox",
			expectedValues: []string{
				"ChatAAT",
				"Mozilla Firefox",
			},
		},
		{
			name:       "title with platform host",
			windowName: "Window Time tracker - SomeApp — Mozilla Firefox",
			expectedValues: []string{
				"Window Time tracker",
				"SomeApp",
				"Mozilla Firefox",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			values := parseWindowName(testCase.windowName)

			if !reflect.DeepEqual(values, testCase.expectedValues) {
				println(len(values))
				t.Errorf("expected: %v, got: %v\n", testCase.expectedValues, values)
			}
		})
	}
}
