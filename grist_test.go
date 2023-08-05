// Copyright 2023 Cover Whale Insurance Solutions Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type GristTest struct {
	expected Data
	data     Data
	document DocumentID
	key      string
	filter   json.RawMessage
}

type Data struct {
	Records []Record `json:"records"`
}

type Record struct {
	ID     int         `json:"id"`
	Fields []TestField `json:"fields"`
}

type TestField struct {
	Foo string `json:"foo"`
}

func getHandler(gt GristTest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if gt.document == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if gt.key == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		u := r.URL.Query().Get("filter")

		if gt.filter != nil && u != string(gt.filter) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// always filter only the first item. Lazy test.
		if gt.filter != nil {
			gt.data.Records = gt.data.Records[:1]
		}

		if err := json.NewEncoder(w).Encode(gt.data); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func postHandler(gt GristTest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data Data

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if !reflect.DeepEqual(gt.data, data) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(gt.expected); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func TestGetDocument(t *testing.T) {
	tt := []struct {
		name     string
		key      string
		document DocumentID
		expected Data
		data     Data
		err      error
		filter   json.RawMessage
	}{
		{
			name: "document request",
			key:  "test", document: DocumentID("document1"),
			data: Data{
				Records: []Record{
					{
						ID: 1,
						Fields: []TestField{
							{
								Foo: "test",
							},
						},
					},
				},
			},
			expected: Data{
				Records: []Record{
					{
						ID: 1,
						Fields: []TestField{
							{
								Foo: "test",
							},
						},
					},
				},
			},
		},
		{
			name: "filtered document request",
			key:  "test", document: "document1",
			data: Data{
				Records: []Record{
					{
						ID: 1,
						Fields: []TestField{
							{
								Foo: "test",
							},
						},
					},
					{
						ID: 1,
						Fields: []TestField{
							{
								Foo: "test",
							},
						},
					},
				},
			},
			expected: Data{
				Records: []Record{
					{
						ID: 1,
						Fields: []TestField{
							{
								Foo: "test",
							},
						},
					},
				},
			},
			filter: json.RawMessage(`{"id": [1]}`),
		},
		{name: "no api key", document: "document1", err: fmt.Errorf("%s", http.StatusText(http.StatusUnauthorized))},
		{name: "no document", document: "", err: fmt.Errorf("%s", http.StatusText(http.StatusNotFound))},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			gt := GristTest{
				expected: v.expected,
				document: v.document,
				key:      v.key,
				data:     v.data,
				filter:   v.filter,
			}

			s := httptest.NewServer(getHandler(gt))
			defer s.Close()

			c := NewClient(
				SetURL(s.URL),
				SetAPIKey(v.key),
			)

			res, err := c.GetRecordsWithOptions(
				SetDocument(v.document),
				SetFilter(v.filter),
			)
			if v.err != nil && v.err == nil {
				t.Errorf("expected no errors but got %v", err)
			}

			// expected error here so we are good
			if v.err != nil && err != nil {
				return
			}

			var d Data
			if err := json.Unmarshal(res, &d); err != nil {
				t.Errorf("error unmarshaling: %v", err)
			}

			if !reflect.DeepEqual(d, v.expected) {
				t.Errorf("expected \n%#v\nbut got \n%#v", v.expected, d)
			}

		})
	}
}

func TestCreateRecord(t *testing.T) {
	tt := []struct {
		name     string
		data     Data
		expected Data
		err      error
	}{
		{
			name: "normal post",
			data: Data{
				Records: []Record{
					{
						Fields: []TestField{
							{
								Foo: "test",
							},
						},
					},
				},
			},
			expected: Data{
				Records: []Record{
					{
						ID: 1,
					},
				},
			},
		},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			gt := GristTest{
				data:     v.data,
				expected: v.expected,
			}

			s := httptest.NewServer(postHandler(gt))
			defer s.Close()

			c := NewClient(
				SetURL(s.URL),
			)

			d, err := json.Marshal(v.data)
			if err != nil {
				t.Errorf("error marshaling: %v", err)
			}

			res, err := c.CreateRecord("test", "test", bytes.NewReader(d))
			if err != nil && v.err == nil {
				t.Errorf("expected no errors but got %v", err)
			}

			// expected error here so we are good
			if v.err != nil && err != nil {
				return
			}

			var respData Data
			if err := json.Unmarshal(res, &respData); err != nil {
				t.Errorf("error unmarshaling: %v", err)
			}

			if !reflect.DeepEqual(respData, v.expected) {
				t.Errorf("expected \n%#v\nbut got \n%#v", v.expected, respData)
			}

		})
	}
}
