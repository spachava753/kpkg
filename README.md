# kpkg

A binary to install various K8s ecosystem related binaries. Heavily inspired
by [arkade](https://github.com/alexellis/arkade)

I needed a tool to solve my problem of installing a bunch of different binaries that are either necessary, helpful, or
both while working with the Kubernetes ecosystem. Although some tools can be installed using package managers, many
tools cannot be installed using something like apt, yum, etc. I wanted something that was easy to use, and easy to
remove. All tools are installed in the `$HOME/.kpkg` directory, so all installed tools can be removed by deleting the
folder. I wanted something that could install multiple versions of tools. This is especially useful; for example,
installing the right version of the kubectl cli for your cluster.

# Installation

The project is currently a work in progress. However, if you would like to try it out, it can be installed by running `go get github.com/spachava753/kpkg`. 

# Goals

## CLI tool management

- download the latest version of a binary
- download a specific version of a binary
- remove a specific version of binary
- purge all versions of binary
- show installed versions
- show binary installation candidates
- easy to uninstall

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
kpkg get linkerd2 stable-2.9.2
```

You might have multiple versions installed. To set to different version, use the same command

```bash
kpkg get linkerd2 stable-2.9.2
```

To force a re-installation:

```bash
kpkg get linkerd2 stable-2.9.2 --force
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
kpkg rm linkerd2 stable-2.9.2
kpkg rm linkerd2 stable-2.9.2 edge-21.1.2
kpkg rm linkerd2 stable-2.9.*
```

For removing all versions of a binary

```bash
kpkg rm linkerd2 --purge
```

# TODO

- [ ] filter version list based on arch and os
- [ ] add support for detecting if running on arm{5,6,7}
- [ ] add support for checking checksum

# Binary List

These are the possible binary candidates you can install. Checked ones are complete, unchecked is on the roadmap

- [x] linkerd2
- [x] istioctl
- [x] kubectl
- [x] helm
- [x] kind
- [x] k3s
- [x] k3d
- [x] buildx
- [x] civo
- [ ] docker-compose
- [ ] doctl
- [ ] faas-cli
- [ ] gh
- [ ] helmfile
- [ ] hugo
- [ ] inletsctl
- [ ] k3sup
- [ ] k9s
- [ ] kail
- [ ] kops
- [ ] krew
- [ ] kube-bench
- [ ] kubebuilder
- [ ] kubectx
- [ ] kubens
- [ ] kubeseal
- [ ] kustomize
- [ ] mc
- [ ] minikube
- [ ] opa
- [ ] osm
- [ ] pack
- [ ] packer
- [ ] popeye
- [ ] stern
- [ ] terraform
- [ ] vagrant
- [ ] yq