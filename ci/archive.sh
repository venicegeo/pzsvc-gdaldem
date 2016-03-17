#!/bin/bash -ex

pushd `dirname $0`/.. > /dev/null
root=$(pwd -P)
popd > /dev/null

export GOPATH=$root/gogo
mkdir -p $GOPATH

###

export GO15VENDOREXPERIMENT=1

go get -v github.com/venicegeo/pzsvc-gdaldem

go install -v github.com/venicegeo/pzsvc-gdaldem


###

src=$GOPATH/bin/$APP

# gather some data about the repo
source $root/ci/vars.sh

# stage the artifact for a mvn deploy
mv $src $root/$APP.$EXT
