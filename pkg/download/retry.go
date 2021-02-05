package download

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
)

type retryFetcher struct {
	// number of times to retry
	retryCount uint
	// print logs to file, if not nil
	out io.Writer
	FileFetcher
}

func (r *retryFetcher) FetchFile(u string) (string, error) {
	s, err := r.FileFetcher.FetchFile(u)
	var count uint = 1
	var urlErr *url.Error
	for err != nil && errors.As(err, &urlErr) && count <= r.retryCount {
		if e := r.print(fmt.Sprintf("fetching file from url %s failed, retrying count: %d\n", u, count)); e != nil {
			return s, e
		}
		s, err = r.FileFetcher.FetchFile(u)
		count++
	}
	return s, err
}

func (r *retryFetcher) print(message string) error {
	if r.out != nil {
		if _, err := fmt.Fprint(r.out, message); err != nil {
			return err
		}
	}
	return nil
}

func MakeRetryFileFetcher(retry uint, out *os.File, f FileFetcher) (FileFetcher, error) {
	if f == nil {
		return nil, fmt.Errorf("file fetcher param cannot be nil")
	}
	return &retryFetcher{
		retryCount:  retry,
		out:         out,
		FileFetcher: f,
	}, nil
}
