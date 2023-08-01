#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

THIS_DIRECTORY=$(dirname "${BASH_SOURCE}")

PROJECT_DIRECTORY=${THIS_DIRECTORY}/../..

K8S_DESTINATION_LOCATION_API=${PROJECT_DIRECTORY}/k8s/client

DOCKER_ARGUMENTS="run --rm=true --entrypoint /bin/bash -it -v ${PWD}:/go/src/github.com/venezia/k8s-gorilla-session-store -w /go/src/github.com/venezia/k8s-gorilla-session-store"
DOCKER_IMAGE=quay.io/venezia/k8s-golang-code-generator:1.27.3

echo "Ensuring Astro API Destination Directory ( ${K8S_DESTINATION_LOCATION_API} ) Exists..."
mkdir -p ${K8S_DESTINATION_LOCATION_API}

echo docker $DOCKER_ARGUMENTS $DOCKER_IMAGE \
       /go/code-generator/generate-groups.sh \
       "all" \
       github.com/venezia/k8s-gorilla-session-store/k8s/client \
       github.com/venezia/k8s-gorilla-session-store/k8s/api \
       session:v1alpha1 \
       --output-base "/go/src" \
       --go-header-file assets/k8s/boilerplate.txt

docker $DOCKER_ARGUMENTS $DOCKER_IMAGE \
       /go/code-generator/generate-groups.sh \
       "all" \
       github.com/venezia/k8s-gorilla-session-store/k8s/client \
       github.com/venezia/k8s-gorilla-session-store/k8s/api \
       session:v1alpha1 \
       --output-base "/go/src" \
       --go-header-file assets/k8s/boilerplate.txt
