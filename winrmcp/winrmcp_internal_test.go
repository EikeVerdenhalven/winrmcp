package winrmcp

import (
	"os"
	"testing"
	"time"
)

type MockFileInfo struct {
	name    string
	dirflag bool
}

func (m MockFileInfo) Name() string     { return m.name }
func (MockFileInfo) Size() int64        { return 123 }        // length in bytes for regular files; system-dependent for others
func (MockFileInfo) Mode() os.FileMode  { return 0 }          // file mode bits
func (MockFileInfo) ModTime() time.Time { return time.Now() } // modification time
func (m MockFileInfo) IsDir() bool      { return m.dirflag }  // abbreviation for Mode().IsDir()
func (MockFileInfo) Sys() interface{}   { return nil }        // underlying data source (can return nil)

func TestShouldUploadFile(t *testing.T) {
	testdata := []struct {
		info     MockFileInfo
		expected bool
	}{
		{info: MockFileInfo{name: "text.txt", dirflag: false}, expected: true},
		{info: MockFileInfo{name: "testdir", dirflag: true}, expected: false},
		{info: MockFileInfo{name: ".DS_Store", dirflag: false}, expected: false},
	}
	for _, tdata := range testdata {
		if shouldUploadFile(tdata.info) != tdata.expected {
			t.Fail()
		}
	}
}
