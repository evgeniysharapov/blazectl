// Copyright 2019 The Samply Development Community
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

package fhir

import (
	"encoding/json"
	fm "github.com/samply/golang-fhir-models/fhir-models/fhir"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// A Client is a FHIR client which combines an HTTP client with the base URL of
// a FHIR server. At minimum, the BaseURL has to be set. HttpClient can be left at
// its default value.
type Client struct {
	httpClient http.Client
	baseURL    url.URL
	auth       ClientAuth
}

// ClientAuth comprises the authentication information used by the Client in
// order to communicate with a FHIR server.
type ClientAuth struct {
	BasicAuthUser     string
	BasicAuthPassword string
}

func NewClient(fhirServerBaseUrl url.URL, auth ClientAuth) *Client {
	// Ensures subsequent calls to ResolveReference do not overwrite the path of the base URL.
	// To avoid this a trailing slash is required.
	if len(fhirServerBaseUrl.Path) > 0 && !strings.HasSuffix(fhirServerBaseUrl.Path, "/") {
		fhirServerBaseUrl.Path = fhirServerBaseUrl.Path + "/"
	}

	return &Client{
		baseURL: fhirServerBaseUrl,
		auth:    auth,
	}
}

// NewCapabilitiesRequest creates a new capabilities interaction request. Uses
// the base URL from the FHIR client and sets JSON Accept header. Otherwise it's
// identical to http.NewRequest.
func (c *Client) NewCapabilitiesRequest() (*http.Request, error) {
	rel := &url.URL{Path: "metadata"}
	req, err := http.NewRequest("GET", c.baseURL.ResolveReference(rel).String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/fhir+json")
	return req, nil
}

// NewTransactionRequest creates a new transaction/batch interaction request.
// Uses the base URL from the FHIR client and sets JSON Accept and Content-Type
// headers. Otherwise it's identical to http.NewRequest.
func (c *Client) NewTransactionRequest(body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", strings.TrimSuffix(c.baseURL.String(), "/"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/fhir+json")
	req.Header.Add("Content-Type", "application/fhir+json")
	return req, nil
}

// NewBatchRequest creates a new transaction/batch interaction request.
// Uses the base URL from the FHIR client and sets JSON Accept and Content-Type
// headers. Otherwise it's identical to http.NewRequest.
func (c *Client) NewBatchRequest(body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", strings.TrimSuffix(c.baseURL.String(), "/"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/fhir+json")
	req.Header.Add("Content-Type", "application/fhir+json")
	return req, nil
}

// NewResourceRequestWithSearch creates a new resource interaction request with an
// additional FHIR search query and sets JSON Accept header. Otherwise it's
// identical to http.NewRequest.
func (c *Client) NewResourceRequestWithSearch(resourceType string, searchQuery url.Values) (*http.Request, error) {
	rel := &url.URL{Path: resourceType, RawQuery: searchQuery.Encode()}
	req, err := http.NewRequest("GET", c.baseURL.ResolveReference(rel).String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/fhir+json")
	return req, nil
}

// NewPaginatedResourceRequest creates a new resource interaction request based on
// a pagination link received from a FHIR server. It sets JSON Accept header and is
// otherwise identical to http.NewRequest.
func (c *Client) NewPaginatedResourceRequest(paginationURL *url.URL) (*http.Request, error) {
	req, err := http.NewRequest("GET", paginationURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/fhir+json")
	return req, nil
}

// Do calls Do on the HTTP client of the FHIR client.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if len(c.auth.BasicAuthUser) != 0 {
		req.SetBasicAuth(c.auth.BasicAuthUser, c.auth.BasicAuthPassword)
	}

	return c.httpClient.Do(req)
}

// CloseIdleConnections calls CloseIdleConnections on the HTTP client of the
// FHIR client.
func (c *Client) CloseIdleConnections() {
	c.httpClient.CloseIdleConnections()
}

// ReadCapabilityStatement reads and unmarshals a capability statement.
func ReadCapabilityStatement(r io.Reader) (fm.CapabilityStatement, error) {
	var capabilityStatement fm.CapabilityStatement
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return capabilityStatement, err
	}
	if err := json.Unmarshal(body, &capabilityStatement); err != nil {
		return capabilityStatement, err
	}
	return capabilityStatement, nil
}

// ReadBundle reads and unmarshals a bundle.
func ReadBundle(r io.Reader) (fm.Bundle, error) {
	var bundle fm.Bundle
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return bundle, err
	}
	return fm.UnmarshalBundle(body)
}
