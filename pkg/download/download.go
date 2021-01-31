package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// FetchFile downloads a file from a url in a temporary directory
func FetchFile(downloadURL string, f *os.File) error {
	res, err := http.DefaultClient.Get(downloadURL)
	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("incorrect status for downloading tool: %d", res.StatusCode)
	}

	defer res.Body.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return err
	}

	return nil
}
