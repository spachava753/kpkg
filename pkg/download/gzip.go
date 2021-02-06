package download

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type gzipFileFetcher struct {
	// print logs to file, if not nil
	out io.Writer
	FileFetcher
}

func (r *gzipFileFetcher) FetchFile(u string) (string, error) {
	s, err := r.FileFetcher.FetchFile(u)
	if err != nil {
		return s, err
	}
	if filepath.Ext(s) != ".gz" {
		return s, err
	}

	print(fmt.Sprintf("unzipping .gz file %s", s))

	f, err := os.Open(s)
	if err != nil {
		return s, err
	}
	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		return s, err
	}
	defer gzipReader.Close()

	fPath := s[:len(s)-3]

	contents, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return s, err
	}

	if err := ioutil.WriteFile(fPath, contents, os.ModePerm); err != nil {
		return s, err
	}

	return fPath, nil
}

func (r *gzipFileFetcher) print(message string) error {
	if r.out != nil {
		if _, err := fmt.Fprint(r.out, message); err != nil {
			return err
		}
	}
	return nil
}

func MakeGzipFileFetcher(out *os.File, f FileFetcher) (FileFetcher, error) {
	if f == nil {
		return nil, fmt.Errorf("file fetcher param cannot be nil")
	}
	return &gzipFileFetcher{
		out:         out,
		FileFetcher: f,
	}, nil
}
