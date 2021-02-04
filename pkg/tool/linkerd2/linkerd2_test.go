package linkerd2

import (
	"github.com/spachava753/kpkg/pkg/config"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLinkerd2Tool_Install(t *testing.T) {
	type args struct {
		version string
		force   bool
	}
	tests := []struct {
		name       string
		args       args
		homeDir    string
		beforeFunc func() error
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
				version: "latest",
			},
			homeDir: t.TempDir(),
			beforeFunc: func() error {
				p := filepath.Join(t.TempDir(), ".kpkg", "linkerd2", "stable-2.9.2")
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
		{
			name: "override",
			args: args{
				version: "latest",
				force:   true,
			},
			homeDir: t.TempDir(),
			beforeFunc: func() error {
				p := filepath.Join(t.TempDir(), ".kpkg", "linkerd2", "stable-2.9.2")
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.beforeFunc != nil {
				if err := tt.beforeFunc(); err != nil {
					t.Fatalf("could not proceed with test, setup func error: %s", err)
					return
				}
			}
			root, err := config.CreatePath(tt.homeDir)
			if err != nil {
				t.Fatalf("could not init dir")
			}
			l := MakeBinary(filepath.Join(root), runtime.GOOS, runtime.GOARCH)
			if _, err := l.Install(tt.args.version, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("downloadLinkerd2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkerd2Tool_Versions(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		homeDir string
		wantErr bool
	}{
		{
			name:    "List versions",
			want:    nil,
			homeDir: t.TempDir(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := config.CreatePath(tt.homeDir)
			if err != nil {
				t.Fatalf("could not init dir")
			}
			l := MakeBinary(root, runtime.GOOS, runtime.GOARCH)
			_, err = l.Versions()
			if (err != nil) != tt.wantErr {
				t.Errorf("Versions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Versions() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
