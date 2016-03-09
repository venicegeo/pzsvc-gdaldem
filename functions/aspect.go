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

// AspectOptions defines options for gdaldem aspect.
type AspectOptions struct {
	Trigonometric bool `json:"trigonometric"`
	ZeroForFlat   bool `json:"zero_for_flat"`
	GeneralOptions
	AlgorithmOptions
}

// NewAspectOptions constructs AspectOptions with default values.
func NewAspectOptions() *AspectOptions {
	opts := NewGeneralOptions()
	alg := NewAlgorithmOptions()
	return &AspectOptions{
		Trigonometric:    false,
		ZeroForFlat:      false,
		GeneralOptions:   *opts,
		AlgorithmOptions: *alg,
	}
}

// Aspect implements gdaldem aspect.
func Aspect(i, o string, options *json.RawMessage) ([]byte, error) {
	opts := NewAspectOptions()
	if options != nil {
		if err := json.Unmarshal(*options, &opts); err != nil {
			return nil, err
		}
	}

	var args []string
	args = append(args, "aspect")
	args = append(args, i)
	args = append(args, o)
	args = append(args, "-of")
	args = append(args, opts.GeneralOptions.Format)
	if opts.GeneralOptions.ComputeEdges {
		args = append(args, "-compute_edges")
	}
	args = append(args, "-alg")
	args = append(args, opts.AlgorithmOptions.Alg)
	args = append(args, "-b")
	args = append(args, strconv.Itoa(opts.GeneralOptions.Band))
	if opts.GeneralOptions.Quiet {
		args = append(args, "-q")
	}
	if opts.Trigonometric {
		args = append(args, "-trigonometric")
	}
	if opts.ZeroForFlat {
		args = append(args, "-zero_for_flat")
	}
	out, err := exec.Command("gdaldem", args...).CombinedOutput()

	fmt.Println(string(out))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
