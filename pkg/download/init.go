package download

import (
	"net/http"
	"os"
	"time"
)

func InitFileFetcher() (FileFetcher, error) {
	// create a file fetcher for binaries to fetch file
	fileFetcher, err := MakeFileFetcherTempDir(&http.Client{
		Timeout: time.Second * 10,
	})
	if err != nil {
		return nil, err
	}
	fileFetcher, err = MakeRetryFileFetcher(3, os.Stdout, fileFetcher)
	if err != nil {
		return nil, err
	}
	fileFetcher, err = MakeZipFileFetcher(os.Stdout, fileFetcher)
	if err != nil {
		return nil, err
	}
	fileFetcher, err = MakeGzipFileFetcher(os.Stdout, fileFetcher)
	if err != nil {
		return nil, err
	}
	fileFetcher, err = MakeTarFileFetcher(os.Stdout, fileFetcher)
	if err != nil {
		return nil, err
	}
	return fileFetcher, nil
}
