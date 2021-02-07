package download

import (
	"archive/tar"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type tarFileFetcher struct {
	// print logs to file, if not nil
	out io.Writer
	FileFetcher
}

func (r *tarFileFetcher) FetchFile(u string) (string, error) {
	s, err := r.FileFetcher.FetchFile(u)
	if err != nil {
		return s, err
	}
	if filepath.Ext(s) != ".tar" {
		return s, err
	}

	// create a new folder to store the contents of the tar file in
	now := time.Now()
	randPath := now.Unix() * int64(rand.Uint32())
	p := filepath.Join(filepath.Dir(s), fmt.Sprint(randPath))

	if err := os.MkdirAll(p, 0755); err != nil {
		return s, err
	}

	print(fmt.Sprintf("expanding .tar file %s", s))

	f, err := os.Open(s)
	if err != nil {
		return s, err
	}
	defer f.Close()
	tarReader := tar.NewReader(f)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return p, err
		}

		a := filepath.Join(p, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(a, 0755); err != nil {
				return p, err
			}
		case tar.TypeReg:
			outFile, err := os.Create(a)
			if err != nil {
				return p, err
			}
			outFile.Chmod(os.FileMode(header.Mode))
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return p, err
			}
			if err := outFile.Close(); err != nil {
				return p, err
			}

		default:
			print(fmt.Sprintf(
				"extraction failed: uknown type: %b in %s",
				header.Typeflag,
				header.Name),
			)
		}

	}

	return p, nil
}

func (r *tarFileFetcher) print(message string) error {
	if r.out != nil {
		if _, err := fmt.Fprint(r.out, message); err != nil {
			return err
		}
	}
	return nil
}

func MakeTarFileFetcher(out *os.File, f FileFetcher) (FileFetcher, error) {
	if f == nil {
		return nil, fmt.Errorf("file fetcher param cannot be nil")
	}
	return &tarFileFetcher{
		out:         out,
		FileFetcher: f,
	}, nil
}
