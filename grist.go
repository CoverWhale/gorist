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

type ClientOpt func(*Client)

type Client struct {
	// Grist API token
	Token string
	// Grist Server URL
	URL string
	// HTTP Client
	Client *http.Client
	// Global Filter applied to all requests. Can be overridden with request filter
	GlobalFilter json.RawMessage
}

type GristRequest struct {
	Path     string
	Method   string
	Document DocumentID
	Table    TableID
	// Overrides a client global filter
	Filter json.RawMessage
	Data   io.Reader
}

type GristRequestOpt func(*GristRequest)

func SetDocument(d DocumentID) GristRequestOpt {
	return func(r *GristRequest) {
		r.Document = d
	}
}

func SetTable(t TableID) GristRequestOpt {
	return func(r *GristRequest) {
		r.Table = t
	}
}

func SetFilter(f json.RawMessage) GristRequestOpt {
	return func(r *GristRequest) {
		r.Filter = f
	}
}

func NewClient(opts ...ClientOpt) *Client {
	c := &Client{}

	for _, v := range opts {
		v(c)
	}

	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	return c
}

// SetAPIKey sets the API key in the client
func SetAPIKey(key string) ClientOpt {
	return func(c *Client) {
		c.Token = key
	}
}

// SetURL sets the URL in the client
func SetURL(url string) ClientOpt {
	return func(c *Client) {
		c.URL = url
	}
}

// SetHTTPClient overrides the default HTTP client
func SetHTTPClient(h *http.Client) ClientOpt {
	return func(c *Client) {
		c.Client = h
	}
}

// SetClientGlobalFilter adds a filter to all client requests. Can be overridden with a request filter
func SetClientGlobalFilter(f json.RawMessage) ClientOpt {
	return func(c *Client) {
		c.GlobalFilter = f
	}
}

func (c *Client) httpRequest(request GristRequest) (json.RawMessage, error) {
	url := fmt.Sprintf("%s%s", c.URL, request.Path)
	token := fmt.Sprintf("Bearer %s", c.Token)

	req, err := http.NewRequest(request.Method, url, request.Data)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	q := req.URL.Query()

	if c.GlobalFilter != nil {
		q.Add("filter", string(c.GlobalFilter))
	}

	if request.Filter != nil {
		q.Del("filter")
		q.Add("filter", string(request.Filter))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, GristError{
			StatusCode: resp.StatusCode,
			Details:    string(body),
		}
	}

	return body, nil
}

type GristError struct {
	StatusCode int
	Details    string
}

func (g GristError) Error() string {
	return fmt.Sprintf("status: %d, details: %s", g.StatusCode, g.Details)
}
