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

// GeneralOptions defines options for multiple gdaldem functions.
type GeneralOptions struct {
	Format       string `json:"format"`
	ComputeEdges bool   `json:"compute_edges"`
	Band         int    `json:"band"`
	Quiet        bool   `json:"quiet"`
}

// NewGeneralOptions constructs GeneralOptions with default values.
func NewGeneralOptions() *GeneralOptions {
	return &GeneralOptions{
		Format:       "GTiff",
		ComputeEdges: false,
		Band:         1,
		Quiet:        false,
	}
}

// AlgorithmOptions defines options for algorithms.
type AlgorithmOptions struct {
	Alg string `json:"alg"`
}

// NewAlgorithmOptions construct AlgorithmOptions with default values.
func NewAlgorithmOptions() *AlgorithmOptions {
	return &AlgorithmOptions{
		Alg: "Horn",
	}
}
