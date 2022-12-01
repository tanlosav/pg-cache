#!/usr/bin/env bash

podman-compose -f ./deployments/docker-compose-dev.yml down
podman machine stop
