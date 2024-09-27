#!/usr/bin/env bash

/usr/bin/podman build -t alert-jira-gateway:latest .

podman  tag   \
        alert-jira-gateway:latest \
        default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/alert-jira-gateway:v0.1

podman  push  \
        default-route-openshift-image-registry.apps.cid-scp0.sf-rz.de/scp-images/alert-jira-gateway:v0.1



