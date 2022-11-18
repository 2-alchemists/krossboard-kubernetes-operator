/*
Copyright 2022.

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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/exp/slices"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	kbapi "krossboard-kubernetes-operator/api/v1alpha1"
)

// KrossboardReconciler reconciles a Krossboard object
type KrossboardReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=krossboard.krossboard.app,resources=krossboards,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=krossboard.krossboard.app,resources=krossboards/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=krossboard.krossboard.app,resources=krossboards/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Here it compares the state specified by the Krossboard object against
// the actual cluster state, and then perform operations to make the cluster state
// reflect the state specified by the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (kbReconciler *KrossboardReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// Fetch the Krossboard instance
	kb := &kbapi.Krossboard{}
	err := kbReconciler.Get(ctx, req.NamespacedName, kb)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Krossboard resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Krossboard")
		return ctrl.Result{}, err
	}

	// Check if the deployment already exists, if not create a new one
	kbDeployment := &appsv1.Deployment{}
	err = kbReconciler.Get(ctx, types.NamespacedName{Name: kb.Name, Namespace: kb.Namespace}, kbDeployment)
	if err != nil && errors.IsNotFound(err) {
		dep := kbReconciler.deploymentForKrossboard(kb, ctx, req)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = kbReconciler.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Krossboard Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as expected
	kbExpectedReplicas := int32(KbReplicaCount)
	if *kbDeployment.Spec.Replicas != kbExpectedReplicas {
		kbDeployment.Spec.Replicas = &kbExpectedReplicas
		err = kbReconciler.Update(ctx, kbDeployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment",
				"Deployment.Namespace", kbDeployment.Namespace,
				"Deployment.Name", kbDeployment.Name)
			return ctrl.Result{}, err
		}
		// Ask to requeue after 1 minute in order to give enough time for the
		// pods be created on the cluster side and the operand be able
		// to do the next update step accurately.
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Update the Krossboard status with the pod names
	// List the pods for this krossboard's deployment
	kbPods := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(kb.Namespace),
		client.MatchingLabels(labelsForKrossboard(kb.Name)),
	}

	if err = kbReconciler.List(ctx, kbPods, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Krossboard.Namespace", kb.Namespace, "Krossboard.Name", kb.Name)
		return ctrl.Result{}, err
	}

	kbPodsRaw, _ := json.Marshal(kbPods)
	fmt.Println(string(kbPodsRaw))

	kbContainers := getKrossboardContainers(kbPods.Items)
	if !reflect.DeepEqual(kbContainers, kb.Status.KoaInstances) {
		kb.Status.KoaInstances = kbContainers.KoaInstances
		kb.Status.KbComponentInstances = kbContainers.KbComponentInstances
		err := kbReconciler.Status().Update(ctx, kb)
		if err != nil {
			log.Error(err, "Failed to update Krossboard status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// labelsForKrossboard returns the labels for selecting the resources
// belonging to the given krossboard CR name.
func labelsForKrossboard(name string) map[string]string {
	return map[string]string{"app": "krossboard", "krossboard-kubernetes-operator": name}
}

// getKrossboardContainers returns a KrossboardStatus object
func getKrossboardContainers(pods []corev1.Pod) *kbapi.KrossboardStatus {
	kbStatus := &kbapi.KrossboardStatus{}
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			currentInstance := kbapi.KoaInstance{
				Name: container.Name,
			}

			if len(container.Env) > 0 {
				idx := slices.IndexFunc(container.Env, func(e corev1.EnvVar) bool { return e.Name == "KOA_LISTENER_PORT" })
				if idx >= 0 {
					currentInstance.ContainerPort, _ = strconv.ParseInt(container.Env[idx].Value, 10, 64)
				}

				idx = slices.IndexFunc(container.Env, func(e corev1.EnvVar) bool { return e.Name == "KOA_CLUSTER_NAME" })
				if idx >= 0 {
					currentInstance.ClusterName = container.Env[idx].Value
				}

				idx = slices.IndexFunc(container.Env, func(e corev1.EnvVar) bool { return e.Name == "KOA_K8S_API_ENDPOINT" })
				if idx >= 0 {
					currentInstance.ClusterEndpointURL = container.Env[idx].Value
				}
			}

			if currentInstance.ClusterName != "" && currentInstance.ClusterEndpointURL != "" && currentInstance.ContainerPort != 0 {
				kbStatus.KoaInstances = append(kbStatus.KoaInstances, currentInstance)
			} else {
				kbStatus.KbComponentInstances = append(kbStatus.KbComponentInstances,
					kbapi.KbComponentInstance{
						Name:          currentInstance.Name,
						ContainerPort: currentInstance.ContainerPort,
					})
			}
		}
	}
	return kbStatus
}

// SetupWithManager sets up the controller with the Manager.
func (r *KrossboardReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kbapi.Krossboard{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
