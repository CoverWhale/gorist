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
	"time"
)

type DocumentID string

type NewDocRequest struct {
	Name     string `json:"name"`
	IsPinned bool   `json:"isPinned"`
}

type Document struct {
	Name      string        `json:"name"`
	CreatedAt time.Time     `json:"createdAt,omitempty"`
	UpdatedAt time.Time     `json:"updatedAt,omitempty"`
	ID        DocumentID    `json:"id,omitempty"`
	IsPinned  bool          `json:"isPinned"`
	URLID     string        `json:"urlId,omitempty"`
	TrunkID   interface{}   `json:"trunkId,omitempty"`
	Type      interface{}   `json:"type,omitempty"`
	Workspace Workspace     `json:"workspace,omitempty"`
	Aliases   []Aliases     `json:"aliases,omitempty"`
	Forks     []interface{} `json:"forks,omitempty"`
	Access    string        `json:"access,omitempty"`
}

type SnapshotWindow struct {
	Count int    `json:"count,omitempty"`
	Unit  string `json:"unit,omitempty"`
}

type Features struct {
	Workspaces                         bool           `json:"workspaces,omitempty"`
	MaxSharesPerWorkspace              int            `json:"maxSharesPerWorkspace,omitempty"`
	MaxSharesPerDoc                    int            `json:"maxSharesPerDoc,omitempty"`
	SnapshotWindow                     SnapshotWindow `json:"snapshotWindow,omitempty"`
	BaseMaxRowsPerDocument             int            `json:"baseMaxRowsPerDocument,omitempty"`
	BaseMaxAPIUnitsPerDocumentPerDay   int            `json:"baseMaxApiUnitsPerDocumentPerDay,omitempty"`
	BaseMaxDataSizePerDocument         int            `json:"baseMaxDataSizePerDocument,omitempty"`
	BaseMaxAttachmentsBytesPerDocument int            `json:"baseMaxAttachmentsBytesPerDocument,omitempty"`
	GracePeriodDays                    int            `json:"gracePeriodDays,omitempty"`
	BaseMaxAssistantCalls              int            `json:"baseMaxAssistantCalls,omitempty"`
}

type Product struct {
	ID       int      `json:"id,omitempty"`
	Name     string   `json:"name,omitempty"`
	Features Features `json:"features,omitempty"`
}

type BillingAccount struct {
	ID              int         `json:"id,omitempty"`
	Individual      bool        `json:"individual,omitempty"`
	InGoodStanding  bool        `json:"inGoodStanding,omitempty"`
	Status          interface{} `json:"status,omitempty"`
	ExternalID      interface{} `json:"externalId,omitempty"`
	ExternalOptions interface{} `json:"externalOptions,omitempty"`
	Product         Product     `json:"product,omitempty"`
}

type Aliases struct {
	OrgID     int    `json:"orgId,omitempty"`
	URLID     string `json:"urlId,omitempty"`
	DocID     string `json:"docId,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

func (c *Client) CreateDocument(workspace WorkspaceID, name string, isPinned bool) (string, error) {
	id, err := c.createDocument(workspace, name, isPinned)
	if err != nil {
		return "", err
	}

	return string(id), nil

}

func (c *Client) createDocument(workspace WorkspaceID, name string, isPinned bool) (json.RawMessage, error) {
	doc := NewDocRequest{
		Name:     name,
		IsPinned: isPinned,
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	request := GristRequest{
		Path:   fmt.Sprintf("/api/workspaces/%d/docs", workspace),
		Method: http.MethodPost,
		Data:   bytes.NewReader(data),
	}

	resp, err := c.httpRequest(request)
	if err != nil {
		return nil, err
	}

	// grist returns quotes around the ID so strip those off
	return bytes.ReplaceAll(resp, []byte(`"`), nil), nil
}

func (c *Client) GetDocument(id string) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s", id),
		Method: http.MethodGet,
	}
	return c.httpRequest(request)
}

func (c *Client) DeleteDocument(id string) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/docs/%s", id),
		Method: http.MethodDelete,
	}

	return c.httpRequest(request)
}
