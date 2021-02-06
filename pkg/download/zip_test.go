package download

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testZipFileFetcher struct {
	zipFilePath string
}

func (t testZipFileFetcher) FetchFile(_ string) (string, error) {
	return filepath.Abs(t.zipFilePath)
}

func Test_zipFileFetcher_FetchFile_Zip(t *testing.T) {
	basePath := t.TempDir()
	contents, err := ioutil.ReadFile("../../test/testdata/hello.zip")
	if err != nil {
		t.Fatalf("could not read zip file")
		return
	}
	zipFilePath := filepath.Join(basePath, "hello.zip")
	if err := ioutil.WriteFile(zipFilePath, contents, os.ModePerm); err != nil {
		t.Fatalf("could not copy zip file")
		return
	}
	zipff := &zipFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testZipFileFetcher{zipFilePath: zipFilePath},
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

func Test_zipFileFetcher_FetchFile(t *testing.T) {
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
	zipff := &zipFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testZipFileFetcher{zipFilePath: normFilePath},
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
