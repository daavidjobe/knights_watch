package watcher

import (
	"os"
	"testing"
)

func TestDescribe(t *testing.T) {
	e := Event{EventType: Create, Path: "watcher.go"}

	real, _ := os.Stat("watcher.go")
	fake, _ := os.Stat("fake.go")

	testCases := []struct {
		info     os.FileInfo
		expected string
	}{
		{
			fake,
			"?",
		},
		{
			real,
			"CREATE -> watcher.go",
		},
	}
	for _, tc := range testCases {
		e.FileInfo = tc.info
		str, err := e.Describe()
		if err != nil && tc.expected != "?" {
			t.Error(err)
		} else if str != tc.expected {
			t.Errorf("expected e.Describe() to be %s, got %s", tc.expected, str)
		}
	}
}
