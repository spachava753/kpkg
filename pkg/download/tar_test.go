package download

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testTarFileFetcher struct {
	zipFilePath string
}

func (t testTarFileFetcher) FetchFile(_ string) (string, error) {
	return filepath.Abs(t.zipFilePath)
}

func Test_tarFileFetcher_FetchFile_Zip(t *testing.T) {
	basePath := t.TempDir()
	contents, err := ioutil.ReadFile("../../test/testdata/hello.tar")
	if err != nil {
		t.Fatalf("could not read tar file")
		return
	}
	tarFilePath := filepath.Join(basePath, "hello.tar")
	if err := ioutil.WriteFile(tarFilePath, contents, os.ModePerm); err != nil {
		t.Fatalf("could not copy zip file")
		return
	}
	tarff := &tarFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testTarFileFetcher{zipFilePath: tarFilePath},
	}
	expandedFilePath, err := tarff.FetchFile("http://some.url")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	if expandedFilePath == "" {
		t.Errorf("expected a path, got empty string")
	}

	expandedContents, err := ioutil.ReadFile(
		filepath.Join(
			expandedFilePath, "hello.txt",
		),
	)
	if err != nil {
		t.Fatalf("could not read file contents at %s", expandedFilePath)
	}
	if "hello" != string(expandedContents) {
		t.Errorf(
			`expected contents to be "hello", got: %s`,
			string(expandedContents),
		)
	}
}

func Test_tarFileFetcher_FetchFile(t *testing.T) {
	basePath := t.TempDir()
	contents, err := ioutil.ReadFile("../../test/testdata/hello.txt")
	if err != nil {
		t.Fatalf("could not read file")
		return
	}
	normFilePath := filepath.Join(basePath, "hello")
	if err := ioutil.WriteFile(
		normFilePath, contents, os.ModePerm,
	); err != nil {
		t.Fatalf("could not copy file")
		return
	}
	tarff := &tarFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testTarFileFetcher{zipFilePath: normFilePath},
	}
	filePath, err := tarff.FetchFile("http://some.url")
	if err != nil {
		t.Errorf("expected no error, got: %s", err)
	}
	if filePath == "" {
		t.Errorf("expected a path, got empty string")
	}
	if filepath.Base(filePath) != "hello" {
		t.Errorf(
			"expected file to named hello, instead got: %s",
			filepath.Base(filePath),
		)
	}
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("could not read file contents at %s", filePath)
	}
	if "hello" != string(fileContents) {
		t.Errorf(
			`expected contents to be "hello", got: %s`,
			string(fileContents),
		)
	}
}
