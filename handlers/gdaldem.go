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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/venicegeo/pzsvc-gdaldem/functions"
	"github.com/venicegeo/pzsvc-sdk-go/job"
	"github.com/venicegeo/pzsvc-sdk-go/s3"
)

type AppError struct {
	Error   error
	Message string
	Code    int
}

// InputMsg defines the expected input JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type InputMsg struct {
	Source      s3.Bucket        `json:"source,omitempty"`
	Function    *string          `json:"function,omitempty"`
	Options     *json.RawMessage `json:"options,omitempty"`
	Destination s3.Bucket        `json:"destination,omitempty"`
}

// FunctionFunc defines the signature of our function creator.
type FunctionFunc func(InputMsg) ([]byte, error)

// MakeFunction wraps the individual PDAL functions.
// Parse the input and output filenames, creating files as needed. Download the
// input data and upload the output data.
func MakeFunction(fn func(string, string, *json.RawMessage) ([]byte, error)) FunctionFunc {
	return func(msg InputMsg) ([]byte, error) {
		var inputName, outputName string
		var fileIn, fileOut *os.File

		// Split the source S3 key string, interpreting the last element as the
		// input filename. Create the input file, throwing 500 on error.
		inputName = s3.ParseFilenameFromKey(msg.Source.Key)
		fileIn, err := os.Create(inputName)
		if err != nil {
			return nil, err
		}
		defer fileIn.Close()

		// If provided, split the destination S3 key string, interpreting the last
		// element as the output filename. Create the output file, throwing 500 on
		// error.
		if len(msg.Destination.Key) > 0 {
			outputName = s3.ParseFilenameFromKey(msg.Destination.Key)
		}

		// Download the source data from S3, throwing 500 on error.
		err = s3.Download(fileIn, msg.Source.Bucket, msg.Source.Key)
		if err != nil {
			return nil, err
		}

		os.Remove(outputName)

		// Run the PDAL function.
		retval, err := fn(inputName, outputName, msg.Options)
		if err != nil {
			return nil, err
		}

		// If an output has been created, upload the destination data to S3,
		// throwing 500 on error.
		if len(msg.Destination.Key) > 0 {
			fileOut, err = os.Open(outputName)
			if err != nil {
				return nil, err
			}
			defer fileOut.Close()
			err = s3.Upload(fileOut, msg.Destination.Bucket, msg.Destination.Key)
			if err != nil {
				return nil, err
			}
		}

		return retval, nil
	}
}

// GdalDemHandler handles PDAL jobs.
func GdalDemHandler(w http.ResponseWriter, r *http.Request) *AppError {
	// Create the job output message. No matter what happens, we should always be
	// able to populate the StartedAt field.
	var res job.OutputMsg
	res.StartedAt = time.Now()

	var msg InputMsg

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		return &AppError{nil, "No JSON", http.StatusBadRequest}
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &AppError{err, err.Error(), http.StatusInternalServerError}
	}

	// Throw 400 if we cannot unmarshal the body as a valid InputMsg.
	if err := json.Unmarshal(b, &msg); err != nil {
		return &AppError{err, err.Error(), http.StatusBadRequest}
	}

	// Throw 400 if the JobInput does not specify a function.
	if msg.Function == nil {
		return &AppError{nil, "Must provide a function", http.StatusBadRequest}
	}

	// Make/execute the requested function.
	switch *msg.Function {
	case "aspect":
		_, err := MakeFunction(functions.Aspect)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "hillshade":
		_, err := MakeFunction(functions.HillshadeFunction)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "roughness":
		_, err := MakeFunction(functions.Roughness)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "slope":
		_, err := MakeFunction(functions.Slope)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "TPI":
		_, err := MakeFunction(functions.TPI)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "TRI":
		_, err := MakeFunction(functions.TRI)(msg)
		if err != nil {
			return &AppError{err, err.Error(), http.StatusInternalServerError}
		}

		res.FinishedAt = time.Now()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		res.Code = http.StatusOK
		res.Message = "Success"

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	// An unrecognized function will result in 400 error, with message explaining
	// how to list available functions.
	default:
		return &AppError{nil, "Unrecognized function", http.StatusBadRequest}
	}

	return nil
}
