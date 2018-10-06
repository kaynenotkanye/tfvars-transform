#!/bin/bash
rm ./bin/tfvars-transform
echo "Building binary"
GOOS=linux GOARCH=amd64 go build -o ./bin/tfvars-transform tfvars-transform.go
