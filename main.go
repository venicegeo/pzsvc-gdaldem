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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-gdaldem/handlers"
	"github.com/venicegeo/pzsvc-sdk-go/servicecontroller"
)

type appHandler func(http.ResponseWriter, *http.Request) *handlers.AppError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		if awsErr, ok := e.Error.(awserr.Error); ok {
			e.Message = awsErr.Message()
		}

		if err := json.NewEncoder(w).Encode(e); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(e.Code)
	}
}

func main() {

	m := servicecontroller.ResourceMetadata{
		Name: "pzsvc-gdaldem",
		URL:  "http://pzsvc-gdaldem.cf.piazzageo.io/gdaldem",
		//		URL:              "http://localhost:8080/gdaldem",
		Description:      "GDAL's gdaldem: Tools to analyze and visualize DEMs",
		Method:           "POST",
		RequestMimeType:  servicecontroller.ContentTypeJSON,
		ResponseMimeType: servicecontroller.ContentTypeJSON,
	}
	if err := servicecontroller.RegisterService(m); err != nil {
		log.Println(err)
	}

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Hi!")
	})

	// Setup the PDAL service.
	router.Handler("POST", "/gdaldem", appHandler(handlers.GdalDemHandler))

	var defaultPort = os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}

	log.Println("Starting on port ", defaultPort)
	log.Println(os.Getenv("PATH"))
	log.Println(os.Getenv("LD_LIBRARY_PATH"))
	if err := http.ListenAndServe(":"+defaultPort, router); err != nil {
		log.Println(err)
	}
}
