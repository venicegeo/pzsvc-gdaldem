#!/bin/bash -ex

pushd `dirname $0`/.. > /dev/null
root=$(pwd -P)
popd > /dev/null

# Run the test
newman -c $root/postman/pzsvc-gdaldem-black-box-tests.json.postman_collection -e $root/postman/pzsvc-gdaldem.json.postman_environment.cf -s
