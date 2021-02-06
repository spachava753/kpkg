package tool

import (
	"github.com/spachava753/kpkg/pkg/config"
	"github.com/spachava753/kpkg/pkg/util"
	"os"
	"path/filepath"
	"testing"
)

func TestLinkedVersion(t *testing.T) {
	type args struct {
		basePath string
		binary   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		setup   func(string) (string, error)
	}{
		{
			name: "Check installed version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "v1.1",
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return root, err
				}

				// make the symlink
				err = os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a"))
				return root, err
			},
		},
		{
			name: "Check installed version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "stable-v1.1",
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "stable-v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return root, err
				}

				// make the symlink
				err = os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a"))
				return root, err
			},
		},
		{
			name: "unknown binary",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "",
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				return root, nil
			},
		},
		{
			name: "broken symlink",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "",
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "stable-v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return root, err
				}

				// make the symlink
				if err := os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a")); err != nil {
					return root, err
				}

				err = os.Remove(binaryFilePath)
				return root, err
			},
		},
		{
			name: "improper base path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "pass in file",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "",
			wantErr: true,
			setup: func(basePath string) (string, error) {
				fPath := filepath.Join(basePath, ".kpkg")
				if _, err := os.Create(fPath); err != nil {
					return "", err
				}
				return fPath, nil
			},
		},
		{
			name: "empty binary name",
			args: args{
				basePath: t.TempDir(),
				binary:   "",
			},
			want:    "",
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				return root, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p string
			if tt.setup != nil {
				var err error
				if p, err = tt.setup(tt.args.basePath); err != nil {
					t.Fatalf("setup func returned an err: %s", err)
					return
				}
			}
			got, err := LinkedVersion(p, tt.args.binary)
			if (err != nil) != tt.wantErr {
				t.Errorf("LinkedVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LinkedVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveVersions(t *testing.T) {
	type args struct {
		basePath string
		binary   string
		versions []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(string) (string, error)
	}{
		{
			name: "happy path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
				},
			},
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return root, err
				}

				// make the symlink
				return root, nil
			},
		},
		{
			name: "happy path multiple versions",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
					"v1.2",
				},
			},
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath1 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath1, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath1 := filepath.Join(binaryPath1, "a")
				_, err = os.Create(binaryFilePath1)
				if err != nil {
					return root, err
				}

				// make a fake binary
				binaryPath2 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath2, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath2 := filepath.Join(binaryPath2, "a")
				_, err = os.Create(binaryFilePath2)
				if err != nil {
					return root, err
				}

				return root, err
			},
		},
		{
			name: "remove installed version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
				},
			},
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return root, err
				}

				// make the symlink
				err = os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a"))
				return root, err
			},
		},
		{
			name: "remove multiple versions with installed version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
					"v1.2",
				},
			},
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath1 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath1, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath1 := filepath.Join(binaryPath1, "a")
				_, err = os.Create(binaryFilePath1)
				if err != nil {
					return root, err
				}

				// make a fake binary
				binaryPath2 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath2, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath2 := filepath.Join(binaryPath2, "a")
				_, err = os.Create(binaryFilePath2)
				if err != nil {
					return root, err
				}

				// make the symlink
				err = os.Symlink(binaryFilePath2, filepath.Join(root, "bin", "a"))
				return root, err
			},
		},
		{
			name: "unknown version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
				},
			},
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make the symlink
				return root, nil
			},
		},
		{
			name: "unknown version with binary path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
				},
			},
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}

				// make the symlink
				return root, nil
			},
		},
		{
			name: "known version without binary file",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{
					"v1.1",
				},
			},
			wantErr: false,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return root, err
				}

				// make the symlink
				return root, nil
			},
		},
		{
			name: "no versions",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: []string{},
			},
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				return root, err
			},
		},
		{
			name: "nil versions",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: nil,
			},
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				return root, err
			},
		},
		{
			name: "bad base path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				versions: nil,
			},
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				return root, err
			},
		},
		{
			name: "missing binary name",
			args: args{
				basePath: t.TempDir(),
				versions: []string{
					"v1.1",
				},
			},
			wantErr: true,
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				return root, err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var p string
			if tt.setup != nil {
				var err error
				if p, err = tt.setup(tt.args.basePath); err != nil {
					t.Fatalf("setup func returned an err: %s", err)
					return
				}
			}
			if err := RemoveVersions(p, tt.args.binary, tt.args.versions); (err != nil) != tt.wantErr {
				t.Errorf("RemoveVersions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstalled(t *testing.T) {
	type args struct {
		basePath string
		binary   string
		version  string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(basePath string) error
		want    bool
		wantErr bool
	}{
		{
			name: "installed",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				version:  "v1.1",
			},
			setup: func(basePath string) error {
				p := filepath.Join(basePath, "a", "v1.1")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				if _, err := os.Create(filepath.Join(p, "a")); err != nil {
					return err
				}
				return nil
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not installed",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				version:  "v1.1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "invalid path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				version:  "v1.1",
			},
			setup: func(basePath string) error {
				if _, err := os.Create(filepath.Join(basePath, "a")); err != nil {
					return err
				}
				return nil
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				version:  "v1.1",
			},
			setup: func(basePath string) error {
				p := filepath.Join(basePath, "a")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				if _, err := os.Create(filepath.Join(p, "v1.1")); err != nil {
					return err
				}
				return nil
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
				version:  "v1.1",
			},
			setup: func(basePath string) error {
				p := filepath.Join(basePath, "a", "v1.1", "a")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				return nil
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := config.CreatePath(tt.args.basePath)
			if err != nil {
				t.Fatalf("could not create .kpkg dir: %s", err)
				return
			}
			if tt.setup != nil {
				if err := tt.setup(root); err != nil {
					t.Fatalf("setup func failed: %s", err)
					return
				}
			}
			got, err := Installed(root, tt.args.binary, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Installed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Installed() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListInstalled(t *testing.T) {
	type args struct {
		basePath string
		binary   string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(basePath string) error
		want    []string
		wantErr bool
	}{
		{
			name: "no version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want: nil,
			setup: func(basePath string) error {
				p := filepath.Join(basePath, "a")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "no version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "one version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want: []string{"v1.1"},
			setup: func(basePath string) error {
				p := filepath.Join(basePath, "a", "v1.1")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				if _, err := os.Create(filepath.Join(p, "a")); err != nil {
					return err
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "multiple versions",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want: []string{"v1.1", "v1.2", "v1.3"},
			setup: func(basePath string) error {
				binaryPath := filepath.Join(basePath, "a")
				{
					versionPath := filepath.Join(binaryPath, "v1.1")
					if err := os.MkdirAll(versionPath, os.ModePerm); err != nil {
						return err
					}
					if _, err := os.Create(filepath.Join(versionPath, "a")); err != nil {
						return err
					}
				}

				{
					versionPath := filepath.Join(binaryPath, "v1.2")
					if err := os.MkdirAll(versionPath, os.ModePerm); err != nil {
						return err
					}
					if _, err := os.Create(filepath.Join(versionPath, "a")); err != nil {
						return err
					}
				}

				{
					versionPath := filepath.Join(binaryPath, "v1.3")
					if err := os.MkdirAll(versionPath, os.ModePerm); err != nil {
						return err
					}
					if _, err := os.Create(filepath.Join(versionPath, "a")); err != nil {
						return err
					}
				}

				return nil
			},
			wantErr: false,
		},
		{
			name: "invalid path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want: nil,
			setup: func(basePath string) error {
				if _, err := os.Create(filepath.Join(basePath, "a")); err != nil {
					return err
				}
				return nil
			},
			wantErr: true,
		},
		{
			name: "invalid path",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want: nil,
			setup: func(basePath string) error {
				p := filepath.Join(basePath, "a", "v1.1", "a")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				return nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := config.CreatePath(tt.args.basePath)
			if err != nil {
				t.Fatalf("could not create .kpkg dir: %s", err)
				return
			}
			if tt.setup != nil {
				if err := tt.setup(root); err != nil {
					t.Fatalf("setup func failed: %s", err)
					return
				}
			}
			got, err := ListInstalled(root, tt.args.binary)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ListInstalled() len(got) = %v, len(want) %v", len(got), len(tt.want))
				return
			}
			for _, v := range got {
				if !util.ContainsString(tt.want, v) {
					t.Errorf("version %s expected, not returned", v)
				}
			}
		})
	}
}

/*func TestInstall(t *testing.T) {
	type args struct {
		version string
		force   bool
	}
	tests := []struct {
		name       string
		args       args
		homeDir    string
		beforeFunc func(basePath string) error
		wantErr    bool
	}{
		{
			name: "Download latest linkerd2 cli for linux/amd64",
			args: args{
				version: "latest",
			},
			homeDir: t.TempDir(),
			wantErr: false,
		},
		{
			name: "Download stable-2.9.2 linkerd2 cli for linux/amd64",
			args: args{
				version: "stable-2.9.2",
			},
			homeDir: t.TempDir(),
			wantErr: false,
		},
		{
			name: "Download edge-21.1.3 linkerd2 cli for linux/amd64",
			args: args{
				version: "edge-21.1.3",
			},
			homeDir: t.TempDir(),
			wantErr: false,
		},
		{
			name: "do not override",
			args: args{
				version: "stable-2.9.2",
			},
			homeDir: t.TempDir(),
			beforeFunc: func(basePath string) error {
				p := filepath.Join(basePath, ".kpkg", "linkerd2", "stable-2.9.2")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				if _, err := os.Create(filepath.Join(p, "linkerd2")); err != nil {
					return err
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "override",
			args: args{
				version: "latest",
				force:   true,
			},
			homeDir: t.TempDir(),
			beforeFunc: func(basePath string) error {
				p := filepath.Join(basePath, ".kpkg", "linkerd2", "stable-2.9.2")
				if err := os.MkdirAll(p, os.ModePerm); err != nil {
					return err
				}
				if _, err := os.Create(filepath.Join(p, "linkerd2")); err != nil {
					return err
				}
				return nil
			},
			wantErr: true,
		},
	}
	// create a file fetcher for binaries to fetch file
	fileFetcher, err := download.MakeFileFetcherTempDir(&http.Client{
		Timeout: time.Second * 10,
	})
	if err != nil {
		t.Fatalf("failed to setup filefetcher: %s", err)
	}
	fileFetcher, err = download.MakeRetryFileFetcher(3, os.Stdout, fileFetcher)
	if err != nil {
		t.Fatalf("failed to setup filefetcher: %s", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeFunc != nil {
				if err := tt.beforeFunc(tt.homeDir); err != nil {
					t.Fatalf("could not proceed with test, setup func error: %s", err)
					return
				}
			}
			root, err := config.CreatePath(tt.homeDir)
			if err != nil {
				t.Fatalf("could not init dir")
			}
			if _, err := Install(filepath.Join(root), tt.args.version, tt.args.force, , runtime.GOOS, runtime.GOARCH, fileFetcher)(); (err != nil) != tt.wantErr {
				t.Errorf("downloadLinkerd2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}*/

func TestPurge(t *testing.T) {
	type args struct {
		basePath string
		binary   string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(basePath string) (string, error)
		wantErr bool
	}{
		{
			name: "purge",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath1 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath1, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath1 := filepath.Join(binaryPath1, "a")
				_, err = os.Create(binaryFilePath1)
				if err != nil {
					return root, err
				}

				// make a fake binary
				binaryPath2 := filepath.Join(root, "a", "v1.2")
				if err := os.MkdirAll(binaryPath2, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath2 := filepath.Join(binaryPath2, "a")
				_, err = os.Create(binaryFilePath2)
				if err != nil {
					return root, err
				}

				// make the symlink
				err = os.Symlink(binaryFilePath2, filepath.Join(root, "bin", "a"))
				return root, err
			},
			wantErr: false,
		},
		{
			name: "purge empty",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}
				return root, err
			},
			wantErr: false,
		},
		{
			name: "purge broken link",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath1 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath1, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath1 := filepath.Join(binaryPath1, "a")
				_, err = os.Create(binaryFilePath1)
				if err != nil {
					return root, err
				}

				// make the symlink
				err = os.Symlink(binaryFilePath1, filepath.Join(root, "bin", "a"))
				if err != nil {
					return root, err
				}
				return root, os.Remove(binaryFilePath1)
			},
			wantErr: false,
		},
		{
			name: "purge no link",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			setup: func(basePath string) (string, error) {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return "", err
				}

				// make a fake binary
				binaryPath1 := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath1, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath1 := filepath.Join(binaryPath1, "a")
				_, err = os.Create(binaryFilePath1)
				if err != nil {
					return root, err
				}

				// make a fake binary
				binaryPath2 := filepath.Join(root, "a", "v1.2")
				if err := os.MkdirAll(binaryPath2, os.ModePerm); err != nil {
					return root, err
				}
				binaryFilePath2 := filepath.Join(binaryPath2, "a")
				_, err = os.Create(binaryFilePath2)
				if err != nil {
					return root, err
				}
				return root, err
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.args.basePath
			if tt.setup != nil {
				var err error
				if p, err = tt.setup(tt.args.basePath); err != nil {
					t.Fatalf("setup func failed: %s", err)
					return
				}
			}
			if err := Purge(p, tt.args.binary); (err != nil) != tt.wantErr {
				t.Errorf("Purge() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
