/*
Copyright 2015-2016, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
pzsvc-gdaldem provides an endpoint for accepting PDAL requests.

Examples

  $ curl -v -X POST -H "Content-Type: application/json" \
    -d '{"source":{"bucket":"venicegeo-sample-data","key":"pointcloud/samp11-utm.laz"},"function":"info"}' http://hostIP:8080/pdal
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Hi!")
	})

	// Setup the PDAL service.
	// router.POST("/gdaldem", handlers.GdalDemHandler)

	var defaultPort = os.Getenv("PORT")
	log.Println("PORT is ", defaultPort)
	if defaultPort == "" {
		defaultPort = "8080"
	}

	log.Println("Starting on ", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, router); err != nil {
		log.Fatal(err)
	}
}
