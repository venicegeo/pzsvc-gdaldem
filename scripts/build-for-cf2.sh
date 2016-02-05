#!/bin/bash

echo "Compiling for linux..."
GOOS=linux GOARCH=amd64 go build .

echo "Constructing Docker image"
docker build -t chambbj/pzsvc-gdaldem .
docker push chambbj/pzsvc-gdaldem

echo "Cleaning up..."
rm pzsvc-gdaldem
