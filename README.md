# kpkg
A binary to install various K8s ecosystem related binaries. Heavily inspired by arkade

I really like [arkade](https://github.com/alexellis/arkade). It solved my problem of installing a bunch of different binaries that are either necessary, helpful, or both while working with the Kubernetes ecosystem. However, it was missing a couple of features that I would have loved to have. This is my solution to the problem that arkade solves.

# Goals

## CLI tool management
 - download the latest version of binary
 - download a specific version of binary
 - remove a specific version of binary
 - purge all versions of binary
 - show installed versions
 - show binary installation candidates

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
kpkg get linkerd2 --version 21.1.4
kpkg get linkerd2 --version=21.1.4
kpkg get linkerd2 --version "21.1.4"
```

You might have multiple versions installed. To set to different version, use the same command
```bash
kpkg get linkerd2 --version 21.1.4
```

To force a re-installation:
```bash
kpkg get linkerd2 --version 21.1.4 --force
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
kpkg rm linkerd2 --version 21.1.4
kpkg rm linkerd2 --version "21.1.4,21.1.3"
kpkg rm linkerd2 --version 21.1.4 --version 21.1.3
kpkg rm linkerd2 --version 21.1.*
```

For removing all versions of a binary
```bash
kpkg purge linkerd2
```

# Binary List
These are the possible binary candidates you can install. Checked ones are complete, unchecked is on the roadmap

- [ ] linkerd2
- [ ] kubectl
- [ ] helm
- [ ] kind
- [ ] k3s
- [ ] k3d
- [ ] k3sup
- [ ] k0s
- [ ] istioctl
- [ ] kops
- [ ] eksctl
- [ ] doctl
- [ ] minikube
- [ ] terraform
- [ ] vagrant
- [ ] packer