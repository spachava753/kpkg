package download

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testGzipFileFetcher struct {
	gzipFilePath string
}

func (t testGzipFileFetcher) FetchFile(_ string) (string, error) {
	return filepath.Abs(t.gzipFilePath)
}

func Test_gzipFileFetcher_FetchFile_Zip(t *testing.T) {
	basePath := t.TempDir()
	contents, err := ioutil.ReadFile("../../test/testdata/hello.gz")
	if err != nil {
		t.Fatalf("could not read gzip file")
		return
	}
	zipFilePath := filepath.Join(basePath, "hello.gz")
	if err := ioutil.WriteFile(zipFilePath, contents, os.ModePerm); err != nil {
		t.Fatalf("could not copy gzip file")
		return
	}
	zipff := &gzipFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testGzipFileFetcher{gzipFilePath: zipFilePath},
	}
	unzippedFilePath, err := zipff.FetchFile("http://some.url")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	if unzippedFilePath == "" {
		t.Errorf("expected a path, got empty string")
	}
	if filepath.Base(unzippedFilePath) != "hello" {
		t.Errorf("expected file to named hello, instead got: %s", filepath.Base(unzippedFilePath))
	}
	unzippedContents, err := ioutil.ReadFile(unzippedFilePath)
	if err != nil {
		t.Fatalf("could not read unzipped file contents at %s", unzippedFilePath)
	}
	if "hello" != string(unzippedContents) {
		t.Errorf(`expected unzipped contents to be "hello", got: %s`, string(unzippedContents))
	}
}

func Test_gzipFileFetcher_FetchFile(t *testing.T) {
	basePath := t.TempDir()
	contents, err := ioutil.ReadFile("../../test/testdata/hello.txt")
	if err != nil {
		t.Fatalf("could not read file")
		return
	}
	normFilePath := filepath.Join(basePath, "hello")
	if err := ioutil.WriteFile(normFilePath, contents, os.ModePerm); err != nil {
		t.Fatalf("could not copy file")
		return
	}
	zipff := &gzipFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testGzipFileFetcher{gzipFilePath: normFilePath},
	}
	unzippedFilePath, err := zipff.FetchFile("http://some.url")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	if unzippedFilePath == "" {
		t.Errorf("expected a path, got empty string")
	}
	if filepath.Base(unzippedFilePath) != "hello" {
		t.Errorf("expected file to named hello, instead got: %s", filepath.Base(unzippedFilePath))
	}
	unzippedContents, err := ioutil.ReadFile(unzippedFilePath)
	if err != nil {
		t.Fatalf("could not read file contents at %s", unzippedFilePath)
	}
	if "hello" != string(unzippedContents) {
		t.Errorf(`expected contents to be "hello", got: %s`, string(unzippedContents))
	}
}
