package faascli

import (
	"runtime"
	"testing"
)

func TestFaasCliTool_Versions(t *testing.T) {
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
