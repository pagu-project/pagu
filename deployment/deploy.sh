#!/bin/bash

set -e

# Ensure Database is up and running
docker compose -f ./deployment/docker-compose.db.yml up -d


export DEPLOY_TAG="latest"

TAG=$(git describe --tags --exact-match 2> /dev/null) || echo ""

if [[ -n "$TAG" ]]; then
    export DEPLOY_TAG="stable"
fi

echo "Building ${DEPLOY_TAG} version..."
docker build -t pagu:${DEPLOY_TAG} -f ./deployment/Dockerfile .

docker compose -p pagu_${DEPLOY_TAG} -f ./deployment/docker-compose.yml up -d

## Some cleanup
echo "Cleanup..."

docker builder prune -f
docker image prune -f
