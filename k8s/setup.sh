#!/bin/bash
# This script is used to setup the environment
PROJECT_ID="mofilabs-batch-1"
ZONE="us-central1-a"
REGION="us-central1"
CLUSTER_NAME="batch"
NUM_NODES="3"
MACHINE_TYPE="c2-standard-30"

# gcloud project create $PROJECT_NAME
echo "Set up gcloud config"
echo "Project ID: ${PROJECT_ID}"
echo "Zone: ${ZONE}"
echo "Region: ${REGION}"
echo ""
gcloud config set core/project ${PROJECT_ID}
gcloud config set compute/zone ${ZONE}
gcloud config set compute/region ${REGION}

echo ""
echo "Creating Kubernetes Cluster"

gcloud beta container --project ${PROJECT_ID} \
    clusters create ${CLUSTER_NAME} \
    --zone ${ZONE} \
    --release-channel "regular" \
    --machine-type ${MACHINE_TYPE} \
    --num-nodes ${NUM_NODES} \
    --addons HorizontalPodAutoscaling,HttpLoadBalancing,GcePersistentDiskCsiDriver,GcpFilestoreCsiDriver \
    --enable-managed-prometheus \
    --enable-autoscaling --min-nodes=1 --max-nodes=6 \
    --logging=SYSTEM,WORKLOAD \
    --logging-variant=MAX_THROUGHPUT

echo ""
echo "Cluster Created: ${CLUSTER_NAME}"

echo ""
echo "Create PVC"
kubectl apply -f pvc.yaml
echo "PVC Created"
echo ""