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

type FieldType string

const (
	IntField     FieldType = "Int"
	AnyField     FieldType = "Any"
	TextField    FieldType = "Text"
	NumericField FieldType = "Numeric"
	BoolField    FieldType = "Bool"
	DateField    FieldType = "Date"
	ChoiceField  FieldType = "Choice"
)

type Recalc int

const (
	NewOnly Recalc = iota
	Never
	NewAndUpdate
)

func NewDateTimeField(timezone string) FieldType {
	return FieldType(fmt.Sprintf("DateTime:%s", timezone))
}

func NewRefField(tableID string) FieldType {
	return FieldType(fmt.Sprintf("Ref:%s", tableID))
}

type Columns struct {
	Columns []Column `json:"columns"`
}

type Column struct {
	ID     string      `json:"id"`
	Fields ColumnField `json:"fields"`
}

type ColumnField struct {
	Label               string    `json:"label"`
	Type                FieldType `json:"type"`
	Formula             string    `json:"formula,omitempty"`
	IsFormula           bool      `json:"isFormula,omitempty"`
	WidgetOptions       string    `json:"widgetOptions,omitempty"`
	UntieColIDFromLabel bool      `json:"untieColIdFromLabel,omitempty"`
	RecalcWhen          Recalc    `json:"recalcWhen"`
	VisibleCol          int       `json:"visibleCol,omitempty"`
	RecalcDeps          []int     `json:"recalcDeps,omitempty"`
}

func (f *ColumnField) MarshalJSON() ([]byte, error) {
	type FieldAlias ColumnField
	t := &struct {
		*FieldAlias
	}{
		FieldAlias: (*FieldAlias)(f),
	}

	if f.IsFormula && f.Formula == "" {
		return nil, errors.New("formula must be set for formula field")
	}

	return json.Marshal(t)

}

func (c *Client) GetColumns(document DocumentID, table Table) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s/tables/%s/columns", document, table.ID),
		Method: http.MethodGet,
	}

	return c.httpRequest(request)
}

func (c *Client) CreateColumns(document DocumentID, table Table, columns ...Column) (json.RawMessage, error) {
	return c.writeColumns(http.MethodPost, document, table, columns...)
}

func (c *Client) PatchColumns(document DocumentID, table Table, columns ...Column) (json.RawMessage, error) {
	return c.writeColumns(http.MethodPatch, document, table, columns...)
}

func (c *Client) writeColumns(method string, document DocumentID, table Table, columns ...Column) (json.RawMessage, error) {
	cols := Columns{
		Columns: columns,
	}

	data, err := json.Marshal(cols)
	if err != nil {
		return nil, err
	}

	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s/tables/%s/columns", document, table.ID),
		Method: method,
		Data:   bytes.NewReader(data),
	}

	return c.httpRequest(request)

}
