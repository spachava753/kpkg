package download

import (
	"io/ioutil"
	"testing"
)

func TestFetchFile(t *testing.T) {
	type args struct {
		downloadURL string
		fileName    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "Helm test",
			args: args{
				downloadURL: "https://run.linkerd.io/install",
				fileName:    "linkerd2-installer",
			},
			wantErr: false,
			want: `#!/bin/sh

set -eu

LINKERD2_VERSION=${LINKERD2_VERSION:-stable-2.9.2}
INSTALLROOT=${INSTALLROOT:-"${HOME}/.linkerd2"}

happyexit() {
  echo ""
  echo "Add the linkerd CLI to your path with:"
  echo ""
  echo "  export PATH=\$PATH:${INSTALLROOT}/bin"
  echo ""
  echo "Now run:"
  echo ""
  echo "  linkerd check --pre                     # validate that Linkerd can be installed"
  echo "  linkerd install | kubectl apply -f -    # install the control plane into the 'linkerd' namespace"
  echo "  linkerd check                           # validate everything worked!"
  echo "  linkerd dashboard                       # launch the dashboard"
  echo ""
  echo "Looking for more? Visit https://linkerd.io/2/next-steps"
  echo ""
  exit 0
}

validate_checksum() {
  filename=$1
  SHA=$(curl -sfL "${url}.sha256")
  echo ""
  echo "Validating checksum..."

  case $checksumbin in
    *openssl)
      checksum=$($checksumbin dgst -sha256 "${filename}" | sed -e 's/^.* //')
      ;;
    *shasum)
      checksum=$($checksumbin -a256 "${filename}" | sed -e 's/^.* //')
      ;;
  esac

  if [ "$checksum" != "$SHA" ]; then
    echo "Checksum validation failed." >&2
    return 1
  fi
  echo "Checksum valid."
  return 0
}

OS=$(uname -s)
arch=$(uname -m)
cli_arch=""
case $OS in
  CYGWIN* | MINGW64*)
    OS=windows.exe
    ;;
  Darwin)
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
        echo "There is no linkerd $OS support for $arch. Please open an issue with your platform details."
        exit 1
        ;;
    esac
    ;;
  *)
    echo "There is no linkerd support for $OS/$arch. Please open an issue with your platform details."
    exit 1
    ;;
esac
OS=$(echo $OS | tr '[:upper:]' '[:lower:]')

checksumbin=$(command -v openssl) || checksumbin=$(command -v shasum) || {
  echo "Failed to find checksum binary. Please install openssl or shasum."
  exit 1
}


tmpdir=$(mktemp -d /tmp/linkerd2.XXXXXX)
srcfile="linkerd2-cli-${LINKERD2_VERSION}-${OS}"
if [ -n "${cli_arch}" ]; then
  srcfile="${srcfile}-${cli_arch}"
fi
dstfile="${INSTALLROOT}/bin/linkerd-${LINKERD2_VERSION}"
url="https://github.com/linkerd/linkerd2/releases/download/${LINKERD2_VERSION}/${srcfile}"

if [ -e "${dstfile}" ]; then
  if validate_checksum "${dstfile}"; then
    echo ""
    echo "Linkerd ${LINKERD2_VERSION} was already downloaded; making it the default 🎉"
    echo ""
    echo "To force re-downloading, delete '${dstfile}' then run me again."
    (
      rm -f "${INSTALLROOT}/bin/linkerd"
      ln -s "${dstfile}" "${INSTALLROOT}/bin/linkerd"
    )
    happyexit
  fi
fi

(
  cd "$tmpdir"

  echo "Downloading ${srcfile}..."
  curl -fLO "${url}"
  echo "Download complete!"

  if ! validate_checksum "${srcfile}"; then
    exit 1
  fi
  echo ""
)

(
  mkdir -p "${INSTALLROOT}/bin"
  mv "${tmpdir}/${srcfile}" "${dstfile}"
  chmod +x "${dstfile}"
  rm -f "${INSTALLROOT}/bin/linkerd"
  ln -s "${dstfile}" "${INSTALLROOT}/bin/linkerd"
)


rm -r "$tmpdir"

echo "Linkerd ${LINKERD2_VERSION} was successfully installed 🎉"
echo ""
happyexit
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create temp file
			f, err := ioutil.TempFile(t.TempDir(), tt.args.fileName)
			if err != nil {
				t.Fatalf("could not create file %s", tt.args.fileName)
			}
			// add clean up for temp file
			defer f.Close()
			if err := FetchFile(tt.args.downloadURL, f); (err != nil) != tt.wantErr {
				t.Errorf("FetchFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			downloadedContents, err := ioutil.ReadFile(f.Name())
			if err != nil {
				t.Fatalf("could not read from file created temp file: %s", err)
			}
			if tt.want != string(downloadedContents) {
				t.Errorf("FetchFile() result = %s, want %s", string(downloadedContents), tt.want)
			}
		})
	}
}
