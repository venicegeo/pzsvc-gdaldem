[![GoDoc](https://godoc.org/github.com/venicegeo/pzsvc-gdaldem?status.svg)](https://godoc.org/github.com/venicegeo/pzsvc-gdaldem)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/venicegeo/pzsvc-gdaldem/blob/master/LICENSE)

# pzsvc-gdaldem

Needs updates

# Install

[GDAL](http://www.gdal.org/) is a commonly used library and set of tools for raster processing. The [gdaldem](http://www.gdal.org/gdaldem.html) program provides tools for working with elevation rasters, e.g., creating hillshades. We have created a [Dockerfile](https://github.com/venicegeo/dockerfiles/blob/master/gdaldem/Dockerfile) that generates a Docker image consisting of gdaldem. It can be built with the following commands

```console
$ git clone https://github.com/venicegeo/dockerfiles/gdaldem
$ cd gdaldem
$ docker build -t venicegeo/gdaldem .
```

This in turn serves as the base image for our microservice, which is written in Go.

`pzsvc-gdaldem` uses [Glide](https://github.com/Masterminds/glide) to manage its dependencies. Assuming you are on a Mac OS X, Glide can be easily installed via [Homebrew](https://github.com/Homebrew/homebrew) (alternative installation instruction can be found on the Glide webpage).

```console
$ brew install glide
```

We also make use of [Go 1.5's vendor/ experiment](https://medium.com/@freeformz/go-1-5-s-vendor-experiment-fd3e830f52c3#.ueuy8ao53), so you'll need to make sure you are running Go 1.5+.

Installing `pzsvc-gdaldem` is as simple as cloning the repo, installing dependencies, and running the provided `build.sh`.

```console
$ git clone https://github.com/venicegeo/pzsvc-gdaldem
$ cd pzsvc-gdaldem
$ glide install
$ scripts/build.sh
```

The build script will first compile the Go code in a temporary container. The resulting static Go binary is then copied into our `venicegeo/pzsvc-gdaldem` image during the `docker build` step.

Finally, the service is started on port 8080, mounting your `~/.aws/credentials` to the image with `run.sh`.

```console
$ scripts/run.sh
```
