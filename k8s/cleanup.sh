#!/bin/bash
# This script is used to cleanup the environment
PROJECT_ID="mofilabs-batch-1"
ZONE="us-central1-a"
REGION="us-central1"
CLUSTER_NAME="batch"
NUM_NODES="3"
MACHINE_TYPE="c2-standard-30"


echo "Delete PVC"
kubectl delete -f pvc.yaml
echo "Deleting Cluster"
gcloud container clusters delete ${CLUSTER_NAME} --zone ${ZONE} --project ${PROJECT_ID} --quiet