package cmd

import (
	"github.com/spachava753/kpkg/pkg/download"
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
)

// register
func GetTools(basePath string, os, arch string, fileFetcher download.FileFetcher) []tool.Binary {
	return []tool.Binary{
		linkerd2.MakeBinary(basePath, os, arch, fileFetcher),
	}
}
