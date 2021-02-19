#!/bin/bash
export ROOT_HOME=/home/jooho/dev/Managed_Git/operator-projects   #Update (For demo, use "/tmp")
export TMP_HOME=/tmp

# Set the release version for operator sdk
export SDK_RELEASE_VERSION=v1.4.0  #Update  latest version v1.4.0
export KUSTOMIZE_VERSION=3.10.0
export OPM_VERSION=1.16.1

# Operator Info
export REPO_URL=github.com/jooho/isv-must-gather-operator
export NEW_OP_NAME=isv-must-gather-operator
export NEW_OP_HOME=${ROOT_HOME}/${NEW_OP_NAME}
export NAMESPACE=${NEW_OP_NAME}

# Images Info
export IMG_TAG=v0.0.1
export IMG=quay.io/jooholee/${NEW_OP_NAME}:${IMG_TAG}
export BUNDLE_TAG=v0.0.1
export BUNDLE_IMG=quay.io/jooholee/${NEW_OP_NAME}:${BUNDLE_TAG}
export INDEX_TAG=v0.0.1
export INDEX_IMG=quay.io/jooholee/${NEW_OP_NAME}-index:{INDEX_TAG}
export CHANNELS=alpha
export DEFAULT_CHANNEL=alpha

# CRD Info
export DOMAIN=operator.com
export GROUP=isv
export VERSION=v1alpha1
export KIND=MustGather

