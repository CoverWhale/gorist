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
	"errors"
	"fmt"
	"net/http"
)

type TableID string

type Tables struct {
	Tables []Table `json:"tables"`
}

type Table struct {
	ID      TableID     `json:"id"`
	Columns []Column    `json:"columns"`
	Fields  TableFields `json:"fields"`
}

type TableFields struct {
	TableRef int  `json:"tableRef"`
	OnDemand bool `json:"onDemand"`
}

// ListTables gets all tables for the specified document ID
func (c *Client) ListTables(document DocumentID) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s/tables", document),
		Method: http.MethodGet,
	}

	return c.httpRequest(request)
}

// CreateTables creates tables in the specified document
func (c *Client) CreateTables(document DocumentID, tables ...Table) (json.RawMessage, error) {
	return c.createTables(document, tables...)
}

func (c *Client) createTables(document DocumentID, tables ...Table) (json.RawMessage, error) {
	for _, v := range tables {
		if len(v.Columns) < 1 {
			return nil, errors.New("columns required to create table")
		}
	}

	return c.writeTable(http.MethodPost, document, tables...)

}

// Patch for tables doesn't follow the spec and doesn't seem to work
//func (c *Client) PatchTables(document Document, tables ...Table) (json.RawMessage, error) {
//	for _, v := range tables {
//		if v.Fields == (TableFields{}) {
//			return nil, errors.New("fields required to patch table")
//		}
//	}
//
//	return c.writeTable(http.MethodPatch, document, tables...)
//}

func (c *Client) writeTable(method string, document DocumentID, tables ...Table) (json.RawMessage, error) {
	t := Tables{
		Tables: tables,
	}
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s/tables", document),
		Method: method,
		Data:   bytes.NewReader(data),
	}

	return c.httpRequest(request)

}
