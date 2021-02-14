package download

import (
	"archive/tar"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type tarFileFetcher struct {
	// print logs to file, if not nil
	out io.Writer
	FileFetcher
}

func (r *tarFileFetcher) FetchFile(u string) (string, error) {
	s, err := r.FileFetcher.FetchFile(u)
	if err != nil || filepath.Ext(s) != ".tar" {
		return s, err
	}

	// create a new folder to store the contents of the tar file in
	now := time.Now()
	randPath := now.Unix() * int64(rand.Uint32())
	p := filepath.Join(filepath.Dir(s), fmt.Sprint(randPath))

	if err := os.MkdirAll(p, 0755); err != nil {
		return s, err
	}

	print(fmt.Sprintf("expanding .tar file %s\n", s))

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

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(a, filepath.Clean(p)+string(os.PathSeparator)) {
			return "", fmt.Errorf("%s: illegal file path", a)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(a, 0755); err != nil {
				return p, err
			}
		case tar.TypeReg:
			outFile, err := os.Create(a)
			if err != nil {
				if !os.IsNotExist(err) {
					return p, err
				}
				// probably means that the a parent folder is not created while expanding the archive
				// example: the github cli tar releases
				b := filepath.Dir(a)
				if _, err := os.Stat(b); err != nil && os.IsNotExist(err) {
					if err := os.MkdirAll(b, 0755); err != nil {
						return "", err
					}
					// try creating the file again
					outFile, err = os.Create(a)
					if err != nil {
						return p, err
					}
				}
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
