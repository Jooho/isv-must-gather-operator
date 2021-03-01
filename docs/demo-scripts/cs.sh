echo "
apiVersion: operators.coreos.com/v1alpha1 
kind: CatalogSource 
metadata: 
  name: isv-must-gather-operator-catalog 
  namespace: openshift-marketplace 
spec: 
  sourceType: grpc 
  image: quay.io/jooholee/isv-must-gather-operator-index:0.2.0" | oc create -f - 
