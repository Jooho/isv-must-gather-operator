# How to use ISV Must Gather Operator

## Flow
- Deploy CatalogSource
- Install ISV Must-Gather Operator via OpenShift Console
- Create MustGather CR


## Steps
### Deploy CatalogSource
~~~
git clone https://github.com/Jooho/isv-must-gather-operator.git

cd isv-must-gather-operator

# Login an OpenShift Cluster with cluster admin user
oc login 

make cs-deploy
~~~

### Install ISV Must-Gather Operator via OpenShift Console
  ![Image](images/isv-must-gather-operator-1.png)
  ![Image](images/isv-must-gather-operator-2.png)
  ![Image](images/isv-must-gather-operator-3.png)
  ![Image](images/isv-must-gather-operator-4.png)


### Create MustGather CR
~~~
# Login an OpenShift Cluster with admin user
oc login

oc new-project test

make cr-deploy

or

echo "apiVersion: isv.operator.com/v1alpha1
kind: MustGather
metadata:
  name: mustgather-test
spec: 
  mustGatherImgURL: \"quay.io/jooholee/isv-smoke-must-gather:0.2.0\""|o create -f -
~~~
