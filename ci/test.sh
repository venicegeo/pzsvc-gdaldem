#!/bin/bash -ex

pushd `dirname $0`/.. > /dev/null
root=$(pwd -P)
popd > /dev/null

export GOPATH=$root/gogo
mkdir -p $GOPATH

###
export GO15VENDOREXPERIMENT=1

go get -v github.com/venicegeo/pzsvc-gdaldem

go test -v $(go list github.com/venicegeo/pzsvc-gdaldem/... | grep -v /vendor/)

###
