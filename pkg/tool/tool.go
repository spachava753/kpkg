package tool

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/util"
)

// Binary is an interface that all binaries must implement
type Binary interface {
	// Name returns the name of the binary
	Name() string

	// ShortDesc returns a short description of the binary
	ShortDesc() string

	// LongDesc returns a long description of the binary
	LongDesc() string

	// Given a version, it returns a url to fetch the file
	MakeUrl(version string) (string, error)

	// Versions lists the possible installation candidates for a source like Github releases.
	// It gives a sorted slice, based of stability. For example, all of the stable releases will appear
	// before a beta release, and all of the beta releases will appear before the alpha releases. The sorting
	//  is implementation specific.
	Versions(max uint) ([]string, error)

	// Extract takes the downloaded artifacts and does some processing on it to extract the binary.
	// This is useful for binaries like helm, where a tar file is downloaded,
	// containing a LICENSE file, a README file and the helm binary. This method would allow for
	// extracting just the binary before installing. Note that this is called after the artifacts
	// have already been downloaded, so this means that the downloaded file will already by
	// unzipped and/or un-tarred.
	Extract(artifactPath, version string) (string, error)
}

func Install(
	basePath, version string, force, windows bool, max uint, b Binary,
	f download.FileFetcher,
) (s string, err error) {
	binary := b.Name()
	if windows {
		binary = binary + ".exe"
	}

	fmt.Printf("installing %s...\n", binary)

	// check that the version exists
	fmt.Println("verifying version info")
	versions, err := b.Versions(max)
	if err != nil {
		return "", err
	}

	if version != "latest" {
		v, err := semver.NewVersion(version)
		if err != nil {
			return "", err
		}
		version = v.String()
		if !util.ContainsString(versions, version) {
			return "", fmt.Errorf(
				"version %s is not valid for binary %s", version, binary,
			)
		}
	}

	if version == "latest" {
		version = versions[0]
	}

	binaryBasePath := filepath.Join(basePath, binary)
	binaryVersionPath := filepath.Join(binaryBasePath, version)
	binaryPath := filepath.Join(binaryVersionPath, binary)
	binaryLinkPath := filepath.Join(basePath, "bin", binary)

	// check if installed already
	fmt.Println("checking for local installation")
	installed, err := Installed(basePath, binary, version)
	if err != nil {
		return "", err
	}

	if installed {
		// since we already have it installed, set the symlink to this
		fmt.Println("tool already installed!")
		if !force {
			fmt.Println("setting symlink")
			if err = os.Remove(binaryLinkPath); err != nil {
				if !os.IsNotExist(err) {
					return "", fmt.Errorf(
						"could not remove symlink to path %s: %w",
						binaryLinkPath, err,
					)
				}
			}
			return binaryPath, os.Symlink(
				binaryPath, filepath.Join(basePath, "bin", binary),
			)
		}
		// since force is enabled, remove the file and continue
		fmt.Println("removing local installation")
		if err := os.Remove(binaryPath); err != nil {
			return "", err
		}
	}

	// construct the url to fetch the release
	url, err := b.MakeUrl(version)
	if err != nil {
		return "", err
	}

	// download CLI
	fmt.Println("downloading from tool from ", url)
	tmpFilePath, err := f.FetchFile(url)
	if err != nil {
		return "", err
	}
	// cleanup temp file
	defer func() {
		if e := os.Remove(tmpFilePath); e == nil {
			err = e
		}
	}()

	fmt.Println("extracting...")
	tmpFilePath, err = b.Extract(tmpFilePath, version)
	if err != nil {
		return "", err
	}
	if tmpFilePath == "" {
		return "", fmt.Errorf("extraction failed, file path is an emtpy string")
	}

	// copy to our bin path
	fmt.Println("installing...")
	// create binary file
	if _, err := os.Stat(binaryVersionPath); os.IsNotExist(err) {
		if err := os.MkdirAll(binaryVersionPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	// copy the downloaded binary to path
	contents, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(binaryPath, contents, os.ModePerm); err != nil {
		return "", err
	}

	// create symlink to bin path
	if info, err := os.Stat(binaryLinkPath); err == nil {
		if info.IsDir() {
			return "", fmt.Errorf(
				"could not remove symlink, path %s is a dir", binaryLinkPath,
			)
		}
		if err = os.Remove(binaryLinkPath); err != nil {
			return "", fmt.Errorf(
				"could not remove symlink to path %s: %w", binaryLinkPath, err,
			)
		}
	}
	return binaryPath, os.Symlink(
		binaryPath, filepath.Join(basePath, "bin", binary),
	)
}

// RemoveVersions will remove the binary version at the provided path
// basePath is the path where the .kpkg folder is located
// binary is the binary name
// versions is a list of versions to remove
func RemoveVersions(basePath string, binary string, versions []string) error {
	// check that supplied versions are valid
	if len(versions) == 0 {
		return fmt.Errorf("not enough versions were passed in")
	}

	installedVersion, err := LinkedVersion(basePath, binary)
	if err != nil {
		return err
	}

	for _, v := range versions {
		if installedVersion == v {
			return fmt.Errorf(
				"cannot uninstalled version %s, currently in use. Please install another version first",
				v,
			)
		}
		if err := os.RemoveAll(filepath.Join(basePath, binary, v)); err != nil {
			return err
		}
	}
	return nil
}

// Purge will remove all binary versions at the provided path
func Purge(basePath, binary string) error {
	// remove the symlink if exists
	if err := os.Remove(filepath.Join(basePath, "bin", binary)); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.RemoveAll(filepath.Join(basePath, binary)); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}

// LinkedVersion checks if a binary is symlinked, and returns the version symlinked
// It returns an error if the binary is not symlinked
// The symlink we are checking could exist in three states, "existing", "broken", and "not existing"
// "existing" is a happy symlink. Symlink exists, and the link is not broken
// "broken" is a sad symlink. Symlink exists, and the link not broken
// "not existing" is depressed symlink. Symlink does not exist, but it's ok
func LinkedVersion(basePath, binary string) (string, error) {
	// check that a valid binary was passed in
	if binary == "" {
		return "", errors.New("binary name cannot be empty")
	}

	// check that given base path exists
	info, err := os.Stat(basePath)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("given path is not a dir: %s", basePath)
	}

	symPath := filepath.Join(basePath, "bin", binary)

	// symlink doesn't exist
	if _, err := os.Readlink(symPath); err != nil {
		return "", nil
	}
	// returns an err for broken symlink
	linkPath, err := filepath.EvalSymlinks(
		filepath.Join(
			basePath, "bin", binary,
		),
	)
	if err != nil {
		return "", err
	}
	// returns version for happy symlink ;)
	dirs := strings.Split(linkPath, string(os.PathSeparator))
	return dirs[len(dirs)-2], nil
}

// Installed checks if a binary version is already downloaded.
// While walking through the paths, it also checks to make sure that
// the path is valid. If a binary is found at the end of the path,
// it returns true
func Installed(basePath, binary, version string) (bool, error) {
	binaryPath := filepath.Join(basePath, binary)
	binaryVersionPath := filepath.Join(binaryPath, version)
	binaryFilePath := filepath.Join(binaryVersionPath, binary)

	if binaryPathInfo, err := os.Stat(binaryPath); err == nil {
		if !binaryPathInfo.IsDir() {
			return false, fmt.Errorf("path %s contains a file", binaryPath)
		}
		if binaryVersionPathInfo, err := os.Stat(binaryVersionPath); err == nil {
			if !binaryVersionPathInfo.IsDir() {
				return false, fmt.Errorf(
					"path %s contains a file", binaryVersionPath,
				)
			}
			if binaryFilePathInfo, err := os.Stat(binaryFilePath); err == nil {
				if binaryFilePathInfo.IsDir() {
					return false, fmt.Errorf(
						"path %s contains a dir", binaryFilePath,
					)
				}
				return true, nil
			}
		}
	}
	return false, nil
}

// ListToolVersionsInstalled will list all of the downloaded versions of a binary
func ListToolVersionsInstalled(basePath, binary string) ([]string, error) {
	var installedVersions []string

	binaryPath := filepath.Join(basePath, binary)
	binaryPathInfo, err := os.Stat(binaryPath)
	if err == nil {
		if !binaryPathInfo.IsDir() {
			return installedVersions, fmt.Errorf(
				"path %s contains a file", binaryPath,
			)
		}
	}
	if err != nil {
		if os.IsNotExist(err) {
			return installedVersions, nil
		}
		return installedVersions, err
	}

	versions, err := ioutil.ReadDir(binaryPath)
	if err != nil {
		return installedVersions, err
	}

	for _, v := range versions {
		if !v.IsDir() {
			continue
		}
		i, err := Installed(basePath, binary, v.Name())
		if err != nil {
			return installedVersions, err
		}
		if i {
			installedVersions = append(installedVersions, v.Name())
		}
	}

	return installedVersions, nil
}

// ListInstalled returns the list of installed binaries
func ListInstalled(basePath string) ([]string, error) {
	var installedBinaries []string

	basePath = filepath.Join(basePath, "bin")

	binaryPathInfo, err := os.Stat(basePath)
	if err == nil {
		if !binaryPathInfo.IsDir() {
			return installedBinaries, fmt.Errorf(
				"path %s contains a file", basePath,
			)
		}
	}
	if err != nil {
		if os.IsNotExist(err) {
			return installedBinaries, nil
		}
		return installedBinaries, err
	}

	binaries, err := ioutil.ReadDir(basePath)
	if err != nil {
		return installedBinaries, err
	}

	for _, v := range binaries {
		if v.IsDir() {
			continue
		}
		installedBinaries = append(installedBinaries, v.Name())
	}

	return installedBinaries, nil
}
