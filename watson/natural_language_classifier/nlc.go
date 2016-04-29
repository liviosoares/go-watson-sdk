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

package natural_language_classifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultNLCVersion = "v1"

// Connects to instance of Watson Natural Language Classifier service
func NewClient(cfg watson.Config) (Client, error) {
	nlc := Client{version: "/" + defaultNLCVersion}
	if len(cfg.Version) > 0 {
		nlc.version = "/" + cfg.Version
	}
	cfg.Credentials.ServiceName = "natural_language_classifier"
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	nlc.watsonClient = client
	return nlc, nil
}

type Classifiers struct {
	Classifiers []Classifier `json:"classifiers"`
}

type Classifier struct {
	// User-supplied name for the classifier ,
	Name string `json:"name,omitempty"`
	// The language used for the classifier ,
	Language string `json:"language,omitempty"`
	// Link to the classifier
	URL string `json:"url"`
	// Unique identifier for this classifier
	ClassifierId string `json:"classifier_id"`
	// Date and time (UTC) the classifier was created
	Created string `json:"created,omitempty"`
}

// Calls 'GET /v1/dialogs' to extract list of dialogs
func (c Client) ListClassifiers() ([]Classifier, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/classifiers", nil, nil)
	if err != nil {
		return nil, err
	}
	var cl Classifiers
	err = json.Unmarshal(body, &cl)
	if err != nil {
		return nil, err
	}
	return cl.Classifiers, nil
}

// Used to specify metadata in classifier creation (see CreateClassifer())
type ClassifierMetadata struct {
	// User-supplied name for the classifier
	Name string `json:"name"`
	// The language used for the classifier
	Language string `json:"language"`
}

type ClassifierStatus struct {
	// User-supplied name for the classifier
	Name string `json:"name,omitempty"`
	// Link to the classifier
	URL string `json:"url,omitempty"`
	// The state of the classifier = ['Non Existent', 'Training', 'Failed', 'Available', 'Unavailable'],
	Status string `json:"status,omitempty"`
	// Unique identifier for this classifier ,
	ClassifierId string `json:"classifier_id,omitempty"`
	// Date and time (UTC) the classifier was created
	Created string `json:"created,omitempty"`
	// Additional detail about the status
	StatusDescription string `json:"status_description,omitempty"`
	// The language used for the classifier
	Language string `json:"language,omitempty"`
}

// Calls 'POST /v1/classifiers' to create a classifier
// The metadata identifies the language of the data, and an optional name to identify the classifier.
// Training data in CSV format. Each text value must have at least one class. The data can include up to 15,000 records.
func (c Client) CreateClassifier(metadata ClassifierMetadata, training_csv io.Reader) (ClassifierStatus, error) {
	j, err := json.Marshal(metadata)
	if err != nil {
		return ClassifierStatus{}, err
	}

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("training_metadata", string(j))
	// now, dump the training CSV file as data for
	part, err := w.CreateFormField("training_data")
	if err != nil {
		return ClassifierStatus{}, err
	}
	_, err = io.Copy(part, training_csv)
	if err != nil {
		return ClassifierStatus{}, err
	}
	w.Close()
	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())

	b, err := c.watsonClient.MakeRequest("POST", c.version+"/classifiers", buf, headers)
	if err != nil {
		return ClassifierStatus{}, err
	}
	fmt.Println(string(b))
	var s ClassifierStatus
	err = json.Unmarshal(b, &s)
	return s, err
}

// Calls 'GET /v1/classifiers/{classifier_id}' to get information about a classifier
func (c Client) GetClassifierStatus(classifier_id string) (ClassifierStatus, error) {
	b, err := c.watsonClient.MakeRequest("GET", c.version+"/classifiers/"+classifier_id, nil, nil)
	if err != nil {
		return ClassifierStatus{}, err
	}
	var s ClassifierStatus
	err = json.Unmarshal(b, &s)
	return s, err
}

// Calls 'DELETE /v1/classifiers/{classifier_id}' to delete a classifier
func (c Client) DeleteClassifier(classifier_id string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+"/classifiers/"+classifier_id, nil, nil)
	return err
}

type Classification struct {
	//  (string, optional): Unique identifier for this classifier ,
	ClassifierId string `json:"classifier_id,omitempty"`
	// Link to the classifier
	URL string `json:"url,omitempty"`
	// The submitted phrase
	Text string `json:"text,omitempty"`
	// The class with the highest confidence
	TopClass string `json:"top_class,omitempty"`
	// An array of up to ten class-confidence pairs sorted in descending order of confidence
	Classes []Class `json:"classes,omitempty"`
}

type Class struct {
	// A decimal percentage that represents the confidence that Watson has in this class. Higher values represent higher confidences.
	Confidence float64 `json:"confidence,omitempty"`
	// Class label
	ClassName string `json:"class_name,omitempty"`
}

func (c Client) Classify(classifier_id string, text string) (Classification, error) {
	q := url.Values{}
	q.Set("text", text)
	b, err := c.watsonClient.MakeRequest("GET", c.version+"/classifiers/"+classifier_id+"/classify?"+q.Encode(), nil, nil)
	if err != nil {
		return Classification{}, err
	}
	var class Classification
	err = json.Unmarshal(b, &class)
	return class, err
}
