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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) UpdateRecord(document DocumentID, table TableID, r io.Reader) (json.RawMessage, error) {
	path := fmt.Sprintf("/api/docs/%s/tables/%s/records", document, table)
	request := GristRequest{
		Path:   path,
		Method: http.MethodPatch,
		Data:   r,
	}
	return c.httpRequest(request)
}


func (c *Client) GetRecordsWithOptions(opts ...GristRequestOpt) (json.RawMessage, error) {
	var r GristRequest

	for _, opt := range opts {
		opt(&r)
	}

	r.Path = fmt.Sprintf("/api/docs/%s/tables/%s/records", r.Document, r.Table)

	return c.getRecords(r)
}

func (c *Client) GetRecords(document DocumentID, table TableID) (json.RawMessage, error) {

	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s/tables/%s/records", document, table),
		Method: http.MethodGet,
	}

	return c.getRecords(request)
}

func (c *Client) GetFilteredRecords(document DocumentID, table TableID, filter json.RawMessage) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s/tables/%s/records", document, table),
		Method: http.MethodGet,
	}

	if filter != nil {
		request.Filter = filter
	}

	return c.getRecords(request)

}

func (c *Client) getRecords(r GristRequest) (json.RawMessage, error) {
	return c.httpRequest(r)
}

func (c *Client) CreateRecord(document DocumentID, table TableID, r io.Reader) (json.RawMessage, error) {
	path := fmt.Sprintf("/api/docs/%s/tables/%s/records", document, table)
	request := GristRequest{
		Path:   path,
		Method: http.MethodPost,
		Data:   r,
	}
	return c.httpRequest(request)
}
