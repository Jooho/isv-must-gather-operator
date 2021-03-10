export productVerion=4.6
echo "
apiVersion: v1
kind: Namespace
metadata:
  name: openshift-local-storage
---
apiVersion: operators.coreos.com/v1alpha2
kind: OperatorGroup
metadata:
  name: openshift-local-operator-group
  namespace: openshift-local-storage
spec:
  targetNamespaces:
    - openshift-local-storage
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: local-storage-operator
  namespace: openshift-local-storage
spec:
  channel: \"${productVersion}\" 
  installPlanApproval: Automatic
  name: local-storage-operator
  source: redhat-operators
  sourceNamespace: openshift-marketplace" |oc create -f -

sleep 120

oc project openshift-local-storage

echo "
apiVersion: \"local.storage.openshift.io/v1\"
kind: \"LocalVolume\"
metadata:
  name: \"local-disks\"
  namespace: \"openshift-local-storage\" 
spec:
  nodeSelector: 
    nodeSelectorTerms:
    - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - worker-0.bell.tamlab.brq.redhat.com
          - worker-1.bell.tamlab.brq.redhat.com
          - worker-2.bell.tamlab.brq.redhat.com
  storageClassDevices:
    - storageClassName: \"local-sc\"
      volumeMode: Filesystem 
      fsType: xfs 
      devicePaths: 
        - /dev/vdb" | oc create -f -

sleep 90

cat <<EOF | oc apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: nfsprovisioner-catalog
  namespace: openshift-marketplace
spec:
  sourceType: grpc
  image: quay.io/jooholee/nfs-provisioner-operator-index:0.0.1 
EOF

sleep 60

oc new-project nfs-provisioner

cat <<EOF | oc apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: nfs-provisioner-9mrzq
  namespace: nfs-provisioner
spec:
  targetNamespaces:
  - nfs-provisioner
EOF

cat <<EOF | oc apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: nfs-provisioner-operator
  namespace: nfs-provisioner
spec:
  channel: alpha
  installPlanApproval: Automatic
  name: nfs-provisioner-operator
  source: nfsprovisioner-catalog
  sourceNamespace: openshift-marketplace
  startingCSV: nfs-provisioner-operator.v0.0.1
EOF

sleep 60

echo "
apiVersion: cache.jhouse.com/v1alpha1
kind: NFSProvisioner
metadata:
  name: nfsprovisioner-sample
  namespace: nfs-provisioner
spec:
  storageSize: \"80G\"
  scForNFSPvc: local-sc
  SCForNFSProvisioner: nfs"|oc create -f -
