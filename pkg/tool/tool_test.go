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
		setup   func(string) error
	}{
		{
			name: "Check installed version",
			args: args{
				basePath: t.TempDir(),
				binary:   "a",
			},
			want:    "v1.1",
			wantErr: false,
			setup: func(basePath string) error {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return err
				}

				// make the symlink
				return os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a"))
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
			setup: func(basePath string) error {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "stable-v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return err
				}

				// make the symlink
				return os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a"))
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
			setup: func(basePath string) error {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "stable-v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return err
				}

				return nil
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
			setup: func(basePath string) error {
				root, err := config.CreatePath(basePath)
				if err != nil {
					return err
				}

				// make a fake binary
				binaryPath := filepath.Join(root, "a", "stable-v1.1")
				if err := os.MkdirAll(binaryPath, os.ModePerm); err != nil {
					return err
				}
				binaryFilePath := filepath.Join(binaryPath, "a")
				_, err = os.Create(binaryFilePath)
				if err != nil {
					return err
				}

				// make the symlink
				if err := os.Symlink(binaryFilePath, filepath.Join(root, "bin", "a")); err != nil {
					return err
				}

				return os.Remove(binaryFilePath)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				if err := tt.setup(tt.args.basePath); err != nil {
					t.Fatalf("setup func returned an err: %s", err)
					return
				}
			}
			got, err := InstalledVersion(filepath.Join(tt.args.basePath, ".kpkg"), tt.args.binary)
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
