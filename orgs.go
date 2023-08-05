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
	"net/http"
)

type Org struct {
	Name               string         `json:"name,omitempty"`
	CreatedAt          string         `json:"createdAt,omitempty"`
	UpdatedAt          string         `json:"updatedAt,omitempty"`
	ID                 int            `json:"id,omitempty"`
	Domain             string         `json:"domain,omitempty"`
	Host               interface{}    `json:"host,omitempty"`
	Owner              Owner          `json:"owner,omitempty"`
	BillingAccount     BillingAccount `json:"billingAccount,omitempty"`
	IsSupportWorkspace bool           `json:"isSupportWorkspace,omitempty"`
	Access             string         `json:"access,omitempty"`
	Public             bool           `json:"public,omitempty"`
	Docs               []Document     `json:"docs,omitempty"`
	OrgDomain          string         `json:"orgDomain,omitempty"`
}

type Owner struct {
	ID      int     `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Picture string  `json:"picture,omitempty"`
	Options Options `json:"options,omitempty"`
	Ref     string  `json:"ref,omitempty"`
}

type Options struct {
	Locale string `json:"locale,omitempty"`
}

// ListOrgs gets all organizations associated with the current token
func (c *Client) ListOrgs() (json.RawMessage, error) {
	request := GristRequest{
		Path:   "/api/orgs",
		Method: http.MethodGet,
	}

	return c.httpRequest(request)
}

// GetOrg gets the details of an organization based on the ID
func (c *Client) GetOrg(id string) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/orgs/%s", id),
		Method: http.MethodGet,
	}

	return c.httpRequest(request)
}

// GetOrgWorkspacesAndDocuments retrieves all workspaces and documents in the workspaces for an organization
func (c *Client) GetOrgWorkspacesAndDocuments(id string) (json.RawMessage, error) {
	request := GristRequest{
		Path:   fmt.Sprintf("/api/orgs/%s/workspaces", id),
		Method: http.MethodGet,
	}

	return c.httpRequest(request)
}
