#!/bin/sh

# Modified Linkerd2 bash install script: curl -sL run.linkerd.io/install

set -eu

KPKG_VERSION=${KPKG_VERSION:-0.4.1}
KPKG_ROOT=${KPKG_ROOT:-"${HOME}/.kpkg"}

happyexit() {
  echo ""
  echo "Add kpkg CLI and kpkg installed binaries to your path with:"
  echo ""
  echo "  export PATH=\$PATH:${KPKG_ROOT}/bin:${KPKG_ROOT}/original"
  echo ""
  echo "Now run:"
  echo ""
  echo "  kpkg get        # validate that CLI was installed"
  echo ""
  echo "For more info: Visit https://github.com/spachava753/kpkg"
  echo ""
  exit 0
}

OS=$(uname -s)
arch=$(uname -m)
cli_arch=""
case $OS in
  CYGWIN* | MINGW64*)
    ;;
  Darwin)
    case $arch in
      x86_64)
        cli_arch=amd64
        ;;
      armv8*)
        cli_arch=arm64
        ;;
      aarch64*)
        cli_arch=arm64
        ;;
      amd64|arm64)
        cli_arch=$arch
        ;;
      *)
        echo "There is no kpkg $OS support for $arch. Please open an issue with your platform details."
        exit 1
        ;;
    esac
    ;;
  Linux)
    case $arch in
      x86_64)
        cli_arch=amd64
        ;;
      armv8*)
        cli_arch=arm64
        ;;
      aarch64*)
        cli_arch=arm64
        ;;
      armv*)
        cli_arch=arm
        ;;
      amd64|arm64)
        cli_arch=$arch
        ;;
      *)
        echo "There is no kpkg $OS support for $arch. Please open an issue with your platform details."
        exit 1
        ;;
    esac
    ;;
  *)
    echo "There is no kpkg support for $OS/$arch. Please open an issue with your platform details."
    exit 1
    ;;
esac
OS=$(echo $OS | tr '[:upper:]' '[:lower:]')

tmpdir=$(mktemp -d /tmp/kpkg.XXXXXX)
srcfile="kpkg_${OS}_${cli_arch}.zip"
dstfile="${KPKG_ROOT}/original/kpkg"
url="https://github.com/spachava753/kpkg/releases/download/${KPKG_VERSION}/${srcfile}"

(
  cd "$tmpdir"

  echo "Downloading ${srcfile} at ${url}..."
  curl -fLO "${url}"
  echo "Download complete!"
  echo ""
  echo "Unzipping... 📤"
  unzip "${srcfile}"
  echo "Unzipped! ☑"
  echo ""
)

srcfile="kpkg"

(
  echo "Installing... 💣"
  mkdir -p "${KPKG_ROOT}/original"
  rm -f "${dstfile}"
  mv "${tmpdir}/${srcfile}" "${dstfile}"
  chmod +x "${dstfile}"
  echo "Installed! 💥"
)


rm -r "$tmpdir"

(
  echo "exporting PATH"
  export PATH=${PATH}:${KPKG_ROOT}/bin:${KPKG_ROOT}/original
)

echo "kpkg was successfully installed 🎉"
echo ""
happyexit
