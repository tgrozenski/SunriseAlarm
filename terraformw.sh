#!/bin/bash
set -e

cd "$(git rev-parse --show-toplevel)/infrastructure/terraform"

echo "Applying GCP resources..."
cd gcp
terraform apply "$@"

echo "Applying AWS resources..."
cd ../aws
terraform apply "$@"

echo "Success"
