oc delete mustgather --all -n test-sue
oc delete subscription isv-must-gather-operator -n openshift-operators
oc delete catalogsource isv-must-gather-operator-catalog -n openshift-marketplace

for i in $(oc get csvs  --all-namespaces|grep isv-must |awk '{print $1}');do oc delete csvs isv-must-gather-operator.v0.2.0 -n $i ;done
oc delete crd mustgathers.isv.operator.com
