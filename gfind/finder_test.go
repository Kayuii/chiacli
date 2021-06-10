package gfind

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"testing"
)

func testFind(t *testing.T, filePaths []string, search *regexp.Regexp, want []string) {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, filePath := range filePaths {
		file := filepath.Join(tempDir, filePath)
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatal(err)
		}
		err := ioutil.WriteFile(file, []byte("foobar"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	finder := NewFinder(search)
	got, err := finder.Find(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	var absWant []string
	for _, w := range want {
		absWant = append(absWant, filepath.Join(tempDir, w))
	}

	sort.Strings(got)
	sort.Strings(absWant)

	if !reflect.DeepEqual(got, absWant) {
		t.Errorf("got: %v\n want: %v\n", got, absWant)
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		filePaths []string
		search    *regexp.Regexp
		want      []string
	}{
		{
			[]string{"src/foo/foo.go", "src/bar/bar.go", "src/duck/feathered/feathered.go"},
			regexp.MustCompile("src/foo/foo.go"),
			[]string{"src/foo/foo.go"},
		},
		{
			[]string{"src/foo/foo.go", "src/bar/baz/foo.go", "src/duck/feathered/feathered.go"},
			regexp.MustCompile("foo.go"),
			[]string{"src/foo/foo.go", "src/bar/baz/foo.go"},
		},
		{
			[]string{"src/foo/foo.go"},
			regexp.MustCompile("xyz.go"),
			[]string{},
		},
	}

	for _, test := range tests {
		testFind(t, test.filePaths, test.search, test.want)
	}
}
