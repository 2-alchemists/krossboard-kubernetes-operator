# Create a KUBECONFIG secret for target Kubernetes

The Kubernetes clusters handled by Krossboard Operator are set via a secret within the installation namespace (Secret Name: `krossboard-secrets`, Secret Key: `kubeconfig`)

> Learn how to create a [KUBECONFIG with minimal permissions for Krossboard](./docs/create-kubeconfig-with-minimal-permissions.md).

This page provides a general guide to create a secret based on one or more KUBECONFIG resources.

When you have several KUBECONFIG resources, the idea of the procedure described hereafyer is to create a merged version before creating the secret.

First, set an environment variable with a colon-seperated list of KUBECONFIG files.

```bash
export KB_KUBECONFIG=/path/to/kubeconfig1:/path/to/kubeconfig2:...
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
