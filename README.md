<!-- vscode-markdown-toc -->
- [Overview](#overview)
- [Key Features](#key-features)
- [Installation](#installation)
  - [Create a namespace for Krossboard](#create-a-namespace-for-krossboard)
  - [Create a secret for KUBECONFIG](#create-a-secret-for-kubeconfig)
  - [Deploy Krossboard Operator](#deploy-krossboard-operator)
- [Day2 Documentation](#day2-documentation)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

# Overview

[Krossboard](https://www.krossboard.app/) is a multi-cluster and cross-distribution Kubernetes usage accounting and analytics software.

Krossboard Operator (`krossboard-kubernetes-operator`) provides an operator and custom resources (CRD) to deploy and manage instances of Krossboard as Kubernetes pods.

The `Krossboard` CRD (see [krossboard.yaml](https://github.com/2-alchemists/krossboard-kubernetes-operator/blob/main/config/latest/krossboard.yaml) allows to define the container images of each Krossboard components: krossboard-api, krossboard-ui, krossboard-consolidator, krossboard-kubeconfig-handler, kube-opex-analytics, etc.

Each instance of Krossboard allows to track the usage of a set of Kubernetes clusters listed in a KUBECONFIG resource.

The KUBECONFIG resource is set via a secret within the installation namespace:

* Secret Name: `krossboard-secrets` (by default).
* Secret Key: `kubeconfig`

>  You can use a different secret while setting the `krossboardSecretName` accordingly (see [krossboard.yaml](https://github.com/2-alchemists/krossboard-kubernetes-operator/blob/main/config/latest/krossboard.yaml)).


![](krossboard-architecture-overview.png)

# Key Features

Highlight of Krossboard features:

* **Multi-Kubernetes Data Collection**: Krossboard periodically collects raw metrics related to containers, pods and nodes from each Kubernetes cluster it handles. The built-in data collection period is 5 minutes.
* **Powerful Analytics Processing**: Krossboard internally processes raw metrics to produce insightful Kubernetes usage accounting and analytics metrics. The analytics data are tracked on a hourly-basis, per namespace, per cluster, and globally.
* **Insightful Usage Accounting**: Krossboard periodically processes usage accounting, per namespace and per cluster. By the default, the UI displays accounting the following periods without any additioanl configuration: daily accounting for the last 14 days, monthly for the ast 12 months.
* **REST API**: This exposes the generated analytics data it generates to third-party systems.
* **Easy to deploy**: The Krossboard operator is deployable in a couple of minutes.

# Installation

## <a name='Createanamespaceforkrossboard'></a>Create a namespace for Krossboard

Krossboard Operator is namespace-scoped, it must be deployed in the `krossboard` namespace.

```
kubectl create namespace krossboard
```

## <a name='CreateasecretforKUBECONFIG'></a>Create a secret for KUBECONFIG
The Kubernetes clusters handled by Krossboard Operator are set via a secret within the installation namespace:

* Secret Name: `krossboard-secrets` (by default)
* Secret Key: `kubeconfig`

The `kubeconfig` key must be set with a base64-encoded KUBECONFIG content. The associated credentials can be of type token, client certificate, or username + password.

[Learn how to create a KUBECONFIG with minimal permissions for Krossboard](./docs/create-kubeconfig-with-minimal-permissions.md).


**Create a secret from a single KUBECONFIG file**

Given a KUBECONFIG file, you can create a secret for Krossboard Operator as follows. 

```bash
kubectl -n krossboard \
    create secret generic krossboard-secrets  \
    --from-file=kubeconfig=/path/to/kubeconfig \
    --type=Opaque
```

Replace `/path/to/kubeconfig` with the path of the KUBECONFIG file.

**Create a secret from multiple KUBECONFIG file**

If you have several KUBECONFIG files, you can merge them by proceeding as hereafter.

Set an environment variable with a comma-seperated list of KUBECONFIG files.

```bash
export KB_KUBECONFIG=/path/to/kubeconfig1;/path/to/kubeconfig2;...
```

Generate a secret file with the resulting KUBECONFIG.

```bash
kubectl -n krossboard \
    create secret generic krossboard-secrets \
    --from-file=kubeconfig=<(KUBECONFIG=$KB_KUBECONFIG kubectl config view --raw) \
    --type=Opaque --dry-run=client -oyaml > ./krossboard-secrets.yaml
```

Review the secret file and apply it the secret file.

```bash
kubectl apply -f ./krossboard-secrets.yaml
```

## <a name='DeployKrossboardKubernetesOperator'></a>Deploy Krossboard Operator
The following command deploy the latest version of the operator.

```bash
kubectl apply -k config/latest/
```

# Day2 Documentation

* https://krossboard.app/
* [Krossboard Enterprise Support](https://krossboard.app/#pricing) 
