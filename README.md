![](krossboard-architecture-overview.png)
---

# Overview
`krossboard-kubernetes-operator` provides an operator and custom resources (CRD) to manage one or more  [Krossboard](https://www.krossboard.app/) instances within a Kubernetes cluster.

The main CRD (`Krossboard`) allows to deploy an instance of Krossboard: API, UI and kube-opex-analytics instances based on a KUBECONFIG defined via a secret (name: `krossboard`, key: `kubeconfig`).

Key features of Krossboard:

* **Multi-Kubernetes Data Collection**: It periodically collects raw metrics related to containers, pods and nodes from each Kubernetes cluster it handles. The built-in data collection period is 5 minutes.
* **Powerful Analytics Processing**: It internally processes raw metrics to produce insightful Kubernetes usage accounting and analytics metrics. The analytics data are tracked on a hourly-basis, per namespace, per cluster, and globally.
* **Insightful Usage Accounting**: It periodically processes usage accounting, per namespace and per cluster. By the default, the UI displays accounting the following periods without any additioanl configuration: daily accounting for the last 14 days, monthly for the ast 12 months.
* **REST API**: This exposes the generated analytics data it generates to third-party systems.
* **Easy to deploy**: The Krossboard operator is deployable in a couple of minutes.

# Getting Started

## Create a namespace for krossboard

Krossboard Kubernetes Operator is namespaced-scoped, and is expected to be deployed in a namespace named `krossboard`.

```
kubectl create namespace krossboard
```

## Create a secret for KUBECONFIG
The Kubernetes clusters handled by the Krossboard operator shall be defined through a secret (name: `krossboard-secrets`, key: `kubeconfig`).

> Krossboard currently support KUBECONFIG where Kubernetes users credentials are based on either token, client certificate, or username + password.

### Create a secret from a KUBECONFIG file

Given a KUBECONFIG file, you can create such a secret as follows (replace `/path/to/kubeconfig` with the path towards the KUBECONFIG file). 

```
kubectl -n krossboard \
    create secret generic krossboard-secrets  \
    --from-file=kubeconfig=/path/to/kubeconfig \
    --type=Opaque
```

### Create a secret from multiple KUBECONFIG files

If you have several KUBECONFIG files, you can merge them by proceeding as described hereafter.

Set an environment variable with a comma-seperated list of KUBECONFIG files.

```
export KB_KUBECONFIG=/path/to/kubeconfig1;/path/to/kubeconfig2;...
```

Create a secret file with the resulting merged KUBECONFIG.
```
kubectl -n krossboard \
    create secret generic krossboard-secrets \
    --from-file=kubeconfig=<(KUBECONFIG=$KB_KUBECONFIG kubectl config view --raw) \
    --type=Opaque --dry-run=client -oyaml > ./krossboard-secrets.yaml
```

Create a secret for Krossboard with the merged KUBECONFIG.
```
kubectl apply -f ./krossboard-secrets.yaml
```

## Deploy Krossboard Kubernetes Operator
The following command deploy the latest version of the operator.

```
kubectl apply -k config/latest/
```

# Additioan resources

* https://krossboard.app/
* [Krossboard Enterprise Support](https://krossboard.app/#pricing) 
