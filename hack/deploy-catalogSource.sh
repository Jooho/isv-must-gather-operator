
source env.sh

cat << EOF | oc apply -f - 
apiVersion: operators.coreos.com/v1alpha1 
kind: CatalogSource 
metadata: 
  name: isv-must-gather-operator-catalog 
  namespace: openshift-marketplace 
spec: 
  sourceType: grpc 
  image: ${INDEX_IMG} 
EOF