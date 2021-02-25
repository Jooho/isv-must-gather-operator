# CONTRIBUTING

## Local Test
~~~
make install              // Deploy CRD on the cluster
make run                  // Deploy isv-must-gather operator on local
make t-cr or t-cr-event   // t-cr (full testing), t-cr-event(only events object)
~~~

## Cluster Test
~~~
make t-deploy-all     // deploy MustGather Operator on the cluster
make t-cr or t-cr-event      // t-cr (full testing), t-cr-event(only events object)

make t-undeploy  //delete MustGather CR, all MustGather related objects
~~~

- Additional make macros
  - `make t-del-cr` = Delete CR only, Not operator
  - `make t-deploy` = Deploy operator only, Not build/push operator image
  - `make p-image`  = Build/push operator image

## Bundle image
~~~
vi ./env.sh
export BUNDLE_TAG=0.2.0  //Update Version 

make bundle-image        // Build bundle image
make bundle-validate     // validate the built bundle image
make bundle-run          // Deploy isv must-gather operator
make bundle-clean        // Delete 
~~~
- Additional make macros
  - `make bundle-build` = Build bundle image
  - `make bundle-push`  = Push bundle image
## Index image
~~~
make index-image         // Build index image
make cs-deploy           // Deploy catalogSource
make cs-undeploy         // Delete MustGather CR in all namespaces, delete subscriptions/csv/installplans/catalogsource
~~~
- Additional make macros
  - `make index-build` = Build index image
  - `make index push`  = Push index image



## Practical Methods

- Full flow to build all images (operator,bundle,index)
  ~~~
  make p-image bundle-image index-image
  ~~~

## Clean repo Before commit
~~~
make clean
~~~
## Image Description for manifests
ISV Must Gather Operator help you to gather debugging data without download any binary. One thing you should do is creating a MustGather Operand
Example MustGather Operand
~~~
apiVersion: isv.operator.com/v1alpha1
kind: MustGather
  metadata:
    name: mustgather
spec:
  mustGatherImgURL: "quay.io/jooholee/isv-smoke-must-gather:0.2"        <== Update this parameter
~~~
You can find mustGatherImgURL from the installed operator page that you are using.

In order to watch the log status, you can execute the following command:
~~~
oc logs isv-cli-pod -f
~~~
  
After finished gathering, you can download mustgather.tar via browser or wget.
You can check downloadURL with this command:
~~~    
oc get mustgather -o jsonpath="{ .items[*].status.downloadURL}"
~~~ 
