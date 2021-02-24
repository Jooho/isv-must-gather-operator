package defaults

const (
	//DestDir is for webserver context-root and also mustgather.tar path. Do not change this DestDir.
	DestDir              = "/opt/download"
	//IsvCliImgVersion will be updated by make command
	IsvCliImgVersion     = "0.2.1"
	//MustGatherImgVersion will be updated by make command
	MustGatherImgVersion = "all"
)

var (
	// NodeSelector is for the isv-cli image. This image must run on one of workers Not masters
	NodeSelector = map[string]string{"node-role.kubernetes.io/worker=": ""}

	//ServiceAccount is the default sc name for deploying isv-cli image and this need admin role level because isv-cli image will deploy must-gather image that need admin permission
	ServiceAccount = "isv-cli-sa"

	//RoleBinding is the default rb name for mapping between ServiceAccount and RoleBinding.
	RoleBinding = "isv-cli-rb"

	//Deployment is for isv-cli image
	Deployment = "isv-cli-deploy"

	//Pod is the default pod name for isv-cli pod image
	Pod = "isv-cli-pod"

	// Service is the default service name
	Service = "isv-cli-svc"

	//Route is the default route name
	Route = "isv-cli-route"

	//IsvCliImg is the container image contains isv-cli binary. Version have to map among isv-cli binary and isv-cli image and base must-gather image.
	IsvCliImg = "quay.io/jooholee/isv-cli:" + IsvCliImgVersion

	//MustGatherImgURL should be updated based on ISV operator. This is the default that gather namespace scope data only.
	MustGatherImgURL = "quay.io/jooholee/isv-smoke-must-gather:" + MustGatherImgVersion
)
