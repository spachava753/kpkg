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
	//contents, err := ioutil.ReadFile("../../test/testdata/hello.tar")
	contents, err := ioutil.ReadFile("../../test/testdata/helm-v3.5.2-linux-arm64.tar")
	if err != nil {
		t.Fatalf("could not read tar file")
		return
	}
	//tarFilePath := filepath.Join(basePath, "hello.tar")
	tarFilePath := filepath.Join(basePath, "helm-v3.5.2-linux-arm64.tar")
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

	expandedContents, err := ioutil.ReadFile(expandedFilePath)
	if err != nil {
		t.Fatalf("could not read file contents at %s", expandedFilePath)
	}
	if "hello" != string(expandedContents) {
		t.Errorf(`expected contents to be "hello", got: %s`, string(expandedContents))
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
	if err := ioutil.WriteFile(normFilePath, contents, os.ModePerm); err != nil {
		t.Fatalf("could not copy file")
		return
	}
	zipff := &zipFileFetcher{
		out:         os.Stdout,
		FileFetcher: &testTarFileFetcher{zipFilePath: normFilePath},
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
