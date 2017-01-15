package winrmcp

import "testing"

func TestEmptyWinPath(t *testing.T) {
	if winPath("") != "" {
		t.Error("Empty String should be preserved!")
	}
}

func TestWinPathWithSpaces(t *testing.T) {
	var testPathes = []struct {
		input    string
		expected string
	}{
		{"C:\\Users\\test user\\", "'C:\\Users\\test user\\'"},
		{"\"C:\\Users\\test user\\\"", "'C:\\Users\\test user\\'"},
		{"'C:\\Users\\test user\\'", "'C:\\Users\\test user\\'"},
	}

	for _, testpair := range testPathes {
		actual := winPath(testpair.input)
		if actual != testpair.expected {
			t.Errorf("Expected : \"%s\" [actual \"%s\"]", testpair.expected, actual)
		}
	}
}

func TestWinPathDelimiterReplacement(t *testing.T) {
	var testPathes = []struct {
		input    string
		expected string
	}{
		{"C:/Users/testuser\\", "C:\\Users\\testuser\\"},
		{"C:\\Users/testuser/", "C:\\Users\\testuser\\"},
		{"C:\\Users/test user", "'C:\\Users\\test user'"},
	}

	for _, testpair := range testPathes {
		actual := winPath(testpair.input)
		if actual != testpair.expected {
			t.Errorf("Expected : \"%s\" [actual \"%s\"]", testpair.expected, actual)
		}
	}
}
