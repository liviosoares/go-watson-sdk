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

// Package document_conversion provides an interface to Watson Document Conversion service.
package document_conversion

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v1"
const defaultMinorVersion = "2015-12-15"
const defaultUrl = "https://gateway.watsonplatform.net/document-conversion/api"

// Connects to instance of Watson Concept Insights service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "document_conversion"
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

const (
	AnswerUnits    = "ANSWER_UNITS"
	NormalizedHtml = "NORMALIZED_HTML"
	NormalizedText = "NORMALIZED_TEXT"
)

func (c Client) Convert(conversion_target string, config_options map[string]interface{}, file io.Reader, content_type string) ([]byte, error) {
	config := make(map[string]interface{})
	for k, v := range config_options {
		config[k] = v
	}
	config["conversion_target"] = conversion_target
	config_json, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	w.WriteField("config", string(config_json))
	// write the file out
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="file"`)
	h.Set("Content-Type", content_type)
	part, err := w.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())

	return c.watsonClient.MakeRequest("POST", c.version+"/convert_document?version="+defaultMinorVersion, buf, headers)
}
