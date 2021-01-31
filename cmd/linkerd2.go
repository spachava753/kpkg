package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/get/linkerd2"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"runtime"
)

func MakeLinkerd2() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "linkerd2",
		Short: "linkerd2 is a well-known service mesh",
		Long:  `linkerd2 is a well-known service mesh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := cmd.Flags().GetString(CliVersionFlag)
			if err != nil {
				return err
			}
			return downloadLinkerd2(v, runtime.GOOS, runtime.GOARCH)
		},
	}
	return cmd
}

func downloadLinkerd2(version, opsys, arch string) error {
	urlConstructor := linkerd2.MakeUrlConstructor(version, opsys, arch)
	url, err := urlConstructor.Construct()
	if err != nil {
		return err
	}

	// create a temp file to download the CLI to
	tmpF, err := ioutil.TempFile(os.TempDir(), "linkerd2")
	defer tmpF.Close()
	defer os.Remove(tmpF.Name())
	if err != nil {
		return err
	}

	// download CLI
	err = download.FetchFile(url, tmpF)
	if err != nil {
		return err
	}

	// copy to our bin path
	hDir, err := homedir.Dir()
	if err != nil {
		return err
	}

	// create binary file
	path := hDir + "/.kpkg/linkerd2/" + version
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	// copy the downloaded binary to path
	contents, err := ioutil.ReadFile(tmpF.Name())
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path+"/linkerd2", contents, os.ModePerm); err != nil {
		return err
	}

	return nil
}
