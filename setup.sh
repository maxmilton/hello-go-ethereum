#!/bin/bash
set -Eeuo pipefail

# PROJECT SETUP:
#  1. Download the solc compiler
#  2. Install abigen tool from go-ethereum
#  3. Generate contract ABI, bin, and types files
#
# - https://docs.soliditylang.org/en/latest/installing-solidity.html#static-binaries

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
BIN_DIR="${SCRIPT_DIR}/.bin"
SOLC_ARCH=linux-amd64

msg() { echo >&2 -e "${1-}"; }
die() { msg "$@"; exit 1; }

type curl >/dev/null 2>&1 || die "curl is required but not available!"
type jq >/dev/null 2>&1 || die "jq is required but not available!"
type sha256sum >/dev/null 2>&1 || die "sha256sum is required but not available!"

msg "Getting solc version info..."

solc_info=$(curl -s https://raw.githubusercontent.com/ethereum/solc-bin/gh-pages/${SOLC_ARCH}/list.json) || die "Failed to aquire solc version info!"
solc_latest=$(echo "${solc_info}" | jq '.latestRelease')
solc_build=$(echo "${solc_info}" | jq ".builds[] | select(.version == ${solc_latest})")
solc_version=$(echo "${solc_build}" | jq --raw-output '.path')

test -z "${solc_version}" && die "Failed to extract solc version info!"

if test "${BIN_DIR}/solc" -ef "${BIN_DIR}/${solc_version}"; then
  msg "${solc_version} is already the current version"
else
  solc_sha256=$(echo "${solc_build}" | jq --raw-output '.sha256')

  msg "Downloading ${solc_version}..."

  mkdir -p "${BIN_DIR}" || die "Failed to create .bin directory!"
  curl -L "https://github.com/ethereum/solc-bin/raw/gh-pages/${SOLC_ARCH}/${solc_version}" -o "${BIN_DIR}/${solc_version}" || die "Failed to download!"
  echo "${solc_sha256}  ${BIN_DIR}/${solc_version}" | shasum -c - || die "Failed to verify checksum!"
  chmod +x "${BIN_DIR}/${solc_version}" || die "Failed to make solc executable!"
  rm -f "${BIN_DIR}/solc" || die "Failed to remove old symlink!"
  ln -s "${solc_version}" "${BIN_DIR}/solc" || die "Failed to set new symlink!"
fi


msg "Processing golang dependencies..."

go mod download github.com/ethereum/go-ethereum
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

go generate

msg "Done"
