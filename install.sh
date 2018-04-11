#!/usr/bin/env bash 
set -euf -o pipefail

arch=$(uname -m)
if [[ $arch == x86_64* ]]; then
    cpu="amd64"
elif  [[ $arch == arm* ]]; then
    cpu="arm64"
fi

HELM_ANNOTATE_VERSION=${HELM_ANNOTATE_VERSION:-"0.6"}
dest_dir="${HELM_PLUGIN_DIR:-"$(helm home)/plugins/helm-annotate"}"
file="${dest_dir}/helm-annotate.tar.gz"
os=$(uname -s | tr '[:upper:]' '[:lower:]')
url="https://github.com/Tradeshift/helm-annotate/releases/download/v${HELM_ANNOTATE_VERSION}/helm-annotate_${HELM_ANNOTATE_VERSION}_${os}_${cpu}.tar.gz"

mkdir -p ${dest_dir}

if command -v wget; then
  wget -O "${file}" "${url}"
elif command -v curl; then
  curl -o "${file}" "${url}"
fi

tar zxvf "${file}" -C "${dest_dir}"
chmod +x "${dest_dir}/helm-annotate"
