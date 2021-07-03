package cmd

import (
	"github.com/spachava753/kpkg/pkg/tool"
	"github.com/spachava753/kpkg/pkg/tool/argocd"
	"github.com/spachava753/kpkg/pkg/tool/buildx"
	"github.com/spachava753/kpkg/pkg/tool/civo"
	"github.com/spachava753/kpkg/pkg/tool/clairctl"
	"github.com/spachava753/kpkg/pkg/tool/copilot"
	"github.com/spachava753/kpkg/pkg/tool/dive"
	"github.com/spachava753/kpkg/pkg/tool/dockercompose"
	"github.com/spachava753/kpkg/pkg/tool/doctl"
	"github.com/spachava753/kpkg/pkg/tool/eksctl"
	"github.com/spachava753/kpkg/pkg/tool/faascli"
	"github.com/spachava753/kpkg/pkg/tool/fzf"
	"github.com/spachava753/kpkg/pkg/tool/gh"
	"github.com/spachava753/kpkg/pkg/tool/golangcilint"
	"github.com/spachava753/kpkg/pkg/tool/goreleaser"
	"github.com/spachava753/kpkg/pkg/tool/helm"
	"github.com/spachava753/kpkg/pkg/tool/helmfile"
	"github.com/spachava753/kpkg/pkg/tool/hugo"
	"github.com/spachava753/kpkg/pkg/tool/inletsctl"
	"github.com/spachava753/kpkg/pkg/tool/istioctl"
	"github.com/spachava753/kpkg/pkg/tool/k3d"
	"github.com/spachava753/kpkg/pkg/tool/k3s"
	"github.com/spachava753/kpkg/pkg/tool/k3sup"
	"github.com/spachava753/kpkg/pkg/tool/k9s"
	"github.com/spachava753/kpkg/pkg/tool/kail"
	"github.com/spachava753/kpkg/pkg/tool/kind"
	"github.com/spachava753/kpkg/pkg/tool/kops"
	"github.com/spachava753/kpkg/pkg/tool/kpkg"
	"github.com/spachava753/kpkg/pkg/tool/krew"
	"github.com/spachava753/kpkg/pkg/tool/kubebench"
	"github.com/spachava753/kpkg/pkg/tool/kubebuilder"
	"github.com/spachava753/kpkg/pkg/tool/kubectl"
	"github.com/spachava753/kpkg/pkg/tool/kubectx"
	"github.com/spachava753/kpkg/pkg/tool/kubens"
	"github.com/spachava753/kpkg/pkg/tool/kubeprompt"
	"github.com/spachava753/kpkg/pkg/tool/kubeseal"
	"github.com/spachava753/kpkg/pkg/tool/kustomize"
	"github.com/spachava753/kpkg/pkg/tool/linkerd2"
	"github.com/spachava753/kpkg/pkg/tool/mc"
	"github.com/spachava753/kpkg/pkg/tool/minikube"
	"github.com/spachava753/kpkg/pkg/tool/nerdctl"
	"github.com/spachava753/kpkg/pkg/tool/opa"
	"github.com/spachava753/kpkg/pkg/tool/osm"
	"github.com/spachava753/kpkg/pkg/tool/pack"
	"github.com/spachava753/kpkg/pkg/tool/packer"
	"github.com/spachava753/kpkg/pkg/tool/polaris"
	"github.com/spachava753/kpkg/pkg/tool/popeye"
	"github.com/spachava753/kpkg/pkg/tool/stern"
	"github.com/spachava753/kpkg/pkg/tool/terraform"
	"github.com/spachava753/kpkg/pkg/tool/terrascan"
	"github.com/spachava753/kpkg/pkg/tool/trivy"
	"github.com/spachava753/kpkg/pkg/tool/vagrant"
	"github.com/spachava753/kpkg/pkg/tool/virtctl"
	"github.com/spachava753/kpkg/pkg/tool/yq"
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
		faascli.MakeBinary(os, arch),
		gh.MakeBinary(os, arch),
		helmfile.MakeBinary(os, arch),
		hugo.MakeBinary(os, arch),
		inletsctl.MakeBinary(os, arch),
		k3sup.MakeBinary(os, arch),
		kail.MakeBinary(os, arch),
		k9s.MakeBinary(os, arch),
		kops.MakeBinary(os, arch),
		krew.MakeBinary(os, arch),
		kubebench.MakeBinary(os, arch),
		kubebuilder.MakeBinary(os, arch),
		kubectx.MakeBinary(os, arch),
		kubens.MakeBinary(os, arch),
		kubeseal.MakeBinary(os, arch),
		kustomize.MakeBinary(os, arch),
		mc.MakeBinary(os, arch),
		minikube.MakeBinary(os, arch),
		osm.MakeBinary(os, arch),
		pack.MakeBinary(os, arch),
		packer.MakeBinary(os, arch),
		popeye.MakeBinary(os, arch),
		stern.MakeBinary(os, arch),
		vagrant.MakeBinary(os, arch),
		yq.MakeBinary(os, arch),
		goreleaser.MakeBinary(os, arch),
		copilot.MakeBinary(os, arch),
		nerdctl.MakeBinary(os, arch),
		argocd.MakeBinary(os, arch),
		trivy.MakeBinary(os, arch),
		golangcilint.MakeBinary(os, arch),
		clairctl.MakeBinary(os, arch),
		terrascan.MakeBinary(os, arch),
		eksctl.MakeBinary(os, arch),
		virtctl.MakeBinary(os, arch),
		dive.MakeBinary(os, arch),
		kpkg.MakeBinary(os, arch),
		kubeprompt.MakeBinary(os, arch),
		fzf.MakeBinary(os, arch),
		polaris.MakeBinary(os, arch),
	}
}
