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

package visual_recognition

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.ibm.com/lsoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultVisualRecognitionVersion = "v2"
const defaultMinorVersion = "2015-12-02"

// Connects to instance of Watson Concept Insights service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultVisualRecognitionVersion}
	cfg.Credentials.ServiceName = "visual_recognition"
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
	Id      string    `json:"classifier_id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Owner   string    `json:"owner,omitempty"`
	Status  string    `json:"status,omitempty"`
}

// Calls 'GET /v2/classifiers' to retrieve a list of classifiers
func (c Client) ListClassifiers() (ClassifierList, error) {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)
	q.Set("verbose", "true")
	b, err := c.watsonClient.MakeRequest("GET", c.version+"/classifiers?"+q.Encode(), nil, nil)
	if err != nil {
		return ClassifierList{}, err
	}
	var classifiers ClassifierList
	err = json.Unmarshal(b, &classifiers)
	return classifiers, err
}

// Calls 'GET /v2/classifiers/{classifier_id}' to retrieve classifier details
func (c Client) GetClassifier(id string) (Classifier, error) {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)
	b, err := c.watsonClient.MakeRequest("GET", c.version+"/classifiers/"+id+"?"+q.Encode(), nil, nil)
	if err != nil {
		return Classifier{}, err
	}
	var classifier Classifier
	err = json.Unmarshal(b, &classifier)
	return classifier, err
}

// Calls 'POST /v2/classifiers' to create a classifier. Train a new classifier on the uploaded image data. Upload a compressed (.zip) file of images (.jpg,
// .png, or .gif) with positive examples that show your classifier and another compressed file with negative examples that are similar to but do NOT show your classifier.
func (c Client) CreateClassifier(name string, positive io.Reader, negative io.Reader) (Classifier, error) {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("name", name)
	// positive  examples
	part, err := w.CreateFormFile("positive_examples", "positive.zip")
	if err != nil {
		return Classifier{}, err
	}
	_, err = io.Copy(part, positive)
	if err != nil {
		return Classifier{}, err
	}
	// negative  examples
	part, err = w.CreateFormFile("negative_examples", "negative.zip")
	if err != nil {
		return Classifier{}, err
	}
	_, err = io.Copy(part, negative)
	if err != nil {
		return Classifier{}, err
	}
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())

	b, err := c.watsonClient.MakeRequest("POST", c.version+"/classifiers?"+q.Encode(), buf, headers)
	if err != nil {
		return Classifier{}, err
	}
	var classifier Classifier
	err = json.Unmarshal(b, &classifier)
	return classifier, err
}

// Calls 'DELETE /v2/classifiers/{classifier_id}' to delete a classifier
func (c Client) DeleteClassifier(id string) error {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+"/classifiers/"+id+"?"+q.Encode(), nil, nil)
	return err
}

type ClassifierResult struct {
	Images []ImageResult `json:"images,omitempty"`
}

type ImageResult struct {
	Image  string            `json:"image,omitempty"`
	Scores []ClassifierScore `json:"scores,omitempty"`
}

type ClassifierScore struct {
	// Classifier name
	Name string `json:"name,omitempty"`
	// Classifier id
	Id string `json:"classifier_id,omitempty"`
	// Score of the classifier on the input images.
	Score float64 `json:"score,omitempty"`
}

// Calls 'POST /v2/classify' to classify one or more images.
//
// You can upload a single image or a compressed file (.zip) with multiple images in .jpeg, .png, or .gif format. You can analyze images against all
// classifiers or against an array of classifiers you upload in a JSON file.
// For each image, the response includes a score for a classifier if the score meets the minimum threshold of 0.5. Scores range from 0 - 1 with a higher
// score indicating greater correlation. If no score meets the threshold for an image, no classifiers are returned.
//
// If classifiers is set to 'nil', all classifiers are queried.
func (c Client) Classify(upload io.Reader, classifiers []string) (ClassifierResult, error) {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	if len(classifiers) > 0 {
		ids := struct {
			Ids []string `json:"classifier_ids,omitempty"`
		}{Ids: classifiers}
		classifiers_json, err := json.Marshal(ids)
		if err != nil {
			return ClassifierResult{}, err
		}
		w.WriteField("classifier_ids", string(classifiers_json))
	}
	// positive  examples
	part, err := w.CreateFormFile("images_file", "file.jpg")
	if err != nil {
		return ClassifierResult{}, err
	}
	_, err = io.Copy(part, upload)
	if err != nil {
		return ClassifierResult{}, err
	}
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())

	b, err := c.watsonClient.MakeRequest("POST", c.version+"/classify?"+q.Encode(), buf, headers)
	if err != nil {
		return ClassifierResult{}, err
	}
	var result ClassifierResult
	err = json.Unmarshal(b, &result)
	return result, err
}
