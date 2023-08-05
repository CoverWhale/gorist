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
	"strconv"
	"time"
)

type WorkspaceID int

type WorkspaceName struct {
	Name string `json:"name,omitempty"`
}

type Workspace struct {
	WorkspaceName
	CreatedAt          time.Time   `json:"createdAt,omitempty"`
	UpdatedAt          string      `json:"updatedAt,omitempty"`
	ID                 WorkspaceID `json:"id,omitempty"`
	IsSupportWorkspace bool        `json:"isSupportWorkspace,omitempty"`
	Org                Org         `json:"org,omitempty"`
	Access             string      `json:"access,omitempty"`
	Owner              Owner       `json:"owner,omitempty"`
	Docs               []Document  `json:"docs,omitempty"`
}

// NewWorkspace creates a new worksapce in the specified organization
func (c *Client) CreateWorkspace(orgID int, name string) (int, error) {
	wn := WorkspaceName{
		Name: name,
	}

	data, err := json.Marshal(wn)
	if err != nil {
		return 0, err
	}

	request := GristRequest{
		Path:   fmt.Sprintf("/api/orgs/%d/workspaces", orgID),
		Method: http.MethodPost,
		Data:   bytes.NewReader(data),
	}

	id, err := c.httpRequest(request)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(id))
}

func (c *Client) GetWorkspace(id WorkspaceID) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/workspaces/%d", id),
		Method: http.MethodGet,
	}

	return c.httpRequest(request)
}
