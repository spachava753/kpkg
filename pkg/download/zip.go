package download

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type zipFileFetcher struct {
	// print logs to file, if not nil
	out io.Writer
	FileFetcher
}

func (r *zipFileFetcher) FetchFile(u string) (string, error) {
	s, err := r.FileFetcher.FetchFile(u)
	if err != nil {
		return s, err
	}
	if filepath.Ext(s) != ".zip" {
		return s, err
	}

	print(fmt.Sprintf("unzipping .zip file %s\n", s))

	zipReader, err := zip.OpenReader(s)
	if err != nil {
		return s, err
	}
	defer zipReader.Close()

	if len(zipReader.File) != 1 {
		return s, fmt.Errorf("cannot handle unzipping archive file at %s", s)
	}

	f := zipReader.File[0]
	bPath := s[:len(s)-4]

	rc, err := f.Open()
	if err != nil {
		return s, err
	}
	defer rc.Close()

	contents, err := ioutil.ReadAll(rc)
	if err != nil {
		return s, err
	}

	if err := ioutil.WriteFile(bPath, contents, os.ModePerm); err != nil {
		return s, err
	}

	return bPath, nil
}

func (r *zipFileFetcher) print(message string) error {
	if r.out != nil {
		if _, err := fmt.Fprint(r.out, message); err != nil {
			return err
		}
	}
	return nil
}

func MakeZipFileFetcher(out *os.File, f FileFetcher) (FileFetcher, error) {
	if f == nil {
		return nil, fmt.Errorf("file fetcher param cannot be nil")
	}
	return &zipFileFetcher{
		out:         out,
		FileFetcher: f,
	}, nil
}
