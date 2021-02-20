package download

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func Test_basicFileFetcher_FetchFile(t *testing.T) {
	testResp := "hello"
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		// ignore error
		_, _ = rw.Write([]byte(testResp))
	}))
	fPath := filepath.Join(t.TempDir(), "testResp")
	defer server.Close()

	b := &basicFileFetcher{fPath, server.Client()}
	tmpFilePath, err := b.FetchFile(server.URL)
	if err != nil {
		t.Errorf("encountered error when fetching file: %s", err)
	}
	if tmpFilePath == "" {
		t.Errorf("returned path was empty")
	}
	contents, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		t.Errorf("failed to read file at %s: %s", tmpFilePath, err)
	}
	if testResp != string(contents) {
		t.Errorf("wanted: %s; got: %s", testResp, string(contents))
	}
}
