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
	"strings"

	krossboardv1alpha1 "krossboard-kubernetes-operator/api/v1alpha1"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

const KbReplicaCount = 1

// deploymentForKrossboard returns a krossboard Deployment object.
func (kbReconciler *KrossboardReconciler) deploymentForKrossboard(
	ctx context.Context, m *krossboardv1alpha1.Krossboard, req ctrl.Request,
) *appsv1.Deployment {
	log := ctrllog.FromContext(ctx).WithName("krossboard-kubernetes-operator")

	kbLabels := labelsForKrossboard(m.Name)
	kbDataDir := "/krossboard/db"
	kbK8sSecretsDir := "/krossboard/secrets" //nolint:gosec
	kbCredsDir := "/krossboard/creds"        //nolint:gosec
	kbRuntimeDir := "/krossboard/run"

	kbProcessorImage := m.Spec.KrossboardDataProcessorImage
	if strings.TrimSpace(kbProcessorImage) == "" {
		kbProcessorImage = "krossboard/krossboard-data-processor:latest"
	}

	kbUIImage := m.Spec.KrossboardUIImage
	if strings.TrimSpace(kbUIImage) == "" {
		kbUIImage = "krossboard/krossboard-ui:latest"
	}

	koaImage := m.Spec.KoaImage
	if strings.TrimSpace(koaImage) == "" {
		koaImage = "rchakode/kube-opex-analytics:latest"
	}

	kbSecretName := m.Spec.KrossboardSecretName
	if strings.TrimSpace(kbSecretName) == "" {
		kbSecretName = "krossboard-secrets" //nolint:gosec
	}

	kbPersistentVolumeClaim := m.Spec.KrossboardPersistentVolumeClaim
	if strings.TrimSpace(kbPersistentVolumeClaim) == "" {
		kbPersistentVolumeClaim = "krossboard-data-pvc"
	}

	secretQuery := ktypes.NamespacedName{
		Name:      kbSecretName,
		Namespace: req.NamespacedName.Namespace,
	}

	kbSecretResult := &corev1.Secret{}
	err := kbReconciler.Client.Get(ctx, secretQuery, kbSecretResult)
	if err != nil {
		log.Error(err, "failed to get secret", "Secret.Name", kbSecretName)
		return &appsv1.Deployment{}
	}

	managedClusters := &map[string]*ManagedCluster{}

	// Credentials of managed clusters are primarily expected in the 'managedClusters' key
	secretDataB64, secretKeyFound := kbSecretResult.Data["managedClusters"]
	if secretKeyFound { //nolint:nestif
		log.Info("using managedClusters key in secret", "Secret.Name", kbSecretName)
		err := json.Unmarshal(secretDataB64, managedClusters)
		if err != nil {
			log.Error(err, "failed decoding managedClusters key", "Secret.Name", kbSecretName)
			return &appsv1.Deployment{}
		}
	} else {
		// alternatively the 'kubeconfig' key can be set with a KUBECONFIG content
		secretDataB64, secretKeyFound = kbSecretResult.Data["kubeconfig"]
		if secretKeyFound {
			log.Info("using kubeconfig key in secret", "Secret.Name", kbSecretName)
			*managedClusters, err = (&KubeConfigManager{}).GetManagedClustersFromData(secretDataB64)
			if err != nil {
				log.Error(err, "failed parsing kubeconfig key in secret", "Secret.Name", kbSecretName)
				return &appsv1.Deployment{}
			}
		} else {
			log.Error(err, "neither kubeconfig nor managedCluster key found in secret", "Secret.Name", kbSecretName)
			return &appsv1.Deployment{}
		}
	}

	log.Info("parsing of managed clusters completed", "count", len(*managedClusters))

	koaContainers := []corev1.Container{}
	koaContainerPort := int32(5483)
	clusterNames := make([]string, 0, len(*managedClusters))
	for _, managedCluster := range *managedClusters {
		clusterNames = append(clusterNames, managedCluster.Name)
		koaContainerName := fmt.Sprintf("kube-opex-analytics-%v", uuid.New().String())
		koaContainers = append(koaContainers,
			corev1.Container{
				Image: koaImage,
				Name:  koaContainerName,
				Ports: []corev1.ContainerPort{{
					ContainerPort: koaContainerPort,
					Name:          "http",
					Protocol:      "TCP",
				}},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "krossboard-data-vol",
						MountPath: kbDataDir,
						ReadOnly:  false,
					},
					{
						Name:      "krossboard-creds-vol",
						MountPath: kbCredsDir,
						ReadOnly:  false,
					},
				},
				Env: []corev1.EnvVar{
					{
						Name:  "KOA_LISTENER_PORT",
						Value: fmt.Sprint(koaContainerPort),
					},
					{
						Name:  "KOA_CLUSTER_NAME",
						Value: managedCluster.Name,
					},
					{
						Name:  "KOA_DB_LOCATION",
						Value: fmt.Sprintf("%s/%s", kbDataDir, managedCluster.Name),
					},
					{
						Name:  "KOA_K8S_API_ENDPOINT",
						Value: managedCluster.APIEndpoint,
					},
					{
						Name:  "KOA_K8S_API_VERIFY_SSL",
						Value: "true",
					},
					{
						Name:  "KOA_K8S_CACERT",
						Value: fmt.Sprintf("%s/%s/cacert.pem", kbCredsDir, managedCluster.Name),
					},
					{
						Name:  "KOA_K8S_AUTH_TOKEN_FILE",
						Value: fmt.Sprintf("%s/%s/token", kbCredsDir, managedCluster.Name),
					},
				},
			})
		koaContainerPort++
	}

	kbReplicas := int32(KbReplicaCount)
	kbContainerUsername := "krossboard"
	koaContainerFsGroup := int64(4583)
	koaContainerUID := int64(4583)
	selectedClusterNames := strings.Trim(strings.Join(clusterNames, " "), " ")

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &kbReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: kbLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: kbLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "krossboard-kb-kubernetes-operator",
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: &koaContainerFsGroup,
					},
					InitContainers: []corev1.Container{
						{
							Image:   koaImage,
							Name:    "krossboard-init",
							Command: []string{"/bin/bash", "-c", "--"},
							Args: []string{
								fmt.Sprintf("for cn in '%s'; do mkdir -p %s/$cn ; done ;", selectedClusterNames, kbDataDir),
								fmt.Sprintf("chown -R %v:%v %s ;", koaContainerUID, koaContainerUID, kbDataDir),
								fmt.Sprintf("chown -R %v:%v %s ;", kbContainerUsername, kbContainerUsername, kbCredsDir),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "krossboard-data-vol",
									MountPath: kbDataDir,
									ReadOnly:  false,
								},
								{
									Name:      "krossboard-creds-vol",
									MountPath: kbCredsDir,
									ReadOnly:  false,
								},
							},
						},
					},
					Containers: append(
						[]corev1.Container{
							{
								Image: kbUIImage,
								Name:  "krossboard-ui",
								Ports: []corev1.ContainerPort{{
									ContainerPort: 80,
									Name:          "krossboard-ui",
								}},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "krossboard-config-vol",
										MountPath: "/etc/caddy/Caddyfile",
										SubPath:   "Caddyfile",
										ReadOnly:  true,
									},
									{
										Name:      "caddy-config-vol",
										MountPath: "/root/.caddy",
										ReadOnly:  true,
									},
								},
							},
							{
								Image:   kbProcessorImage,
								Name:    "krossboard-api",
								Command: []string{"/app/krossboard-data-processor", "api"},
								Ports: []corev1.ContainerPort{{
									ContainerPort: 1519,
									Name:          "krossboard-api",
								}},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "krossboard-data-vol",
										MountPath: kbDataDir,
										ReadOnly:  true,
									},
									{
										Name:      "krossboard-run-vol",
										MountPath: kbRuntimeDir,
										ReadOnly:  true,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "KROSSBOARD_ROOT_DIR",
										Value: kbDataDir,
									},
									{
										Name:  "KROSSBOARD_RUN_DIR",
										Value: kbRuntimeDir,
									},
									{
										Name:  "KROSSBOARD_RAWDB_DIR",
										Value: kbDataDir,
									},
									{
										Name:  "KROSSBOARD_HISTORYDB_DIR",
										Value: kbDataDir,
									},
								},
							},
							{
								Image:   kbProcessorImage,
								Name:    "krossboard-consolidator",
								Command: []string{"/bin/bash", "-c", "--"},
								Args:    []string{"while true; do /app/krossboard-data-processor consolidator; sleep 300; done;"},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "krossboard-data-vol",
										MountPath: kbDataDir,
										ReadOnly:  false,
									},
									{
										Name:      "krossboard-run-vol",
										MountPath: kbRuntimeDir,
										ReadOnly:  false,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "KROSSBOARD_SELECTED_CLUSTER_NAMES",
										Value: selectedClusterNames,
									},
									{
										Name:  "KROSSBOARD_ROOT_DIR",
										Value: kbDataDir,
									},
									{
										Name:  "KROSSBOARD_RAWDB_DIR",
										Value: kbDataDir,
									},
									{
										Name:  "KROSSBOARD_HISTORYDB_DIR",
										Value: kbDataDir,
									},
									{
										Name:  "KROSSBOARD_RUN_DIR",
										Value: kbRuntimeDir,
									},
								},
							},
							{
								Image:   kbProcessorImage,
								Name:    "krossboard-credentials-handler",
								Command: []string{"/bin/bash", "-c", "--"},
								Args:    []string{"while true; do /app/krossboard-data-processor cluster-credentials-handler; sleep 300; done;"},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "krossboard-secrets-vol",
										MountPath: kbK8sSecretsDir,
										ReadOnly:  true,
									},
									{
										Name:      "krossboard-creds-vol",
										MountPath: kbCredsDir,
										ReadOnly:  false,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name:  "KROSSBOARD_CREDENTIALS_DIR",
										Value: kbCredsDir,
									},
									{
										Name:  "KUBECONFIG",
										Value: fmt.Sprintf("%s/kubeconfig", kbK8sSecretsDir),
									},
								},
							},
						},
						koaContainers...,
					),
					Volumes: []corev1.Volume{
						{
							Name: "caddy-config-vol",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "krossboard-config-vol",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "krossboard-config",
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "Caddyfile",
											Path: "Caddyfile",
										},
									},
								},
							},
						},
						{
							Name: "krossboard-secrets-vol",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: kbSecretName,
									Items: []corev1.KeyToPath{
										{
											Key:  "kubeconfig",
											Path: "kubeconfig",
										},
									},
								},
							},
						},
						{
							Name: "krossboard-data-vol",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: kbPersistentVolumeClaim,
									ReadOnly:  false,
								},
							},
						},
						{
							Name: "krossboard-creds-vol",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "krossboard-run-vol",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	// If UseGKEIdentity is enabled, target compatinle node pools
	if m.Spec.UseGKEIdentity {
		dep.Spec.Template.Spec.NodeSelector = map[string]string{
			"iam.gke.io/gke-metadata-server-enabled": "true",
		}
	}

	_ = ctrl.SetControllerReference(m, dep, kbReconciler.Scheme)

	return dep
}
