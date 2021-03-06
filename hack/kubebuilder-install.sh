#!/bin/sh

#  Copyright 2020 The Kubernetes Authors.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

#
# This file will be  fetched as: curl -L https://git.io/getLatestKubebuilder | sh -
# so it should be pure bourne shell, not bash (and not reference other scripts)
#
# The script fetches the latest kubebuilder release candidate and untars it.
# It lets users to do curl -L https://git.io//getLatestKubebuilder | KUBEBUILDER_VERSION=1.0.5 sh -
# for instance to change the version fetched.

# Check if the program is installed, otherwise exit
function command_exists () {
  if ! [ -x "$(command -v $1)" ]; then
    echo "Error: $1 program is not installed." >&2
    exit 1
  fi
}

# Determine OS
OS="$(uname)"
case $OS in
  Darwin)
    OSEXT="darwin"
    ;;
  Linux)
    OSEXT="linux"
    ;;
  *)
    echo "Only OSX and Linux OS are supported !"
    exit 1
    ;;
esac

HW=$(uname -m)
case $HW in
    x86_64)
      ARCH=amd64 ;;
    *)
      echo "Only x86_64 machines are supported !"
      exit 1
      ;;
esac

# Check if curl, tar commands/programs exist
command_exists curl
command_exists tar

if [ "x${KUBEBUILDER_VERSION}" = "x" ] ; then
  KUBEBUILDER_VERSION=$(curl -L -s https://api.github.com/repos/kubernetes-sigs/kubebuilder/releases/latest | \
                  grep tag_name | sed "s/ *\"tag_name\": *\"\\(.*\\)\",*/\\1/")
  if [ -z "$KUBEBUILDER_VERSION" ]; then
    echo "\nUnable to fetch the latest version tag. This may be due to network access problem"
    exit 0
  fi
fi

KUBEBUILDER_VERSION=${KUBEBUILDER_VERSION#"v"}
KUBEBUILDER_VERSION_NAME="kubebuilder_${KUBEBUILDER_VERSION}"
KUBEBUILDER_DIR=/usr/local/kubebuilder

# Check if folder containing kubebuilder executable exists and is not empty
if [ -d "$KUBEBUILDER_DIR" ]; then
  if [ "$(ls -A $KUBEBUILDER_DIR)" ]; then
    echo "\n/usr/local/kubebuilder folder is not empty. Please delete or backup it before to install ${KUBEBUILDER_VERSION_NAME}"
    exit 1
  fi
fi

TMP_DIR=$(mktemp -d)
pushd $TMP_DIR

# Downloading Kubebuilder compressed file using curl program
URL="https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${OSEXT}_${ARCH}"
echo "Downloading ${KUBEBUILDER_VERSION_NAME}\nfrom $URL"
curl -L "$URL" -o ${TMP_DIR}/kubebuilder
chmod 775 ${TMP_DIR}/kubebuilder

if [[ -z "${PROW}" ]]; then 
  MOVE="sudo mv"
else
  MOVE="mv"
fi

cd ${TMP_DIR}
echo "Moving files to $KUBEBUILDER_DIR folder\n"
sudo mkdir -p /usr/local/kubebuilder/bin
sudo mv -f kubebuilder /usr/local/kubebuilder/bin

echo "Add kubebuilder to your path; e.g copy paste in your shell and/or edit your ~/.profile file"
echo "export PATH=\$PATH:/usr/local/kubebuilder/bin"
popd
rm -rf $TMP_DIR