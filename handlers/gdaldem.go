/*
Copyright 2016, RadiantBlue Technologies, Inc.

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

package handlers

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-gdaldem/functions"
	"github.com/venicegeo/pzsvc-sdk-go/job"
	"github.com/venicegeo/pzsvc-sdk-go/utils"
)

// GdalDemHandler handles PDAL jobs.
func GdalDemHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Create the job output message. No matter what happens, we should always be
	// able to populate the StartedAt field.
	var res job.OutputMsg
	res.StartedAt = time.Now()

	msg := job.GetInputMsg(w, r, res)

	// Throw 400 if the JobInput does not specify a function.
	if msg.Function == nil {
		http.Error(w, "must provide a function", http.StatusBadRequest)
		return
	}

	// If everything is okay up to this point, we will echo the JobInput in the
	// JobOutput and mark the job as Running.
	res.Input = msg
	job.Update(job.Running, r)

	// Make/execute the requested function.
	switch *msg.Function {
	case "aspect":
		utils.MakeFunction(functions.Aspect)(w, r, &res, msg)

	case "hillshade":
		utils.MakeFunction(functions.HillshadeFunction)(w, r, &res, msg)

	case "slope":
		utils.MakeFunction(functions.Slope)(w, r, &res, msg)

	// An unrecognized function will result in 400 error, with message explaining
	// how to list available functions.
	default:
		job.BadRequest(w, r, res, "")
		return
	}

	// If we made it here, we can record the FinishedAt time, notify the job
	// manager of success, and return 200.
	res.FinishedAt = time.Now()
	job.Okay(w, r, res, "Success!")
}
