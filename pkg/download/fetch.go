package download

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
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

func (b *basicFileFetcher) FetchFile(urlStr string) (s string, err error) {
	var res *http.Response
	res, err = b.client.Get(urlStr)
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
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	fLoc := filepath.Join(b.filePath, filepath.Base(parsedUrl.Path))
	var f *os.File
	f, err = os.Create(fLoc)
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
		return fLoc, err
	}

	s = fLoc
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
	randPath := time.Now().Unix() * int64(rand.Uint32())
	p := filepath.Join(os.TempDir(), "kpkg", fmt.Sprint(randPath))
	if err := os.MkdirAll(p, 0755); err != nil {
		return nil, err
	}
	return MakeBasicFileFetcher(p, client)
}
