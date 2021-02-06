package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FileFetcher is an interface responsible for fetching files from a url
type FileFetcher interface {
	// FetchFile takes a url returns to location of the donwloaded file
	FetchFile(url string) (string, error)
}

type basicFileFetcher struct {
	filePath string
	client   *http.Client
}

func (b *basicFileFetcher) FetchFile(url string) (s string, err error) {
	var res *http.Response
	res, err = b.client.Get(url)
	if err != nil {
		return s, err
	}
	if res == nil {
		return "", fmt.Errorf("response object is nil")
	}
	if res.Body == nil {
		return "", fmt.Errorf("response body is nil")
	}

	// defer closing the body, and handle the error if pops up
	defer func() {
		if e := res.Body.Close(); e != nil && err == nil {
			err = e
		}
	}()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("incorrect status for downloading tool: %d", res.StatusCode)
		return
	}

	// create a temporary file
	var f *os.File
	f, err = os.Create(b.filePath)
	if err != nil {
		return
	}

	// defer closing the file, and handle the error if pops up
	defer func() {
		if f != nil {
			if e := f.Close(); e != nil && err == nil {
				err = e
			}
		}
	}()

	if _, err := io.Copy(f, res.Body); err != nil {
		return b.filePath, err
	}

	s = b.filePath
	err = nil
	return
}

func MakeBasicFileFetcher(filePath string, client *http.Client) (FileFetcher, error) {
	if filePath == "" || client == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	return &basicFileFetcher{
		filePath: filePath,
		client:   client,
	}, nil
}

func MakeFileFetcherTempDir(client *http.Client) (FileFetcher, error) {
	if client == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	return MakeBasicFileFetcher(filepath.Join(os.TempDir(), "kpkg_artifact"), client)
}
