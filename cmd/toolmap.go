package cmd

import (
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spachava753/kpkg/pkg/tool/buildx"
	"github.com/spachava753/kpkg/pkg/tool/civo"
	"github.com/spachava753/kpkg/pkg/tool/dockercompose"
	"github.com/spachava753/kpkg/pkg/tool/doctl"
	"github.com/spachava753/kpkg/pkg/tool/helm"
	"github.com/spachava753/kpkg/pkg/tool/istioctl"
	"github.com/spachava753/kpkg/pkg/tool/k3d"
	"github.com/spachava753/kpkg/pkg/tool/k3s"
	"github.com/spachava753/kpkg/pkg/tool/kind"
	"github.com/spachava753/kpkg/pkg/tool/kubectl"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
	"github.com/spachava753/kpkg/pkg/tool/opa"
	"github.com/spachava753/kpkg/pkg/tool/terraform"
)

// register tools here
func GetTools(os, arch string) []tool.Binary {
	return []tool.Binary{
		linkerd2.MakeBinary(os, arch),
		kubectl.MakeBinary(os, arch),
		kind.MakeBinary(os, arch),
		helm.MakeBinary(os, arch),
		istioctl.MakeBinary(os, arch),
		k3s.MakeBinary(os, arch),
		k3d.MakeBinary(os, arch),
		buildx.MakeBinary(os, arch),
		civo.MakeBinary(os, arch),
		dockercompose.MakeBinary(os, arch),
		opa.MakeBinary(os, arch),
		terraform.MakeBinary(os, arch),
		doctl.MakeBinary(os, arch),
	}
}
