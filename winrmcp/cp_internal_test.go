package winrmcp

import (
	"regexp"
	"testing"
)

func TestTempFilename(t *testing.T) {
	var validator = regexp.MustCompile(`^winrmcp-[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}\.tmp$`)
	actual, fnameError := tempFileName()
	if fnameError != nil {
		t.Errorf("Can't create temp filename! Error = \"%s\"", fnameError.Error())
		return
	}
	if !validator.MatchString(actual) {
		t.Errorf("Invalid Temp Filename: \"%s\"", actual)
	}
}

func TestChunkSizeCalculation(t *testing.T) {
	testdata := []int{
		1,
		20,
		21,
		24,
		23,
		49,
		13,
	}
	for _, testpoint := range testdata {
		actualSize := chunkSize(testpoint)
		if (actualSize % 3) != 0 {
			t.Error("ChunkSize is not a multiple of 3")
		}
		if actualSize >= 8192 {
			t.Errorf("ChunkSize of %d exceeds maximum console characters", actualSize)
		}
	}
}
