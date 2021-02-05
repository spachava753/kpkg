package tool

import (
	"github.com/spachava753/kpkg/pkg/config"
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
