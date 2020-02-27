F#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

readonly root_dir=$(cd `dirname $0`/../../.. && pwd)

readonly version=$(cat ${root_dir}/VERSION)
readonly git_sha=$(git rev-parse HEAD)
readonly git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
readonly slug=${version}-${git_timestamp}-${git_sha:0:16}

${root_dir}/fats/install.sh helm
${root_dir}/fats/install.sh ytt
${root_dir}/fats/install.sh k8s-tag-resolver
${root_dir}/fats/install.sh yq

helm init --client-only
make clean package

# upload releases
gsutil cp -a public-read target/*.yaml gs://projectriff/release/snapshots/${slug}/
