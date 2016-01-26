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

package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// SlopeOptions defines options for gdaldem slope.
type SlopeOptions struct {
	Percent bool    `json:"percent"`
	Scale   float64 `json:"scale"`
	GeneralOptions
}

// NewSlopeOptions constructs SlopeOptions with default values.
func NewSlopeOptions() *SlopeOptions {
	opts := NewGeneralOptions()
	return &SlopeOptions{
		Percent:        false,
		Scale:          1.0,
		GeneralOptions: *opts,
	}
}

// Slope implements gdaldem slope.
func Slope(w http.ResponseWriter, r *http.Request,
	res *job.OutputMsg, msg job.InputMsg, i, o string) {
	opts := NewSlopeOptions()
	if msg.Options != nil {
		if err := json.Unmarshal(*msg.Options, &opts); err != nil {
			job.BadRequest(w, r, *res, err.Error())
			return
		}
	}

	var args []string
	args = append(args, *msg.Function)
	args = append(args, i)
	args = append(args, o)
	args = append(args, "-of")
	args = append(args, opts.GeneralOptions.Format)
	if opts.GeneralOptions.ComputeEdges {
		args = append(args, "-compute_edges")
	}
	args = append(args, "-alg")
	args = append(args, opts.GeneralOptions.Alg)
	args = append(args, "-b")
	args = append(args, strconv.Itoa(opts.GeneralOptions.Band))
	if opts.GeneralOptions.Quiet {
		args = append(args, "-q")
	}
	if opts.Percent {
		args = append(args, "-p")
	}
	args = append(args, "-s")
	args = append(args, strconv.FormatFloat(opts.Scale, 'f', -1, 64))
	out, err := exec.Command("gdaldem", args...).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
