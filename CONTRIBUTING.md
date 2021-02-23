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

