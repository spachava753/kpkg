package dockercompose

import (
	"runtime"
	"testing"

	"github.com/spachava753/kpkg/test"
)

func TestComposeTool_Versions(t *testing.T) {
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
		t.Run(
			tt.name, func(t *testing.T) {
				l := MakeBinary(runtime.GOOS, runtime.GOARCH)
				_, err := l.Versions(test.TestMaxVersion)
				if (err != nil) != tt.wantErr {
					t.Errorf(
						"Versions() error = %v, wantErr %v", err, tt.wantErr,
					)
					return
				}
			},
		)
	}
}
