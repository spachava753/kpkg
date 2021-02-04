package tool

import (
	"github.com/spachava753/kpkg/pkg/config"
	"os"
	"path/filepath"
	"testing"
)

func TestInstalledVersion(t *testing.T) {
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
			name: "missing symlink",
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
			got, err := InstalledVersion(p, tt.args.binary)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstalledVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InstalledVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
