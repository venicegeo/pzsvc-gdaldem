#!/bin/bash -ex

pushd `dirname $0` > /dev/null
base=$(pwd -P)
popd > /dev/null

# Gather some data about the repo
source $base/vars.sh

# Run the test
newman -c $base/../postman/pzsvc-gdaldem-black-box-tests.json.postman_collection -e $base/../postman/pzsvc-gdaldem.json.postman_environment.cf -s
