// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	jsontesting "k8s.io/kube-openapi/pkg/util/jsontesting"
)

var operation = Operation{
	VendorExtensible: VendorExtensible{
		Extensions: map[string]interface{}{
			"x-framework": "go-swagger",
		},
	},
	OperationProps: OperationProps{
		Description: "operation description",
		Consumes:    []string{"application/json", "application/x-yaml"},
		Produces:    []string{"application/json", "application/x-yaml"},
		Schemes:     []string{"http", "https"},
		Tags:        []string{"dogs"},
		Summary:     "the summary of the operation",
		ID:          "sendCat",
		Deprecated:  true,
		Security: []map[string][]string{
			{
				"apiKey": {},
			},
		},
		Parameters: []Parameter{
			{Refable: Refable{Ref: MustCreateRef("Cat")}},
		},
		Responses: &Responses{
			ResponsesProps: ResponsesProps{
				Default: &Response{
					ResponseProps: ResponseProps{
						Description: "void response",
					},
				},
			},
		},
	},
}

const operationJSON = `{
	"description": "operation description",
	"x-framework": "go-swagger",
	"consumes": [ "application/json", "application/x-yaml" ],
	"produces": [ "application/json", "application/x-yaml" ],
	"schemes": ["http", "https"],
	"tags": ["dogs"],
	"summary": "the summary of the operation",
	"operationId": "sendCat",
	"deprecated": true,
	"security": [ { "apiKey": [] } ],
	"parameters": [{"$ref":"Cat"}],
	"responses": {
		"default": {
			"description": "void response"
		}
	}
}`

func TestIntegrationOperation(t *testing.T) {
	var actual Operation
	if assert.NoError(t, json.Unmarshal([]byte(operationJSON), &actual)) {
		assert.EqualValues(t, actual, operation)
	}

	assertParsesJSON(t, operationJSON, operation)
}

func TestOperationRoundtrip(t *testing.T) {
	cases := []jsontesting.RoundTripTestCase{
		{
			// Show at least one field from each embededd struct sitll allows
			// roundtrips successfully
			Name: "UnmarshalEmbedded",
			JSON: `{
				"description": "a cool description",
				"x-framework": "swagger-go"
			  }`,
			Object: &Operation{
				VendorExtensible{Extensions{
					"x-framework": "swagger-go",
				}},
				OperationProps{
					Description: "a cool description",
				},
			},
		}, {
			Name:   "BasicCase",
			JSON:   operationJSON,
			Object: &operation,
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			require.NoError(t, tcase.RoundTripTest(&Operation{}))
		})
	}
}
