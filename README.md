<!-- vscode-markdown-toc -->
- [Overview](#overview)
- [Installation](#installation)
  - [Create a namespace for Krossboard](#create-a-namespace-for-krossboard)
  - [Create a secret for KUBECONFIG](#create-a-secret-for-kubeconfig)
  - [Deploy Krossboard Operator](#deploy-krossboard-operator)
  - [Deploy a Krossboard Instance](#deploy-a-krossboard-instance)
- [Day2 Documentation](#day2-documentation)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

# Overview

[Krossboard](https://www.krossboard.app/) is a multi-cluster and cross-distribution Kubernetes usage accounting and analytics software. 

> Learn about [Krossboard features](https://github.com/2-alchemists/krossboard#overview).

![](krossboard-architecture-overview.png)

Krossboard Operator provides custom resources (CRD) along with an operator to deploy and manage instances of Krossboard as Kubernetes pods.

The `Krossboard` CRD (see [krossboard.yaml](https://github.com/2-alchemists/krossboard-kubernetes-operator/blob/main/config/latest/krossboard.yaml) allows to define the container images of each Krossboard components: krossboard-api, krossboard-ui, krossboard-consolidator, krossboard-kubeconfig-handler, kube-opex-analytics, etc.

Each instance of Krossboard allows to track the usage of a set of Kubernetes clusters listed in a KUBECONFIG secret:

* Secret Name: `krossboard-secrets`
* Secret Key: `kubeconfig` (base64-encoded KUBECONFIG resource).

A different secret can be used (instead of `krossboard-secrets`). In this case, you must set the parameter `krossboardSecretName` of the Krossboard CRD with the name of the selected secret (see [krossboard.yaml](https://github.com/2-alchemists/krossboard-kubernetes-operator/blob/main/config/latest/krossboard.yaml)).


# Installation

## <a name='CreateanamespaceforKrossboard'></a>Create a namespace for Krossboard

Krossboard Operator is namespace-scoped, it must be deployed in the `krossboard` namespace.

```
kubectl create namespace krossboard
```

## <a name='CreateasecretforKUBECONFIG'></a>Create a secret for KUBECONFIG
The Kubernetes clusters handled by Krossboard Operator are set via a secret within the installation namespace (Secret Name: `krossboard-secrets`, Secret Key: `kubeconfig`)

> Learn how to create a [KUBECONFIG with minimal permissions for Krossboard](./docs/create-kubeconfig-with-minimal-permissions.md).


**Create a secret from a single KUBECONFIG resource**

Given a KUBECONFIG resource, you can create a secret for Krossboard Operator as follows. 

```bash
kubectl -n krossboard \
    create secret generic krossboard-secrets  \
    --from-file=kubeconfig=/path/to/kubeconfig \
    --type=Opaque
```

Replace `/path/to/kubeconfig` with the path of the KUBECONFIG file.

**Create a secret from multiple KUBECONFIG resources**

If you have several KUBECONFIG resources, you can merge them by proceeding as hereafter.

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

## <a name='DeployKrossboardOperator'></a>Deploy Krossboard Operator
The following command deploy the latest version of Korssboard Operator.

```bash
kubectl apply -f config/releases/latest/krossboard/krossboard-kubernetes-operator.yaml
```

## <a name='DeployKrossboardOperator'></a>Deploy a Krossboard Instance
The following command deploy the latest version of Krossboard.

```bash
kubectl -n krossboard apply -k config/releases/latest/krossboard/
```

# Day2 Documentation

* https://krossboard.app/
* [Krossboard Enterprise Support](https://krossboard.app/#pricing) 
