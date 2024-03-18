



# Create KUBECONFIG wih Minimal Permissions for Krossboard

<!-- vscode-markdown-tkubectl -->
- [Create KUBECONFIG wih Minimal Permissions for Krossboard](#create-kubeconfig-wih-minimal-permissions-for-krossboard)
  - [Overview](#overview)
  - [Create RBAC Resources (Service Account, ClusterRole, ClusterRoleBinding)](#create-rbac-resources-service-account-clusterrole-clusterrolebinding)
  - [Get the Service Account Token](#get-the-service-account-token)
  - [Create the KUBECONFIG Resource](#create-the-kubeconfig-resource)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-tkubectl -->

## <a name='Overview'></a>Overview
This document describes step-by-step how to create a KUBECONFIG resource with minimal RBAC permissions for Krossboard.

The minimal RBAC permissions needed for Krossboard are the same as required for [kube-opex-analytics](https://github.com/rchakode/kube-opex-analytics). Therefore, we rely on the same [set of RBAC resources](https://github.com/rchakode/kube-opex-analytics/blob/main/manifests/kustomize/resources/kube-opex-analytics-rbac.yaml).

These RBAC resources enable permissions to retrieve raw usage metrics related to nodes and pods from Kubernetes API. They resources includes a Service Account, a ClusterRole along with a ClusterRoleBinding binding ClusterRole and the Service Account.


## <a name='CreateRBACResourcesServiceAccountClusterRoleClusterRoleBinding'></a>Create RBAC Resources (Service Account, ClusterRole, ClusterRoleBinding)

The below two commands create a namespace `kube-opex-analytics` along with the needed RBAC resources. 


```bash
kubectl create ns kube-opex-analytics
kubectl apply -f https://raw.githubusercontent.com/rchakode/kube-opex-analytics/main/manifests/kustomize/resources/kube-opex-analytics-rbac.yaml
```

The created RBAC resources include the following: 

* Service Account: `kube-opex-analytics`
* ClusterRole: `kube-opex-analytics`
* ClusterRolebinding: `kube-opex-analytics`

## <a name='GettheServiceAccountToken'></a>Get the Service Account Token

The following command outputs the token associated to the service account `kube-opex-analytics`.

```bash
kubectl -n kube-opex-analytics get secret kube-opex-analytics -ojsonpath='{.data.token}'  | base64 -d
```


## <a name='CreatetheKUBECONFIGResource'></a>Create the KUBECONFIG Resource

You can create a KUBECONFIG resource for Krossboard using the template below. 

Make sure to update the following parameters:

* Replace `https://your-k8s-cluster-api:6443` with the URL of the API of the Kubernetes cluster.
* Replace all occurrences of `<cluster-name>` with the name of the cluster.
* Replace `<service-account-token-here>` with the service account token extracted previously.
* Replace `<cluster-cacert>` with a base64-encoded value of the certificate authority data.

```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: <K8s-CLUSTER-CACERT>
    server: https://K8s-CLUSTER-API:6443
  name: <K8s-CLUSTER-NAME>
contexts:
- context:
    cluster:  <K8s-CLUSTER-NAME>
    user:  <K8s-CLUSTER-NAME>_serviceaccount_kube-opex-analytics
  name:  <K8s-CLUSTER-NAME>
current-context:  <K8s-CLUSTER-NAME>
kind: Config
preferences: {}
users:
- name:  <K8s-CLUSTER-NAME>_serviceaccount_kube-opex-analytics
  user:
    token: <KUBE-OPEX-ANALYTICS-SERVICE-ACCOUNT-TOKEN>
```

