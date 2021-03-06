#!/bin/sh

pushd `dirname $0` > /dev/null
base=$(pwd -P)
popd > /dev/null

export GOPATH=$base/gogo
mkdir -p $GOPATH

###

export GO15VENDOREXPERIMENT=1

go get -v github.com/venicegeo/pzsvc-gdaldem

go test -v $(go list github.com/venicegeo/pzsvc-gdaldem/... | grep -v /vendor/)

go install -v github.com/venicegeo/pzsvc-gdaldem

###

exe=$GOPATH/bin/pzsvc-gdaldem

# gather some data about the repo
source $base/vars.sh

# do we have this artifact in s3? If not, upload it.
aws s3 ls $S3URL || aws s3 cp $exe $S3URL

# success if we have an artifact stored in s3.
aws s3 ls $S3URL
