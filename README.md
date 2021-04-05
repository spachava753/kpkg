# kpkg

A tool to install various K8s ecosystem related binaries.

I needed a tool to solve my problem of installing a bunch of different binaries that are either necessary, helpful, or
both while working with the Kubernetes ecosystem. Although some tools can be installed using package managers, many
tools cannot be installed using something like apt, yum, etc. I wanted something that was easy to use, and easy to
remove. All tools are installed in the `$HOME/.kpkg` directory, so all installed tools can be removed by deleting the
folder. I wanted something that could install multiple versions of tools. This is especially useful; for example,
installing the right version of the kubectl cli for your cluster.

# What this tool is not

This tool does not download nor keep track of dependencies. However, this
should not be a problem, as the tools installed usually do not have any dependencies in the first place.

# Installation

`kpkg` can be installed by running `go get github.com/spachava753/kpkg`, or you can run the installation script `curl -sL https://raw.githubusercontent.com/spachava753/kpkg/0.2.0/install.sh | sh`. Optionally, you can just download the zip file from releases: [https://github.com/spachava753/kpkg/releases/latest](https://github.com/spachava753/kpkg/releases/latest)

# Goals

## CLI tool management

- download the latest version of a binary
- download a specific version of a binary
- remove a specific version of binary
- purge all versions of binary
- show installed versions
- show binary installation candidates
- easy to uninstall
- complete parity with [arkade](https://github.com/alexellis/arkade) (meaning all binaries supported by arkade is also
  supported by kpkg)

# UX

The experience of the CLI should look something like this:

For getting a list of all possible binary installs

```bash
kpkg get
```

For installing the latest version

```bash
kpkg get linkerd2
```

For installing a specific version

```bash
kpkg get linkerd2 latest
kpkg get linkerd2 2.9.2
```

You might have multiple versions installed. To set to different version, use the same command

```bash
kpkg get linkerd2 2.9.2
```

To force a re-installation:

```bash
kpkg get linkerd2 2.9.2 --force
```

For listing possible versions of a binary. The output should also show installed versions

```bash
kpkg list linkerd2
```

For listing only installed versions of a binary

```bash
kpkg list linkerd2 --installed
```

For removing a version(s) of a binary. Should fail if current target points to version

```bash
kpkg rm linkerd2 2.9.2
```

For removing all versions of a binary

```bash
kpkg rm linkerd2 --purge
```

# TODO

- [ ] add support for detecting if running on arm{5,6,7}
- [ ] add support for checking checksum
- [ ] add progress bar

# Binary List

By default, the latest version of the binary will be downloaded

Usage:
```bash
  kpkg get [command]
```

Available Commands:

- `argocd`         Declarative continuous deployment for Kubernetes
- `buildx`         Docker CLI plugin for extended build capabilities with BuildKit
- `civo`           Civo CLI is a tool to manage your Civo.com account from the terminal
- `clairctl`       Vulnerability Static Analysis for Containers
- `copilot`        The AWS Copilot CLI is a tool for developers to build, release and operate production ready containerized applications on Amazon ECS and AWS Fargate
- `docker-compose` Define and run multi-container applications with Docker
- `doctl`          The official command line interface for the DigitalOcean API
- `eksctl`         The official CLI for Amazon EKS
- `faas-cli`       openfaas CLI plugin for extended build capabilities with BuildKit
- `gh`             GitHub’s official command line tool
- `golangci-lint`  Fast linters Runner for Go
- `goreleaser`     Deliver Go binaries as fast and easily as possible
- `helm`           The Kubernetes Package Manager
- `helmfile`       Deploy Kubernetes Helm Charts
- `hugo`           The world’s fastest framework for building websites
- `inletsctl`      The fastest way to create self-hosted exit-servers
- `istioctl`       Connect, secure, control, and observe services
- `k3d`            Little helper to run Rancher Lab's k3s in Docker
- `k3s`            Lightweight Kubernetes
- `k3sup`          bootstrap Kubernetes with k3s over SSH < 1 min 🚀
- `k9s`            🐶 Kubernetes CLI To Manage Your Clusters In Style!
- `kail`           kubernetes log viewer
- `kind`           Kubernetes IN Docker - local clusters for testing Kubernetes
- `kops`           Kubernetes Operations (kops) - Production Grade K8s Installation, Upgrades, and Management
- `krew`           📦 Find and install kubectl plugins
- `kube-bench`     Checks whether Kubernetes is deployed according to security best practices as defined in the CIS Kubernetes Benchmark
- `kubebuilder`    SDK for building Kubernetes APIs using CRDs
- `kubectl`        kubectl is a cli to communicate k8s clusters
- `kubectx`        Faster way to switch between clusters in kubectl
- `kubens`         Faster way to switch between namespaces in kubectl
- `kubeseal`       A Kubernetes tool for one-way encrypted Secrets
- `kustomize`      Customization of kubernetes YAML configurations
- `linkerd2`       linkerd2 is a cli to install linkerd2 service mesh
- `mc`             MinIO Client (mc) provides a modern alternative to UNIX commands like ls, cat, cp, mirror, diff, find etc
- `minikube`       Run Kubernetes locally
- `nerdctl`        Docker-compatible CLI for containerd
- `opa`            An open source, general-purpose policy engine
- `osm`            Open Service Mesh (OSM) is a lightweight, extensible, cloud native service mesh that allows users to uniformly manage, secure, and get out-of-the-box observability features for highly dynamic microservice environments
- `pack`           CLI for building apps using Cloud Native Buildpacks
- `packer`         Packer is a tool for creating identical machine images for multiple platforms from a single source configuration
- `popeye`         👀 A Kubernetes cluster resource sanitizer
- `stern`          ⎈ Multi pod and container log tailing for Kubernetes
- `terraform`      Write infrastructure as code using declarative configuration files
- `terrascan`      Detect compliance and security violations across Infrastructure as Code to mitigate risk before provisioning cloud native infrastructure
- `trivy`          A Simple and Comprehensive Vulnerability Scanner for Container Images, Git Repositories and Filesystems. Suitable for CI
- `vagrant`        Vagrant is a tool for building and distributing development environments
- `virtctl`        Kubernetes Virtualization API and runtime in order to define and manage virtual machines
- `yq`             yq is a portable command-line YAML processor
