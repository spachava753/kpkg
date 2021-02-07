package cmd

import (
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spachava753/kpkg/pkg/tool/helm"
	"github.com/spachava753/kpkg/pkg/tool/kind"
	"github.com/spachava753/kpkg/pkg/tool/kubectl"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
)

// register tools here
func GetTools(basePath string, os, arch string) []tool.Binary {
	return []tool.Binary{
		linkerd2.MakeBinary(basePath, os, arch),
		kubectl.MakeBinary(basePath, os, arch),
		kind.MakeBinary(basePath, os, arch),
		helm.MakeBinary(basePath, os, arch),
	}
}
