#!/bin/bash
set -e

cd "$(git rev-parse --show-toplevel)/infrastructure/terraform"

echo "Applying GCP resources..."
cd gcp
terraform apply "$@"

echo "Copying google-services.json to app dir"
terraform output -raw firebase_android_config | base64 -d > "$(git rev-parse --show-toplevel)/app/google-services.json"

echo "Applying AWS resources..."
cd ../aws
terraform apply "$@"

echo "Success"
