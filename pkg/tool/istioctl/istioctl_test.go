package istioctl

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIstioctlTool_Versions(t *testing.T) {
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
			l := MakeBinary(runtime.GOOS, runtime.GOARCH)
			_, err := l.Versions()
			if (err != nil) != tt.wantErr {
				t.Errorf("Versions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_istioctlTool_MakeUrl(t *testing.T) {
	type fields struct {
		arch string
		os   string
	}
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "after 1.5",
			fields: fields{
				arch: runtime.GOARCH,
				os:   runtime.GOOS,
			},
			args: args{
				version: "1.6.0",
			},
			want:    fmt.Sprintf("https://github.com/istio/istio/releases/download/1.6.0/istio-1.6.0-%s-%s.tar.gz", runtime.GOOS, runtime.GOARCH),
			wantErr: false,
		},
		{
			name: "before 1.5",
			fields: fields{
				arch: runtime.GOARCH,
				os:   runtime.GOOS,
			},
			args: args{
				version: "1.4.0",
			},
			want:    fmt.Sprintf("https://github.com/istio/istio/releases/download/1.4.0/istio-1.4.0-%s.tar.gz", runtime.GOOS),
			wantErr: false,
		},
		{
			name: "1.5",
			fields: fields{
				arch: runtime.GOARCH,
				os:   runtime.GOOS,
			},
			args: args{
				version: "1.5.0",
			},
			want:    fmt.Sprintf("https://github.com/istio/istio/releases/download/1.5.0/istio-1.5.0-%s.tar.gz", runtime.GOOS),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := istioctlTool{
				arch: tt.fields.arch,
				os:   tt.fields.os,
			}
			got, err := l.MakeUrl(tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MakeUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_istioctlTool_Extract(t *testing.T) {
	type args struct {
		artifactPath string
		version      string
	}
	tests := []struct {
		name    string
		args    args
		setup   func() (string, error)
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				artifactPath: filepath.Join("..", "..", "..", "test", "testdata", "istio"),
				version:      "1.9.0",
			},
			want:    filepath.Join("..", "..", "..", "test", "testdata", "istio", "istio-1.9.0", "bin", "istioctl"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := istioctlTool{os: "linux", arch: "amd64"}
			got, err := l.Extract(tt.args.artifactPath, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Extract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
