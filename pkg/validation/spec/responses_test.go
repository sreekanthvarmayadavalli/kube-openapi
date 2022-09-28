/*
 Copyright 2022 The Kubernetes Authors.

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

package spec

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/go-openapi/jsonreference"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestResponsesRoundtrip(t *testing.T) {
	type testCase struct {
		Name string
		Responses
	}

	var cases []testCase = []testCase{
		{
			"Decode Ref Object",
			Responses{
				VendorExtensible: VendorExtensible{Extensions: map[string]interface{}{
					"x-<vBŤç,ʡËdSS暙ɑɮ":     "鄲兴ȑʦ衈覻鋕嚮峡jw逓:鮕虫F迢.",
					"x-h":                  "",
					"x-岡ʍ":                 "Đɻ/nǌo鿻曑Œ TĀyĢ",
					"x-绅ƄȆ疩ã[魑銒;苎#砠zPȺ5Aù": "閲ǉǠyư"},
				},
				ResponsesProps: ResponsesProps{
					Default: &Response{
						ResponseProps: ResponseProps{Description: "梱bȿF)渽Ɲō-%x"},
					},
					StatusCodeResponses: map[int]Response{
						200: {
							Refable: Refable{Ref: MustCreateRef("Cat")},
						},
					},
				},
			},
		},
		{
			"Default Full Object",
			Responses{
				VendorExtensible: VendorExtensible{Extensions: map[string]interface{}{
					"x-<vBŤç,ʡËdSS暙ɑɮ":     "鄲兴ȑʦ衈覻鋕嚮峡jw逓:鮕虫F迢.",
					"x-h":                  "",
					"x-岡ʍ":                 "Đɻ/nǌo鿻曑Œ TĀyĢ",
					"x-绅ƄȆ疩ã[魑銒;苎#砠zPȺ5Aù": "閲ǉǠyư"},
				},
				ResponsesProps: ResponsesProps{
					Default: &Response{
						Refable: Refable{Ref: MustCreateRef("Dog")},
					},
					StatusCodeResponses: map[int]Response{
						200: {
							ResponseProps: ResponseProps{
								Description: "梱bȿF)渽Ɲō-%x",
								Headers: map[string]Header{
									"a header": header,
								},
								Schema: &Schema{
									VendorExtensible: VendorExtensible{Extensions: map[string]interface{}{"x-framework": "go-swagger"}},
									SchemaProps: SchemaProps{
										Ref:              MustCreateRef("Cat"),
										Type:             []string{"string"},
										Format:           "date",
										Description:      "the description of this schema",
										Title:            "the title",
										Default:          "blah",
										Maximum:          float64Ptr(100),
										ExclusiveMaximum: true,
										ExclusiveMinimum: true,
										Minimum:          float64Ptr(5),
										MaxLength:        int64Ptr(100),
										MinLength:        int64Ptr(5),
										Pattern:          "\\w{1,5}\\w+",
										MaxItems:         int64Ptr(100),
										MinItems:         int64Ptr(5),
										UniqueItems:      true,
										MultipleOf:       float64Ptr(5),
										Enum:             []interface{}{"hello", "world"},
										MaxProperties:    int64Ptr(5),
										MinProperties:    int64Ptr(1),
										Required:         []string{"id", "name"},
										Items:            &SchemaOrArray{Schema: &Schema{SchemaProps: SchemaProps{Type: []string{"string"}}}},
										AllOf:            []Schema{{SchemaProps: SchemaProps{Type: []string{"string"}}}},
										Properties: map[string]Schema{
											"id":   {SchemaProps: SchemaProps{Type: []string{"integer"}, Format: "int64"}},
											"name": {SchemaProps: SchemaProps{Type: []string{"string"}}},
										},
										AdditionalProperties: &SchemaOrBool{Allows: true, Schema: &Schema{SchemaProps: SchemaProps{
											Type:   []string{"integer"},
											Format: "int32",
										}}},
									},
									SwaggerSchemaProps: SwaggerSchemaProps{
										Discriminator: "not this",
										ReadOnly:      true,
										ExternalDocs: &ExternalDocumentation{
											Description: "the documentation etc",
											URL:         "http://readthedocs.org/swagger",
										},
										Example: []interface{}{
											map[string]interface{}{
												"id":   float64(1),
												"name": "a book",
											},
											map[string]interface{}{
												"id":   float64(2),
												"name": "the thing",
											},
										},
									},
								},
								Examples: map[string]interface{}{
									"example1": "example text",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {

			jsonText, err := json.Marshal(&test.Responses)
			require.NoError(t, err)

			var decoded Responses
			err = json.Unmarshal(jsonText, &decoded)
			require.NoError(t, err)

			if !reflect.DeepEqual(&test.Responses, &decoded) {
				t.Fatal(cmp.Diff(&test.Responses, &decoded, cmpopts.IgnoreUnexported(jsonreference.Ref{})))
			}
		})
	}
}

func TestResponsesDecodeErrors(t *testing.T) {
	type testCase struct {
		Name     string
		JSONText string
	}

	var cases []testCase = []testCase{
		{
			"FailDecodeDefault",
			`{"x-extension":"an extension","default":"wrong type object"}`,
		},
	}

	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			deserialized := Responses{}
			err := json.Unmarshal([]byte(v.JSONText), &deserialized)
			require.Error(t, err, "expefcted an error")
		})
	}
}
