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

// Package conversation provides an interface to Watson Conversation service.
package conversation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v1"
const defaultMinorVersion = "2016-05-19"
const defaultUrl = "https://gateway.watsonplatform.net/conversation-experimental/api"

// Connects to instance of Watson Conversation service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "conversation"
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

type Intent struct {
	Intent     string  `json:"intent,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
}

type IntentExample struct {
	Text     string          `json:"text,omitempty"`
	Entities []EntityExample `json:"entities,omitempty"`
}

type EntityExample struct {
	Entity   string `json:"entity,omitempty"`
	Value    string `json:"value,omitempty"`
	Location []int  `json:"location,omitempty"`
}

type Message struct {
	Input   MessageInput           `json:"input,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type MessageInput struct {
	Text string `json:"text,omitempty"`
}

type MessageOutput struct {
	LogMessages []interface{} `json:"log_messages,omitempty"`
	Text        []string      `json:"text,omitempty"`
	HitNodes    []string      `json:"hit_nodes,omitempty"`
}

type MessageResponse struct {
	Input    MessageInput           `json:"input,omitempty"`
	Intents  []Intent               `json:"intents,omitempty"`
	Entities []EntityExample        `json:"entities,omitempty"`
	Output   MessageOutput          `json:"output,omitempty"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// Calls 'GET /v1/workspaces/{workspace_id}/message' to retrieve response from conversation utterance
func (c Client) Message(workspace_id string, text string, ctx map[string]interface{}) (MessageResponse, error) {
	q := url.Values{}
	q.Set("version", defaultMinorVersion)

	message := &Message{Input: MessageInput{Text: text}, Context: ctx}
	message_json, err := json.Marshal(message)
	if err != nil {
		return MessageResponse{}, err
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")

	body, err := c.watsonClient.MakeRequest("POST", c.version+"/workspaces/"+workspace_id+"/message?"+q.Encode(), bytes.NewReader(message_json), headers)
	if err != nil {
		return MessageResponse{}, err
	}
	var response MessageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return MessageResponse{}, err
	}
	return response, nil
}
