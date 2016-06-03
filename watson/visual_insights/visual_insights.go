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

// Package visual_insights provides an interface to Watson Visual Insights service.
package visual_insights

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v1"
const defaultUrl = "https://gateway.watsonplatform.net/visual-insights-experimental/api"

// Connects to instance of Watson Concept Insights service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "visual_insights"
	}
	if len(cfg.Credentials.Url) == 0 {
		cfg.Credentials.Url = defaultUrl
	}
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	ci.watsonClient = client
	return ci, nil
}

type ClassifierList struct {
	Classifiers []Classifier `json:"classifiers"`
}
type Classifier struct {
	Name string `json:"name"`
}

// Calls 'GET /v1/classifiers' to return a list of all available classifier
func (c Client) ListClassifiers() (ClassifierList, error) {
	b, err := c.watsonClient.MakeRequest("GET", c.version+"/classifiers", nil, nil)
	if err != nil {
		return ClassifierList{}, err
	}
	var cl ClassifierList
	err = json.Unmarshal(b, &cl)
	return cl, err
}

type Summary struct {
	Summary []ClassifierScore `json:"summary,omitempty"`
}

type ClassifierScore struct {
	Name string `json:"name,omitempty"`
	// Score of the classifier on the input images
	Score float64 `json:"score,omitempty"`
}

// Calls 'POST /v1/summary'. It  takes a zip file of images and returns a JSON summary of visual attributes.
func (c Client) Summarize(images_zip io.Reader) (Summary, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	part, err := w.CreateFormFile("images_file", "images.zip")
	if err != nil {
		return Summary{}, err
	}
	_, err = io.Copy(part, images_zip)
	if err != nil {
		return Summary{}, err
	}
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())

	b, err := c.watsonClient.MakeRequest("POST", c.version+"/summary", buf, headers)
	if err != nil {
		return Summary{}, err
	}
	var s Summary
	err = json.Unmarshal(b, &s)
	return s, err

}
