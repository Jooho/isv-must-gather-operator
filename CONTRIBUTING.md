# CONTRIBUTING

## Local Test
~~~
make install
make run
make t-cr or t-cr-event
~~~

## Cluster Test
~~~
make t-deploy-all
make t-cr or t-cr-event

make t-undeploy 
~~~

- Additional make macros
  - `make t-del-cr` = Delete CR only not operator
  - `make t-deploy` = Deploy operator only not build/push operator image
  - `make p-image` = Build/push operator image

## Bundle image
~~~
vi ./env.sh
export BUNDLE_TAG=0.2.0  //Update Version 

make bundle-image
make bundle-validate
make bundle-run
make bundle-clean
~~~
- Additional make macros
  - `make bundle-build` = Build bundle image
  - `make bundle-push`  = Push bundle image
## Index image
~~~
make index-image
make cs-deploy
make cs-undeploy
~~~
- Additional make macros
  - `make index-build` = Build index image
  - `make index push`  = Push index image



## Practical Methods
~~~
# Update operator
make p-image bundle-image index-image
~~~


## Before commit
~~~
make clean
~~~
## Image Description
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
