#!/bin/bash
# ------------------------------------------------------------------------ #
# Copyright (c) 2022 2Alchemists SAS                                       #
#                                                                          #
# This file is part of Krossboard (https://krossboard.app/).               #
#                                                                          #
# The tool is distributed in the hope that it will be useful,              #
# but WITHOUT ANY WARRANTY; without even the implied warranty of           #
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the            #
# Krossboard terms of use: https://krossboard.app/legal/terms-of-use/      #
#--------------------------------------------------------------------------#

set -e

if ! command -v gcloud &> /dev/null; then
  echo "\e[31m[ERROR] gcloud sdk could not be found, please install it => https://cloud.google.com/sdk\e[0m"
  exit 1
fi

if [ -z "$GCP_PROJECT" ]; then
  GCP_PROJECT=$(gcloud config get-value project)
  echo "===> GCP_PROJECT variable is not set..."
  echo -e "\e[35m     ↳ Using current project => $GCP_PROJECT\e[0m"
fi

set -u

GCLOUD_CMD="gcloud --project=$GCP_PROJECT"
echo "==> Configuring IAM permissions for Krossboard..."
KB_SA_NAME='krossboard-sa'
KB_SA_EMAIL=$($GCLOUD_CMD iam service-accounts list --filter="name:$KB_SA_NAME@$GCP_PROJECT" --format="value(EMAIL)")
if [ "$KB_SA_EMAIL" != "" ]; then
  echo -e "\e[35m    ↳ Using existing service account\e[0m"
  echo -e "\e[35m      ↳ ${KB_SA_NAME} => $KB_SA_EMAIL\e[0m"
else
  echo -e "\e[35mCreating a GCP service account ${KB_SA_NAME}...\e[0m"
  $GCLOUD_CMD iam service-accounts create $KB_SA_NAME --display-name $KB_SA_NAME
  retry=0
  until [ "$retry" -ge 15 ]; do
    KB_SA_EMAIL=$($GCLOUD_CMD iam service-accounts list --filter="name:$KB_SA_NAME@$GCP_PROJECT" --format="value(EMAIL)")
    if [ -z "$KB_SA_EMAIL" ]; then
      echo -e "\e[35mWaiting the service account to become ready...\e[0m"
      sleep 1
    else
      echo -e "\e[35mService account email => $KB_SA_EMAIL\e[0m"
      $GCLOUD_CMD projects add-iam-policy-binding "$GCP_PROJECT" --member="serviceAccount:$KB_SA_EMAIL" --role='roles/container.viewer'
      break
    fi
    sleep 1
    retry=$((retry+1)) 
  done
  if [ -z "$KB_SA_EMAIL" ]; then
    echo -e "\e[31m[ERROR] Failed getting created service account\e[0m"
    exit 1
  fi
fi
