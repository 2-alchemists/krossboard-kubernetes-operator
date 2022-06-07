#!/bin/bash

set -e
set -u

docker login

make undeploy || true

make build docker-build docker-push
make deploy

# get SA token
# oc -n kube-opex-analytics get secret $(oc -n kube-opex-analytics get sa kube-opex-analytics -ojsonpath='{.secrets[0].name}') -ojsonpath='{.data.token}'  | base64 -d > token

kubectl -n krossboard create secret generic krossboard-secrets  \
    --from-file=kubeconfig=$HOME/.kube/config-sa \
    --type=Opaque

oc apply -k config/latest/

docker system prune -f -a