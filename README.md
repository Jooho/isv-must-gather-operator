# isv-must-gather-operator



Image description
~~~
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
  ~~~