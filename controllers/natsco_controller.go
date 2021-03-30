/*
Copyright 2021.

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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cecov1alpha1 "github.com/inspursoft/ceco/api/v1alpha1"
)

// NatsCoReconciler reconciles a NatsCo object
type NatsCoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ceco.board.io,resources=natscoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ceco.board.io,resources=natscoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ceco.board.io,resources=natscoes/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NatsCo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *NatsCoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("natsco", req.NamespacedName)

	// Fetch the NatsCo instance
	natsco := &cecov1alpha1.NatsCo{}
	err := r.Get(ctx, req.NamespacedName, natsco)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("NatsCo resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get NatsCo")
		return ctrl.Result{}, err
	}

	// Check if the source deployment already exists, if not create a new one
	source := natsco.Spec.Source
	destinations := natsco.Spec.Destinations
	sourceIP := "127.0.0.1"
	// Check source deployment
	sourceDeployment, err := r.checkDeployment(ctx, source, natsco, sourceIP, log)
	if err != nil {
		return ctrl.Result{}, err
	}
	sourceDeploymentName := sourceDeployment.ObjectMeta.Name
	sourcePodMatchLabels := sourceDeployment.Spec.Selector.MatchLabels
	sourcePodList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(natsco.Namespace),
		client.MatchingLabels(sourcePodMatchLabels),
	}
	if err = r.List(ctx, sourcePodList, listOpts...); err != nil {
		log.Error(err, "Failed to list source pods", sourcePodMatchLabels)
		return ctrl.Result{}, err
	}
	for _, pod := range sourcePodList.Items {
		if pod.Status.ContainerStatuses[0].Ready {
			sourceIP = pod.Status.PodIP
			break
		}
	}
	if sourceIP == "127.0.0.1" {
		log.Info("Source pod not ready.")
		return ctrl.Result{RequeueAfter: time.Second * 3}, nil
	}

	// Check destinations deployments
	for _, d := range destinations {
		if _, err = r.checkDeployment(ctx, d, natsco, sourceIP, log); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Update the NatsCo status with the pod names
	// List the pods for this natsco's deployment
	podList := &corev1.PodList{}
	operatortesterCommonLabels := map[string]string{"app": "natsco", "operatortester_cr": natsco.Name}
	listOpts = []client.ListOption{
		client.InNamespace(natsco.Namespace),
		client.MatchingLabels(operatortesterCommonLabels),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "NatsCo.Namespace", natsco.Namespace, "NatsCo.Name", natsco.Name)
		return ctrl.Result{}, err
	}
	log.Info("Get pod list")
	podStatus := getPodStatus(podList.Items)
	for podName, status := range podStatus {
		if exist := strings.Contains(podName, sourceDeploymentName); exist {
			natsco.Status.Source = podName + ":" + status
			delete(podStatus, podName)
			break
		}
	}
	natsco.Status.Destination = podStatus
	err = r.Status().Update(ctx, natsco)
	if err != nil {
		log.Error(err, "Failed to update NatsCo status")
		return ctrl.Result{}, err
	}
	log.Info("NatsCo.Name", natsco.Name, "Status has updated")

	// Reconcile for any reason other than an error after 5 seconds
	return ctrl.Result{RequeueAfter: time.Second * 5}, nil
}

func (r *NatsCoReconciler) checkDeployment(ctx context.Context, hostAndPath cecov1alpha1.HostAndPath,
	natsco *cecov1alpha1.NatsCo, sourcePodIP string, log logr.Logger) (*appsv1.Deployment, error) {
	found := &appsv1.Deployment{}
	coType := natsco.Spec.CoType
	sourceFile := natsco.Spec.Source.FilePath
	destinationFile := hostAndPath.FilePath
	filePath := strings.Split(destinationFile, "/")
	deploymentName := hostAndPath.Hostname + "-" + filePath[len(filePath)-1]
	err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: natsco.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		natsServs := strings.Join(natsco.Spec.NatsServers, ",")
		deployment := createDeployment(deploymentName, natsco.Name, natsco.Namespace, hostAndPath.Hostname, coType, natsServs, sourceFile, destinationFile, sourcePodIP)
		ctrl.SetControllerReference(natsco, deployment, r.Scheme)
		err = r.Create(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return nil, err
		}
		log.Info("Create a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		found = deployment
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return nil, err
	} else {
		log.Info("Deployment has deployed", "Deployment.Namespace", natsco.Namespace, "Deployment.Name", deploymentName)
	}
	return found, nil
}

// createDeployment for creating the target deployment
func createDeployment(name, operatortesterName, ns, hostname, coType, natsServs, sourceFile, destinationFile, sourcePodIP string) *appsv1.Deployment {
	replicas := int32(1)
	labels := labelsForOperatorTester(operatortesterName, name)
	volumeName := name + "-volume"
	image := "172.31.0.7:5000/source-nats:v0.0.2"
	containerPort := int32(8080)
	if sourcePodIP != "127.0.0.1" {
		image = "172.31.0.7:5000/destination-nats:v0.0.2"
		containerPort = 8081
	}
	hostVolume, hostVolumeMount := volumeForOperatorTester(volumeName, sourceFile, destinationFile, sourcePodIP)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": hostname,
					},
					Containers: []corev1.Container{{
						Image: image,
						Name:  name,
						Env:   envsForOperatorTester(natsServs, sourceFile, destinationFile, sourcePodIP),
						Ports: []corev1.ContainerPort{{
							ContainerPort: containerPort,
							Name:          "http",
						}},
						VolumeMounts: hostVolumeMount,
					}},
					Volumes: hostVolume,
				},
			},
		},
	}
}

// labelsForOperatorTester returns the labels for selecting the resources
// belonging to the given natsco CR name.
func labelsForOperatorTester(operatortesterName, name string) map[string]string {
	return map[string]string{"app": "natsco", "operatortester_cr": operatortesterName, "meta": name}
}

// envsForOperatorTester returns env for conatiner
func envsForOperatorTester(natsServs, sourceFile, destinationFile, sourcePodIP string) []corev1.EnvVar {
	var envs []corev1.EnvVar
	envs = append(envs,
		corev1.EnvVar{Name: "NATS_SERVERS", Value: natsServs},
		corev1.EnvVar{Name: "SOURCE", Value: sourceFile},
		corev1.EnvVar{Name: "DESTINATION", Value: destinationFile},
		corev1.EnvVar{Name: "SOURCE_IP", Value: sourcePodIP + ":8080"})
	return envs
}

func volumeForOperatorTester(volumeName, sourceFile, destinationFile, sourcePodIP string) ([]corev1.Volume, []corev1.VolumeMount) {
	var hostVolume corev1.HostPathVolumeSource
	var hostVolumeMount []corev1.VolumeMount
	hostVolumeType := corev1.HostPathFile
	if sourcePodIP == "127.0.0.1" {
		hostVolume = corev1.HostPathVolumeSource{
			Path: sourceFile,
			Type: &hostVolumeType,
		}
		hostVolumeMount = []corev1.VolumeMount{{
			Name:      volumeName,
			MountPath: sourceFile,
		}}
	} else {
		hostVolumeType = corev1.HostPathDirectoryOrCreate
		dir := filepath.Dir(destinationFile)
		hostVolume = corev1.HostPathVolumeSource{
			Path: dir,
			Type: &hostVolumeType,
		}
		hostVolumeMount = []corev1.VolumeMount{{
			Name:      volumeName,
			MountPath: dir,
		}}
	}
	return []corev1.Volume{{
		Name: volumeName,
		VolumeSource: corev1.VolumeSource{
			HostPath: &hostVolume,
		},
	}}, hostVolumeMount
}

// getPodStatus returns the pod status of the array of pods passed in
func getPodStatus(pods []corev1.Pod) map[string]string {
	podStatus := make(map[string]string)
	for _, pod := range pods {
		podStatus[pod.Name] = strconv.FormatBool(pod.Status.ContainerStatuses[0].Ready)
	}
	return podStatus
}

// SetupWithManager sets up the controller with the Manager.
func (r *NatsCoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cecov1alpha1.NatsCo{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
