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
	ext := filepath.Ext(s)
	if ext != ".gz" && ext != ".tgz" {
		return s, err
	}

	print(fmt.Sprintf("decompressing .gz file %s\n", s))

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
	if ext == ".tgz" {
		// remove one more character
		fPath = fPath[:len(fPath)-1]
		fPath = fPath + ".tar"
	}

	contents, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return s, err
	}

	if err := ioutil.WriteFile(fPath, contents, os.ModePerm); err != nil {
		return s, err
	}

	return fPath, nil
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
