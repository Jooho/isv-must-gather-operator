/*
Copyright 2021 Jooho Lee.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"path"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	// "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	// "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	isvv1alpha1 "github.com/jooho/isv-must-gather-operator/api/v1alpha1"
	"github.com/jooho/isv-must-gather-operator/controllers/defaults"
)

// MustGatherReconciler reconciles a MustGather object
type MustGatherReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=isv.operator.com,resources=mustgathers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=isv.operator.com,resources=mustgathers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=isv.operator.com,resources=mustgathers/finalizers,verbs=update

// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MustGather object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *MustGatherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("mustgather", req.NamespacedName)

	// Fetch the MustGather instance
	mustgather := &isvv1alpha1.MustGather{}
	if err := r.Get(ctx, req.NamespacedName, mustgather); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("MustGather resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get MustGather")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Validate checking

	// Object Creation

	// ServiceAccount
	// Check if the serviceaccount already exists, if not create a new one
	sa := &corev1.ServiceAccount{}

	err := r.Get(ctx, types.NamespacedName{Name: defaults.ServiceAccount, Namespace: mustgather.Namespace}, sa)
	if err != nil && errors.IsNotFound(err) {
		// Define a new serviceaccount
		sa = r.newSA(mustgather)

		log.Info("Creating a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)

		if err := r.Create(ctx, sa); err != nil {
			log.Error(err, "Failed to create a new Serviceaccount", "Serviceaccount.Namespace", sa.Namespace, "Serviceaccount.Name", sa.Name)
			return ctrl.Result{}, err
		}

	}

	// RoleBinding
	rb := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.RoleBinding, Namespace: mustgather.Namespace}, rb)
	if err != nil && errors.IsNotFound(err) {
		rb = r.newRoleBinding(sa.Name, mustgather)

		log.Info("Creating a new RoleBinding", "RoleBinding.Namespace", rb.Namespace, "roleBinding.Name", rb.Name)

		if err := r.Create(ctx, rb); err != nil {
			log.Error(err, "Failed to create a RoleBinding for MustGather", "roleBinding.Namespace", rb.Namespace, "roleBinding.Name", rb.Name)
			return ctrl.Result{}, err
		}
	}

	// Deployment
	dep := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.Pod, Namespace: mustgather.Namespace}, dep)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		// dep = r.newDeployment(sa.Name, mustgather)

		dep = r.newPod(sa.Name, mustgather)

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create a Deployment for MustGather", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

	}

	// Service
	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.Service, Namespace: mustgather.Namespace}, svc)
	if err != nil && errors.IsNotFound(err) {

		svc = r.newService(mustgather)

		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		if err = r.Create(ctx, svc); err != nil {
			log.Error(err, "Failed to create a Service for MustGather", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}

	}

	// Route
	route := &routev1.Route{}
	err = r.Get(ctx, types.NamespacedName{Name: defaults.Route, Namespace: mustgather.Namespace}, route)
	if err != nil && errors.IsNotFound(err) {

		route = r.newRoute(mustgather)

		log.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		if err = r.Create(ctx, route); err != nil {
			log.Error(err, "Failed to create a Service for MustGather", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return ctrl.Result{}, err
		}

	}

	//Update MustGather
	mustgather.Status.DownloadURL = "http://" + route.Spec.Host + "/download"
	if err := r.Status().Update(ctx, mustgather); err != nil {
		log.Error(err, "Failed to update the MustGather Status for download url", "MustGather.Namespace", mustgather.Namespace, "MustGather.Name", mustgather.Name)
		return ctrl.Result{}, err
	}

	// Delete Logic
	// name of our custom finalizer
	const finalizerName = "isv.finalizers.operator.io"

	// examine DeletionTimestamp to determine if object is under deletion

	if mustgather.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.

		if !containsString(mustgather.GetFinalizers(), finalizerName) {
			log.Info("Adding Finalizer for the MustGather")

			controllerutil.AddFinalizer(mustgather, finalizerName)

			if err := r.Update(context.Background(), mustgather); err != nil {
				log.Error(err, "Failed to update CR MustGather to add finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(mustgather.ObjectMeta.Finalizers, finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteRelatedResources(ctx, mustgather); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried

				log.Error(err, "Failed to delete external resoureces")
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			log.Info("Removing Finalizer for the MustGather")
			controllerutil.RemoveFinalizer(mustgather, finalizerName)

			if err := r.Update(context.Background(), mustgather); err != nil {
				log.Error(err, "Failed to update CR MustGather with finalizer to remove finalizer")
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *MustGatherReconciler) newSA(mg *isvv1alpha1.MustGather) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{

		ObjectMeta: metav1.ObjectMeta{
			Name:      defaults.ServiceAccount,
			Namespace: mg.Namespace,
			Labels: map[string]string{
				"app":  "must-gather",
				"bin":  "isv-cli",
				"kind": "sa",
			},
		},
	}
	// Set MustGather instance as the owner and controller
	ctrl.SetControllerReference(mg, sa, r.Scheme)
	return sa
}

func (r *MustGatherReconciler) newRoleBinding(sa string, mg *isvv1alpha1.MustGather) *rbacv1.RoleBinding {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaults.RoleBinding,
			Namespace: mg.Namespace,
			Annotations: map[string]string{
				"oc.openshift.io/command": "isv-cli must-gather",
			},
			Labels: map[string]string{
				"app":  "must-gather",
				"bin":  "isv-cli",
				"kind": "rb",
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: sa,
			},
		},
	}

	// Set MustGather instance as the owner and controller
	ctrl.SetControllerReference(mg, rb, r.Scheme)

	return rb
}

// Delete any related resources associated with the MustGather
func (r *MustGatherReconciler) deleteRelatedResources(ctx context.Context, mg *isvv1alpha1.MustGather) error {

	if err := r.DeleteAllOf(ctx, &corev1.Pod{}, client.InNamespace(mg.Namespace), client.MatchingLabels{"app": "must-gather"}); err != nil {
		return err
	}

	return nil
}

// // newPod return isv-cli pod specification
func (r *MustGatherReconciler) newPod(sa string, mg *isvv1alpha1.MustGather) *corev1.Pod {
	zero := int64(0)
	isvImg := defaults.IsvCliImg
	mgImg := defaults.MustGatherImgURL

	if mg.Spec.MustGatherImgURL != "" {
		mgImg = mg.Spec.MustGatherImgURL
	}

	// nodeSelector := defaults.NodeSelector

	ret := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaults.Pod,
			Namespace: mg.Namespace,
			Labels: map[string]string{
				"app":  "must-gather",
				"bin":  "isv-cli",
				"kind": "pod",
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy:      corev1.RestartPolicyNever,
			ServiceAccountName: sa,
			Volumes: []corev1.Volume{
				{
					Name: "must-gather-download",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            "isv-cli",
					Image:           isvImg,
					ImagePullPolicy: corev1.PullAlways,
					// always force disk flush to ensure that all data gathered is accessible in the copy container
					Command: []string{"/bin/bash", "-c", "isv-cli must-gather --image " + mgImg + " --dest-dir /opt/download --browser ; sync"},
					Env: []corev1.EnvVar{
						{
							Name: "NAMESPACE",
							ValueFrom: &corev1.EnvVarSource{
								FieldRef: &corev1.ObjectFieldSelector{
									FieldPath: "metadata.namespace",
								},
							},
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "must-gather-download",
							MountPath: path.Clean(defaults.DestDir),
							ReadOnly:  false,
						},
					},
				},
				{
					Name:            "copy",
					Image:           isvImg,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         []string{"/bin/bash", "-c", "trap : TERM INT; sleep infinity & wait"},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "must-gather-download",
							MountPath: path.Clean(defaults.DestDir),
							ReadOnly:  false,
						},
					},
				},
			},
			TerminationGracePeriodSeconds: &zero,
			Tolerations: []corev1.Toleration{
				{
					Operator: "Exists",
				},
			},
		},
	}
	// Set MustGather instance as the owner and controller
	ctrl.SetControllerReference(mg, ret, r.Scheme)
	return ret
}

func (r *MustGatherReconciler) newService(mg *isvv1alpha1.MustGather) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaults.Service,
			Namespace: mg.Namespace,
			Labels: map[string]string{
				"app":  "must-gather",
				"bin":  "isv-cli",
				"kind": "svc",
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     "isv-web-server",
					Port:     int32(9000),
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(9000),
						StrVal: "8080",
					},
				},
			},
			Selector: map[string]string{
				"app": "must-gather",
				"bin": "isv-cli",
			},
		},
	}
	ctrl.SetControllerReference(mg, service, r.Scheme)
	return service
}
func (r *MustGatherReconciler) newRoute(mg *isvv1alpha1.MustGather) *routev1.Route {
	route := &routev1.Route{

		ObjectMeta: metav1.ObjectMeta{
			Name:      defaults.Route,
			Namespace: mg.Namespace,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: defaults.Service,
			},
		},
	}

	ctrl.SetControllerReference(mg, route, r.Scheme)
	return route
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *MustGatherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&isvv1alpha1.MustGather{}).
		Complete(r)
}
