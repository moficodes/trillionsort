#!/bin/bash
echo "Creating the generate job..."
kubectl apply -f job-generate.yaml

echo "Waiting for the job to complete..."
kubectl wait --for=condition=Complete job/generate

echo "Creating the join job..."
kubectl apply -f job-join.yaml

echo "Waiting for the job to complete..."
kubectl wait --for=condition=Complete job/join

echo "All done!"