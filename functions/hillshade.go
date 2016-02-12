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
	"os/exec"
	"strconv"
)

// HillshadeOptions defines options for gdaldem hillshade.
type HillshadeOptions struct {
	ZFactor  float64 `json:"zFactor"`
	Scale    float64 `json:"scale"`
	Azimuth  float64 `json:"azimuth"`
	Altitude float64 `json:"altitude"`
	Combined bool    `json:"combined"`
	GeneralOptions
}

// NewHillshadeOptions constructs HillshadeOptions with default values.
func NewHillshadeOptions() *HillshadeOptions {
	opts := NewGeneralOptions()
	return &HillshadeOptions{
		ZFactor:        1.0,
		Scale:          1.0,
		Azimuth:        315.0,
		Altitude:       45.0,
		Combined:       false,
		GeneralOptions: *opts,
	}
}

// HillshadeFunction implements gdaldem hillshade.
func HillshadeFunction(i, o string, options *json.RawMessage) ([]byte, error) {
	opts := NewHillshadeOptions()
	if options != nil {
		if err := json.Unmarshal(*options, &opts); err != nil {
			return nil, err
		}
	}

	var args []string
	args = append(args, "hillshade")
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
	args = append(args, "-z")
	args = append(args, strconv.FormatFloat(opts.ZFactor, 'f', -1, 64))
	args = append(args, "-s")
	args = append(args, strconv.FormatFloat(opts.Scale, 'f', -1, 64))
	args = append(args, "-az")
	args = append(args, strconv.FormatFloat(opts.Azimuth, 'f', -1, 64))
	args = append(args, "-alt")
	args = append(args, strconv.FormatFloat(opts.Altitude, 'f', -1, 64))
	if opts.Combined {
		args = append(args, "-combined")
	}
	out, err := exec.Command("gdaldem", args...).CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
