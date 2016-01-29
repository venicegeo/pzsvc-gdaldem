[![GoDoc](https://godoc.org/github.com/venicegeo/pzsvc-gdaldem?status.svg)](https://godoc.org/github.com/venicegeo/pzsvc-gdaldem)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/venicegeo/pzsvc-gdaldem/blob/master/LICENSE)

# pzsvc-gdaldem

[GDAL](http://www.gdal.org/) is a commonly used library and set of tools for raster processing. The [gdaldem](http://www.gdal.org/gdaldem.html) program provides tools for working with elevation rasters, e.g., creating hillshades. The purpose of this application is to expose `gdaldem` as an API endpoint.


## Install

Go 1.5+ is required. You can download it [here](https://golang.org/dl/).

If you have not already done so, make sure you've setup your Go [workspace](https://golang.org/doc/code.html#Workspaces) and set the necessary environment [variables](https://golang.org/doc/code.html#GOPATH).

For `pzsvc-gdaldem` to function properly, GDAL must be installed on your system. Our manifest.yml file specifies a custom buildpack to ensure that GDAL is available on Cloud Foundry. For local operation, follow installation instructions for your system, e.g., `brew install gdal` on Mac OS X.

To install `pzsvc-gdaldem`, simply clone the repository and issue a `go install` command.

```console
$ git clone https://github.com/venicegeo/pzsvc-gdaldem
$ cd pzsvc-gdaldem
$ go install github.com/venicegeo/pzsvc-gdaldem
```

Assuming `$GOPATH/bin` is on your `$PATH`, the service can easily be started on port 8080.

```console
$ pzsvc-gdaldem
```

The following curl command should return `Hi!`.

```console
$ curl -X GET -H "Cache-Control: no-cache" -H "Postman-Token: 73e58ef4-44ba-8b40-6868-72ff404b41dd" 'http://localhost:8080'
```

## Example

While this one will download the sample `elevation.tif` from S3, create a hillshade with default parameters, and upload the result to S3 as `hillshade.tif`.

```console
curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -H "Postman-Token: d5eb710b-cf94-c8b1-db29-98d4df633895" -d '{
    "source": {
        "bucket": "venicegeo-sample-data",
        "key": "raster/elevation.tif"
    },
    "function":"hillshade",
    "destination": {
        "bucket": "venicegeo-sample-data",
        "key": "temp/hillshade.tif"
    }
}' 'http://localhost:8080/gdaldem'
```

## Deploying

When deployed, `localhost:8080` is replaced with `pzsvc-gdaldem.cf.piazzageo.io`.

All commits to master will be pushed through the VeniceGeo DevOps infrastructure, first triggering a build in Jenkins and, upon success, pushing the resulting binaries to Cloud Foundry.
