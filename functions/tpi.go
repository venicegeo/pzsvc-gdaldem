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
	"fmt"
	"net/http"
	"os/exec"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

// // TPIOptions defines options for dart sampling.
// type TPIOptions struct {
// 	// Radius is minimum distance criteria. No two points in the sampled point
// 	// cloud will be closer than the specified radius.
// 	Percent bool    `json:"percent"`
// 	Scale   float64 `json:"scale"`
// }
//
// // NewTPIOptions constructs TPIOptions with default values.
// func NewTPIOptions() *TPIOptions {
// 	return &TPIOptions{
// 		Percent: false,
// 		Scale:   1.0,
// 	}
// }

// TPI implements pdal height.
func TPI(w http.ResponseWriter, r *http.Request,
	res *job.OutputMsg, msg job.InputMsg, i, o string) {
	// opts := NewTPIOptions()
	// if msg.Options != nil {
	// 	if err := json.Unmarshal(*msg.Options, &opts); err != nil {
	// 		job.BadRequest(w, r, *res, err.Error())
	// 		return
	// 	}
	// }

	var args []string
	args = append(args, *msg.Function)
	args = append(args, i)
	args = append(args, o)
	out, err := exec.Command("gdaldem", args...).CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err.Error())
	}
}
