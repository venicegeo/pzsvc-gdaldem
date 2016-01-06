package handlers

import (
	"net/http"
	"time"

	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/venicegeo/pzsvc-sdk-go/objects"
	"github.com/venicegeo/pzsvc-gdaldem/Godeps/_workspace/src/github.com/venicegeo/pzsvc-sdk-go/utils"
	"github.com/venicegeo/pzsvc-gdaldem/functions"
)

// GdalDemHandler handles PDAL jobs.
func GdalDemHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Create the job output message. No matter what happens, we should always be
	// able to populate the StartedAt field.
	var res objects.JobOutput
	res.StartedAt = time.Now()

	msg := utils.GetJobInput(w, r, res)

	// Throw 400 if the JobInput does not specify a function.
	if msg.Function == nil {
		http.Error(w, "must provide a function", http.StatusBadRequest)
		return
	}

	// If everything is okay up to this point, we will echo the JobInput in the
	// JobOutput and mark the job as Running.
	res.Input = msg
	utils.UpdateJobManager(objects.Running, r)

	// Make/execute the requested function.
	switch *msg.Function {
	case "hillshade":
		utils.MakeFunction(functions.HillshadeFunction)(w, r, &res, msg)

	// An unrecognized function will result in 400 error, with message explaining
	// how to list available functions.
	default:
		utils.BadRequest(w, r, res, "")
		return
	}

	// If we made it here, we can record the FinishedAt time, notify the job
	// manager of success, and return 200.
	res.FinishedAt = time.Now()
	utils.Okay(w, r, res, "Success!")
}
