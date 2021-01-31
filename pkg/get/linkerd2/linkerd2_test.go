package linkerd2

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spachava753/kpkg/pkg/config"
	"testing"
)

func Test_downloadLinkerd2(t *testing.T) {
	hDir, err := homedir.Dir()
	if err != nil {
		t.Fatalf("could not fetch home dir")
	}
	if err := config.CreateBinPath(hDir); err != nil {
		t.Fatalf("could not init dir")
	}
	type args struct {
		version string
		opsys   string
		arch    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Download latest linkerd2 cli for linux/amd64",
			args: args{
				version: "latest",
				opsys:   "linux",
				arch:    "amd64",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Download(tt.args.version, tt.args.opsys, tt.args.arch); (err != nil) != tt.wantErr {
				t.Errorf("downloadLinkerd2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersions(t *testing.T) {
	type args struct {
		installed bool
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "List versions",
			args:    args{installed: false},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Versions(tt.args.installed)
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
