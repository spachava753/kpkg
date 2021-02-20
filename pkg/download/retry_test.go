package download

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

type testUrlErrFetcher struct{}

type timeoutErr struct{}

func (t timeoutErr) Error() string {
	return "timed out"
}

func (t timeoutErr) Timeout() bool {
	return true
}

func (t *testUrlErrFetcher) FetchFile(_ string) (string, error) {
	return "", &url.Error{
		Op:  "",
		URL: "",
		Err: timeoutErr{},
	}
}

func Test_retryFetcher_FetchFile_ZeroRetries(t *testing.T) {
	r := &retryFetcher{
		retryCount:  0,
		out:         os.Stdout,
		FileFetcher: &testUrlErrFetcher{},
	}
	u, err := r.FetchFile("https://some.url")
	if err == nil {
		t.Errorf("expected err, got nil")
	}
	if u != "" {
		t.Errorf(`expected "", got %s`, u)
	}
}

func Test_retryFetcher_FetchFile_OneRetries(t *testing.T) {
	r := &retryFetcher{
		retryCount:  1,
		out:         os.Stdout,
		FileFetcher: &testUrlErrFetcher{},
	}
	u, err := r.FetchFile("https://some.url")
	if err == nil {
		t.Errorf("expected err, got nil")
	}
	if u != "" {
		t.Errorf(`expected "", got %s`, u)
	}
}

func Test_retryFetcher_FetchFile_TwoRetries(t *testing.T) {
	r := &retryFetcher{
		retryCount:  2,
		out:         os.Stdout,
		FileFetcher: &testUrlErrFetcher{},
	}
	u, err := r.FetchFile("https://some.url")
	if err == nil {
		t.Errorf("expected err, got nil")
	}
	if u != "" {
		t.Errorf(`expected "", got %s`, u)
	}
}

func Test_retryFetcher_FetchFile_ThreeRetries(t *testing.T) {
	r := &retryFetcher{
		retryCount:  3,
		out:         os.Stdout,
		FileFetcher: &testUrlErrFetcher{},
	}
	u, err := r.FetchFile("https://some.url")
	if err == nil {
		t.Errorf("expected err, got nil")
	}
	if u != "" {
		t.Errorf(`expected "", got %s`, u)
	}
}

type testErrFetcher struct {
	path string
}

func (t testErrFetcher) FetchFile(url string) (string, error) {
	c := &http.Client{
		Timeout: time.Second,
	}
	if _, err := c.Get(url); err != nil {
		return "", err
	}

	return t.path, nil
}

func Test_retryFetcher_FetchFile_Retry_Success(t *testing.T) {
	testResp := "hello"
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if count > 1 {
			// ignore error
			_, _ = rw.Write([]byte(testResp))
			return
		}
		time.Sleep(time.Second)
		count++
	}))
	r := &retryFetcher{
		retryCount: 3,
		out:        os.Stdout,
		FileFetcher: &testErrFetcher{
			"some/path",
		},
	}
	u, err := r.FetchFile(server.URL)
	if err != nil {
		t.Errorf("expected nil err, got %s", err)
	}
	if u != "some/path" {
		t.Errorf("expected some/path, got %s", u)
	}
}
