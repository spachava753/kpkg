package linkerd2

import (
	"reflect"
	"runtime"
	"testing"
)

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
			l := MakeBinary(runtime.GOOS, runtime.GOARCH)
			_, err := l.Versions()
			if (err != nil) != tt.wantErr {
				t.Errorf("Versions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_sortVersions(t *testing.T) {
	type args struct {
		versions []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				versions: []string{
					"stable-2.9.0",
					"stable-2.9.2",
					"stable-2.9.1",
				},
			},
			want: []string{
				"stable-2.9.2",
				"stable-2.9.1",
				"stable-2.9.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortVersions(tt.args.versions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
