package config

import (
	"errors"
	"os"
	"path"
	"testing"
	"time"
)

func TestCreateBinPath(t *testing.T) {
	type args struct {
		basePath string
	}
	tests := []struct {
		name                  string
		args                  args
		generateBaseDirBefore bool
		wantErr               bool
	}{
		{
			name:    "Create root folder",
			args:    args{basePath: t.TempDir()},
			wantErr: false,
		},
		{
			name:                  "Skip creating root folder",
			args:                  args{basePath: t.TempDir()},
			generateBaseDirBefore: true,
			wantErr:               false,
		},
		{
			name:    "invalid base path",
			args:    args{basePath: "/fdsfa/fasdfa"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var modTime time.Time
			rootDirPath := path.Join(tt.args.basePath, ".kpkg")
			if tt.generateBaseDirBefore {
				if err := os.Mkdir(rootDirPath, os.ModePerm); err != nil {
					t.Fatalf("could not create dir before hand: %s", err)
					return
				}
				info, _ := os.Stat(rootDirPath)
				modTime = info.ModTime()
			}
			if _, err := CreatePath(tt.args.basePath); (err != nil) != tt.wantErr {
				t.Errorf("CreatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				a, err := os.Stat(rootDirPath)
				if err != nil {
					var pathError *os.PathError
					if errors.As(err, &pathError) {
						t.Errorf("root directory was not created")
					}
					t.Fatalf("could not stat dir: %s", err)
				}
				if !a.IsDir() {
					t.Errorf("created file instead of dir")
				}
			}
			if info, _ := os.Stat(rootDirPath); tt.generateBaseDirBefore && modTime != info.ModTime() {
				t.Errorf("overrode root dir")
			}
		})
	}
}
