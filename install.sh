#!/bin/sh
# install.sh — zero-dependency installer for the midaz CLI binary.
# Installs only the binary. Skills are installed separately via:
#   midaz skills install --yes
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
#   curl -fsSL ... | sh -s -- --version 0.6.0 --install-dir ~/.local/bin
set -eu

REPO="SparkssL/Midaz-cli"
BINARY="midaz"
INSTALL_DIR="${HOME}/.local/bin"
VERSION=""

usage() {
  cat <<EOF
Usage: install.sh [OPTIONS]

Options:
  --version VERSION   Install a specific version (default: latest)
  --install-dir DIR   Binary install directory (default: ~/.local/bin)
  -h, --help          Show this help
EOF
}

# Parse arguments
while [ $# -gt 0 ]; do
  case "$1" in
    --version)     VERSION="$2"; shift 2 ;;
    --install-dir) INSTALL_DIR="$2"; shift 2 ;;
    -h|--help)     usage; exit 0 ;;
    *)             echo "Unknown option: $1"; usage; exit 1 ;;
  esac
done

detect_platform() {
  OS="$(uname -s)"
  case "$OS" in
    Linux*)  PLATFORM="linux" ;;
    Darwin*) PLATFORM="darwin" ;;
    *)       echo "Unsupported OS: $OS"; exit 1 ;;
  esac

  ARCH="$(uname -m)"
  case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)             echo "Unsupported architecture: $ARCH"; exit 1 ;;
  esac
}

resolve_version() {
  if [ -n "$VERSION" ]; then
    return
  fi
  echo "Fetching latest release..."
  if command -v curl >/dev/null 2>&1; then
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
  elif command -v wget >/dev/null 2>&1; then
    VERSION=$(wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
  else
    echo "Error: curl or wget required"; exit 1
  fi
  if [ -z "$VERSION" ]; then
    echo "Error: could not determine latest version"; exit 1
  fi
}

download() {
  ARCHIVE="${BINARY}-${VERSION}-${PLATFORM}-${ARCH}.tar.gz"
  URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE}"
  CHECKSUMS_URL="https://github.com/${REPO}/releases/download/v${VERSION}/checksums.txt"

  TMPDIR=$(mktemp -d)
  trap 'rm -rf "$TMPDIR"' EXIT

  echo "Downloading ${BINARY} v${VERSION} (${PLATFORM}/${ARCH})..."
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "${TMPDIR}/${ARCHIVE}" "$URL"
    curl -fsSL -o "${TMPDIR}/checksums.txt" "$CHECKSUMS_URL" 2>/dev/null || true
  else
    wget -q -O "${TMPDIR}/${ARCHIVE}" "$URL"
    wget -q -O "${TMPDIR}/checksums.txt" "$CHECKSUMS_URL" 2>/dev/null || true
  fi

  if [ -f "${TMPDIR}/checksums.txt" ]; then
    EXPECTED=$(grep "${ARCHIVE}" "${TMPDIR}/checksums.txt" | awk '{print $1}')
    if [ -n "$EXPECTED" ]; then
      if command -v sha256sum >/dev/null 2>&1; then
        ACTUAL=$(sha256sum "${TMPDIR}/${ARCHIVE}" | awk '{print $1}')
      elif command -v shasum >/dev/null 2>&1; then
        ACTUAL=$(shasum -a 256 "${TMPDIR}/${ARCHIVE}" | awk '{print $1}')
      else
        echo "Warning: sha256sum/shasum not found, skipping checksum verification"
        ACTUAL=""
      fi
      if [ -n "$ACTUAL" ] && [ "$ACTUAL" != "$EXPECTED" ]; then
        echo "Error: checksum mismatch"
        echo "  expected: $EXPECTED"
        echo "  actual:   $ACTUAL"
        exit 1
      fi
      if [ -n "$ACTUAL" ]; then
        echo "Checksum verified."
      fi
    fi
  fi

  tar -xzf "${TMPDIR}/${ARCHIVE}" -C "${TMPDIR}"

  mkdir -p "$INSTALL_DIR"
  cp "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
  chmod +x "${INSTALL_DIR}/${BINARY}"
  echo "Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
}

ensure_path() {
  case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) return ;;
  esac

  echo ""
  echo "${INSTALL_DIR} is not on your PATH."

  SHELL_NAME=$(basename "${SHELL:-/bin/sh}")
  case "$SHELL_NAME" in
    zsh)  PROFILE="${HOME}/.zshrc" ;;
    bash) PROFILE="${HOME}/.bashrc" ;;
    *)    PROFILE="${HOME}/.profile" ;;
  esac

  if [ -f "$PROFILE" ]; then
    echo "export PATH=\"${INSTALL_DIR}:\$PATH\"" >> "$PROFILE"
    echo "Added to ${PROFILE}. Run: source ${PROFILE}"
  else
    echo "Add this to your shell profile:"
    echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
  fi
}

detect_platform
resolve_version
download
ensure_path

echo ""
echo "Done! Run 'midaz version' to verify."
echo ""
echo "To install agent skills (Claude Code, Codex, etc.), run:"
echo "  midaz skills install --yes"
