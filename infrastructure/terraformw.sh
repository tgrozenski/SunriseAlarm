#!/bin/bash
set -e
cd "$(git rev-parse --show-toplevel)/infrastructure/terraform"

action=apply
if [ "$1" == "destroy" ]; then
  action=destroy
fi

echo "Performing on GCP resources..."
cd gcp
terraform $action

echo "Copying google-services.json to app dir"
terraform output -raw firebase_android_config | base64 -d > "$(git rev-parse --show-toplevel)/app/google-services.json"

echo "Performing on AWS resources..."
cd ../aws
terraform $action

echo "Success"
