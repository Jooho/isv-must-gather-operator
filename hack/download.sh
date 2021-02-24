#!/bin/bash

source ../env.sh

cd ${TMP_HOME}
## Download Binary & Move them in the path
echo "curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${SDK_RELEASE_VERSION}/operator-sdk-linux-amd64"
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${SDK_RELEASE_VERSION}/operator-sdk_linux_amd64
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${SDK_RELEASE_VERSION}/ansible-operator_linux_amd64
curl -LO https://github.com/operator-framework/operator-sdk/releases/download/${SDK_RELEASE_VERSION}/helm-operator_linux_amd64
curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz | tar xz
curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"

curl -LO https://github.com/operator-framework/operator-registry/releases/download/v${OPM_VERSION}/linux-amd64-opm

chmod +x operator-sdk_linux_amd64 && sudo mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk
chmod +x ansible-operator_linux_amd64 && sudo mv ansible-operator_linux_amd64 /usr/local/bin/ansible-operator
chmod +x helm-operator_linux_amd64 && sudo mv helm-operator_linux_amd64 /usr/local/bin/helm-operator

chmod +x linux-amd64-opm && sudo mv linux-amd64-opm /usr/local/bin/opm

chmod +x kustomize  && sudo mv kustomize /usr/local/bin/kustomize
chmod +x kubectl && sudo mv kubectl /usr/local/bin/kubectl


./kubebuilder-install.sh
