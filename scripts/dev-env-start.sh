#!/usr/bin/env bash

podman machine start
podman-compose -f ./deployments/docker-compose-dev.yml up -d
