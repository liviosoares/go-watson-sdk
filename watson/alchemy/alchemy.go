//
// Copyright (C) IBM Corporation 2016, Livio Soares <lsoares@us.ibm.com>
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

// Package alchemy provides an interface to Alchemy based API endpoints.
// This package contains general methods used by the Alchemy specific family of endpoints,
// including alchemy_language, alchemy_data_news and alchemy_vision.
package alchemy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	watsonClient *watson.Client
}

// NewClient uses the cfg configuration to connect to the Alchemy endpoints.
func NewClient(cfg watson.Config) (Client, error) {
	cfg.Credentials.ServiceName = "alchemy_api"
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	alchemy := Client{watsonClient: client}
	return alchemy, nil
}

// DetectAlchemyPath takes in a byte slice, typically encoding a string, and determine which type
// of API to use amongst the 3 variants in AlchemyAPI: URL, HTML or text
func detectAlchemyPath(data []byte) (key string, pathPrefix string, err error) {
	if len(data) == 0 {
		return "", "", errors.New("could not detect data type: empty data")
	}

	if strings.HasPrefix(string(data), "http://") || strings.HasPrefix(string(data), "https://") {
		return "url", "/url/URL", nil
	}

	typ := http.DetectContentType(data)
	if strings.HasPrefix(typ, "text/html") {
		return "html", "/html/HTML", nil
	}

	return "text", "/text/Text", nil
}

// BaseResponse contains common fields returned by a few different Alchemy endpoints
type BaseResponse struct {
	Status           string `json:"status"`
	Usage            string `json:"usage,omitempty"`
	Url              string `json:"url,omitempty"`
	TotalTransaction int    `json:"totalTransactions,string,omitempty"`
	Language         string `json:"language,omitempty"`
	Text             string `json:"text,omitempty"`
	StatusInfo       string `json:"statusInfo,omitempty"`
}

// Call uses POST method to call an Alchemy API endpoint. The method authenticates the call using the ApiKey in the Client object.
// payload is the content passed in the url, html or text keys.
// options can be used to pass additional query parameters to the call.
// out is the object used for unmarshalling the returned JSON
func (c Client) Call(pathSuffix string, payload []byte, options map[string]interface{}, out interface{}) error {
	dataKey, pathPrefix, err := detectAlchemyPath(payload)
	if err != nil {
		return err
	}

	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("apikey", c.watsonClient.Creds.ApiKey)
	q.Set("outputMode", "json")
	q.Set(dataKey, string(payload))

	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	body, err := c.watsonClient.MakeRequest("POST", pathPrefix+pathSuffix, strings.NewReader(q.Encode()), headers)
	// fmt.Println(string(body))
	if err != nil {
		return err
	}

	var baseResponse BaseResponse
	err = json.Unmarshal(body, &baseResponse)
	if err == nil && strings.EqualFold(baseResponse.Status, "error") {
		return errors.New(baseResponse.StatusInfo)
	}

	return json.Unmarshal(body, out)
}

// Get uses the GET method to call an Alchemy API endpoint. The method authenticates the call using the ApiKey in the Client object.
// query can be used to pass additional query parameters to the call.
// out is the object used for unmarshalling the returned JSON
func (c Client) Get(path string, query map[string]interface{}, out interface{}) error {
	q := url.Values{}
	for k, v := range query {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("apikey", c.watsonClient.Creds.ApiKey)
	q.Set("outputMode", "json")

	body, err := c.watsonClient.MakeRequest("GET", path + "?" + q.Encode(), nil, nil)
	// fmt.Println(string(body))
	if err != nil {
		return err
	}

	var baseResponse BaseResponse
	err = json.Unmarshal(body, &baseResponse)
	if err == nil && strings.EqualFold(baseResponse.Status, "error") {
		return errors.New(baseResponse.StatusInfo)
	}

	return json.Unmarshal(body, out)
}

// Disambiguated contains information about a concept tag returned by some of the Alchmey API endpoints.
type Disambiguated struct {
	Name        string   `json:"name,omitempty"`
	SubType     []string `json:"subType,omitempty"`
	Website     string   `json:"website,omitempty"`
	Geo         string   `json:"geo,omitempty"`
	DBpedia     string   `json:"dbpedia,omitempty"`
	Yago        string   `json:"yago,omitempty"`
	OpenCyc     string   `json:"opencyc,omitempty"`
	Umbel       string   `json:"umbel,omitempty"`
	Freebase    string   `json:"freebase,omitempty"`
	CiaFactbook string   `json:"ciaFactbook,omitempty"`
	Census      string   `json:"census,omitempty"`
	GeoNames    string   `json:"geonames,omitempty"`
	MusicBrainz string   `json:"musicBrainz,omitempty"`
	CrunchBase  string   `json:"crunchbase,omitempty"`
}
