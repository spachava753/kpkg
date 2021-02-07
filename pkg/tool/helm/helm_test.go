package helm

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestHelmTool_Extract(t *testing.T) {
	version := "v3.5.2"
	artifactPath, err := filepath.Abs("../../../test/testdata/helm")
	if err != nil {
		t.Fatalf("could not get testdata")
	}
	binaryPath, err := filepath.Abs("../../../test/testdata/helm/linux-arm64/helm")
	if err != nil {
		t.Fatalf("could not get binary path")
	}

	l := helmTool{
		arch: runtime.GOARCH,
		os:   runtime.GOOS,
	}
	got, err := l.Extract(artifactPath, version)
	if err != nil {
		t.Errorf("Extract() error = %v, expected nil", err)
	}
	if got != binaryPath {
		t.Errorf("Extract() got = %v, want %v", got, binaryPath)
	}
}
