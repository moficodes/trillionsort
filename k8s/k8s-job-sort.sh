#!/bin/bash
echo "Creating the split job..."
kubectl apply -f job-split.yaml

echo "Waiting for the job to complete..."
kubectl wait --for=condition=Complete job/split

echo "Creating the sort job..."
kubectl apply -f job-sort.yaml

echo "Waiting for the job to complete..."
kubectl wait --for=condition=Complete job/sort

echo "Creating the merge job..."
kubectl apply -f job-merge.yaml

echo "Waiting for the job to complete..."
kubectl wait --for=condition=Complete job/externalsort

echo "All done!"