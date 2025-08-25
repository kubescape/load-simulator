#!/bin/bash
# This script ensures ko is available for this project
# Set PATH to include the current Go version's bin directory
export PATH="$(go env GOBIN):$PATH"
KO_DOCKER_REPO=kind.local ko apply -f config/
