package cmd

import (
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spachava753/kpkg/pkg/tool/helm"
	"github.com/spachava753/kpkg/pkg/tool/kind"
	"github.com/spachava753/kpkg/pkg/tool/kubectl"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
)

// register tools here
func GetTools(os, arch string) []tool.Binary {
	return []tool.Binary{
		linkerd2.MakeBinary(os, arch),
		kubectl.MakeBinary(os, arch),
		kind.MakeBinary(os, arch),
		helm.MakeBinary(os, arch),
	}
}
