package download

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
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

	// create a new folder to store the contents of the zip file in
	now := time.Now()
	randPath := now.Unix() * int64(rand.Uint32())
	p := filepath.Join(filepath.Dir(s), fmt.Sprint(randPath))

	zipReader, err := zip.OpenReader(s)
	if err != nil {
		return s, err
	}
	defer func() {
		_ = zipReader.Close()
	}()

	for _, f := range zipReader.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(p, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(p)+string(os.PathSeparator)) {
			return "", fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return "", err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return "", err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		rc, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		if err := outFile.Close(); err != nil {
			_ = rc.Close()
			return "", err
		}
		if err := rc.Close(); err != nil {
			return "", err
		}

		if err != nil {
			return "", err
		}
	}

	return p, nil
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
