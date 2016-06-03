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

// Package language_translation provides an interface to Watson Language Translation service.
package language_translation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v2"
const defaultUrl = "https://gateway.watsonplatform.net/language-translation/api"

// Connects to instance of Watson Natural Language Classifier service
func NewClient(cfg watson.Config) (Client, error) {
	lt := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "language_translation"
	}
	if len(cfg.Credentials.Url) == 0 {
		cfg.Credentials.Url = defaultUrl
	}
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	lt.watsonClient = client
	return lt, nil
}

type ModelList struct {
	Models []Model `json:"models"`
}

type Model struct {
	// A globally unique string that identifies the underlying model that is used for translation. This string contains all the information about
	// source language, target language, domain, and various other related configurations.
	ModelId string `json:"model_id"`
	// If a model is trained by a user, there might be an optional “name” parameter attached during training to help the user identify the model.
	Name string `json:"name"`
	// Source language in two letter language code. Use the five letter code when clarifying between multiple supported languages. When model_id
	// is used directly, it will override the source-target language combination. Also, when a two letter language code is used, but no
	// suitable default is found, it returns an error.
	Source string `json:"source"`
	// Target language in two letter language code.
	Target string `json:"target"`
	// If this model is a custom model, this returns the base model that it is trained on. For a base model, this response value is empty.
	BaseModelId string `json:"base_model_id"`
	// The domain of the translation model.
	Domain string `json:"domain"`
	// Whether this model can be used as a base for customization. ,
	Customizable bool `json:"customizable"`
	// Whether this model is considered a default model and is used when the source and target languages are specified without the model_id. ,
	Default bool `json:"default"`
	// Returns the Bluemix ID of the instance that created the model, or an empty string if it is a model that is trained by IBM. ,
	Owner string `json:"owner"`
	// Availability of a model. = ['available', 'training', 'error']}
	Status string `json:"status"`
}

// Calls 'GET /v2/models' to list available standard and custom models by source or target language
func (c Client) ListModels(options map[string]interface{}) (ModelList, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/models?"+q.Encode(), nil, nil)
	if err != nil {
		return ModelList{}, err
	}
	var models ModelList
	err = json.Unmarshal(body, &models)
	return models, err
}

type TrainingStatus struct {
	// The status of training. Possible responses are: training - training/ is still in progress, error - training did not complete because of an
	// error, or available - training completed and the service is now available to use with your custom translation model ,
	Status string `json:"status"`
	// Returns the base model that this translation model was trained on
	BaseModelId string `json:"bade_model_id"`
}

// Calls 'GET /v2/models/{model_id}' to return the training status of the translation mode
func (c Client) GetModelStatus(model_id string) (TrainingStatus, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/models/"+model_id, nil, nil)
	if err != nil {
		return TrainingStatus{}, err
	}
	var status TrainingStatus
	err = json.Unmarshal(body, &status)
	return status, err
}

// Calls 'DELETE /v2/models/{model_id}' to delete a custom translation mode
func (c Client) DeleteModel(model_id string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+"/models/"+model_id, nil, nil)
	return err
}

// Calls 'POST /v2/models' to uploads a TMX glossary file on top of a domain to customize a translation mode
// base_model_id (Required). Specifies the domain model that is used as the base for the training. To see current supported domain models, use ListModels().
// name The model name. Valid characters are letters, numbers, -, and _. No spaces.
// glossary_type should be one of:
//      "forced_glossary"     TMX file with your customizations. Anything that is specified in this file completely overwrites the domain data
//                                translation. You can upload only one glossary with a file size less than 10 MB per call.
//      "parallel_corpus"     TMX file that contains entries that are treated as a parallel corpus instead of a glossary.
//	"monolingual_corpus"  UTF-8 encoded plain text file that is used to customize the target language model.
//
// Returns the model id for the newly created model
func (c Client) CreateModel(base_model_id string, name string, glossary_type string, glossary io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("base_model_id", base_model_id)
	w.WriteField("name", name)
	part, err := w.CreateFormFile(glossary_type, "glossary.tmx")
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, glossary)
	if err != nil {
		return "", err
	}
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())

	body, err := c.watsonClient.MakeRequest("POST", c.version+"/models", buf, headers)
	if err != nil {
		return "", err
	}
	var model struct {
		ModelId string `json:"model_id"`
	}
	err = json.Unmarshal(body, &model)
	return model.ModelId, err
}

type Response struct {
	// Number of words of the complete input text.
	WordCount int `json:"word_count"`
	// Number of characters of the complete input text.
	CharacterCount int `json:"character_count`
	// List of translation output in UTF-8, corresponding to the list of input text.
	Translations []Translation `json:"translations"`
}
type Translation struct {
	// Translation output in UTF-8.
	Translation string
}

// Calls 'POST /v2/translate' to translates the input text from the source language to the target language
// model_id  The unique model_id of the translation model that is used to translate text. The model_id inherently specifies source language, target
//           language, and domain. If the model_id is specified, there is no need for the source and target parameters, and the values are ignored.
// source     Used in combination with target as an alternative way to select the model for translation. When target and source are set, and model_id is not
//            set, the system chooses a default model with the right language pair to  translate (usually the model based on the news domain).
// target     Used in combination with source as an alternative way to select which model is used for translation. When target and source are set, and model_id
//            is not set, the system chooses a default model with the right language pair to translate (usually the model based on the news domain).
//
// returns a Response object
func (c Client) Translate(text string, source string, target string, model_id string) (Response, error) {
	req := map[string]interface{}{
		"text": text,
	}
	if len(model_id) > 0 {
		req["model_id"] = model_id
	}
	if len(source) > 0 {
		req["source"] = source
	}
	if len(target) > 0 {
		req["target"] = target
	}

	req_json, err := json.Marshal(req)
	if err != nil {
		return Response{}, err
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")
	body, err := c.watsonClient.MakeRequest("POST", c.version+"/translate", bytes.NewReader(req_json), headers)
	if err != nil {
		return Response{}, err
	}
	var response Response
	err = json.Unmarshal(body, &response)
	return response, err
}

type IdentifiableLanguageList struct {
	// A list of all languages that the service can identify.
	Languages []IdentifiableLanguage `json:"languages"`
}

type IdentifiableLanguage struct {
	// The code for an identifiable language.
	Language string `json:"language"`
	// The name of the identifiable language.
	Name string `json:"name"`
}

// Calls 'GET /v2/identifiable_languages' to lists all languages that can be identified by the API
func (c Client) ListIdentifiableLanguages() (IdentifiableLanguageList, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/identifiable_languages", nil, nil)
	if err != nil {
		return IdentifiableLanguageList{}, err
	}
	var list IdentifiableLanguageList
	err = json.Unmarshal(body, &list)
	return list, err
}

type IdentifiedLanguages struct {
	// A ranking of identified languages with confidence scores.
	Languages []IdentifiedLanguage `json:"languages"`
}

type IdentifiedLanguage struct {
	// The code for an identified language.
	Language string `json:"language"`
	// The confidence score for the identified language.
	Confidence float64 `json:"confidence"`
}

func (c Client) IdentifyLanguage(text string) (languages IdentifiedLanguages, err error) {
	headers := make(http.Header)
	headers.Set("Accept", "application/json")
	headers.Set("Content-Type", "text/plain")
	body, err := c.watsonClient.MakeRequest("POST", c.version+"/identify", strings.NewReader(text), headers)
	if err != nil {
		return IdentifiedLanguages{}, err
	}
	var langs IdentifiedLanguages
	err = json.Unmarshal(body, &langs)
	return langs, err
}
