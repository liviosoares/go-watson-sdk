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

// Package dialog provides an interface to Watson Dialog service.
package dialog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultDialogVersion = "v1"

// Connects to instance of Watson Dialog service
func NewClient(cfg watson.Config) (Client, error) {
	dialog := Client{version: "/" + defaultDialogVersion}
	if len(cfg.Version) > 0 {
		dialog.version = "/" + cfg.Version
	}
	cfg.Credentials.ServiceName = "dialog"
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	dialog.watsonClient = client
	return dialog, nil
}

type dialogs struct {
	Dialogs       []Dialog `json:"dialogs,omitempty"`
	LanguagePacks []Dialog `json:"language_packs,omitempty"`
}
type Dialog struct {
	DialogId string `json:"dialog_id,omitempty"`
	Name     string `json:"name,omitempty"`
}

// Calls 'GET /v1/dialogs' to extract list of dialogs
func (d Client) ListDialogs() ([]Dialog, error) {
	body, err := d.watsonClient.MakeRequest("GET", d.version+"/dialogs", nil, nil)
	if err != nil {
		return nil, err
	}
	var dialogs dialogs
	err = json.Unmarshal(body, &dialogs)
	if err != nil {
		return nil, err
	}
	return dialogs.Dialogs, nil
}

// Calls 'GET /v1/dialogs' to extract list of language packs
func (d Client) ListLanguagePacks() ([]Dialog, error) {
	body, err := d.watsonClient.MakeRequest("GET", d.version+"/dialogs", nil, nil)
	if err != nil {
		return nil, err
	}
	var dialogs dialogs
	err = json.Unmarshal(body, &dialogs)
	if err != nil {
		return nil, err
	}
	return dialogs.LanguagePacks, nil
}

// Calls 'POST /v1/dialogs' to create new dialog. template is a reader to the
// content.  The file content type is determined by the filename extension:
//   .mct for encrypted Dialog account file,
//   .json fo Watson Dialog document JSON format
//   .xml for Watson Dialog document XML format
func (d Client) CreateDialog(name string, filename string, data io.Reader) (string, error) {
	return d.createOrUpdateDialog("", name, filename, data)
}

// Calls 'PUT /v1/dialogs/{dialog_id}' to update an existing Dialog
func (d Client) UpdateDialog(id string, filename string, data io.Reader) error {
	_, err := d.createOrUpdateDialog(id, "", filename, data)
	return err
}

func (d Client) createOrUpdateDialog(id string, name string, filename string, data io.Reader) (string, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	if len(id) == 0 {
		// first, write out the name of the dialog
		part, err := w.CreateFormField("name")
		if err != nil {
			return "", err
		}
		part.Write([]byte(name))
	}
	// now, dump the template file as data for the dialog
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, data)
	if err != nil {
		return "", err
	}
	w.Close()
	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())
	var b []byte
	if len(id) == 0 {
		b, err = d.watsonClient.MakeRequest("POST", d.version+"/dialogs", buf, headers)
	} else {
		b, err = d.watsonClient.MakeRequest("PUT", d.version+"/dialogs/"+id, buf, headers)
	}
	if err != nil {
		return "", err
	}
	var respId struct {
		Id string `json:"id"`
	}
	err = json.Unmarshal(b, &respId)
	if err != nil {
		return "", err
	}
	return respId.Id, nil
}

// Calls 'GET /v1/dialogs/{dialog_id}' to download dialog files. Valid concept_type values
// are 'application/octet-stream', 'application/wds+json', and 'application/wds+xml'
func (d Client) DownloadDialog(id string, content_type string) ([]byte, error) {
	headers := make(http.Header)
	headers.Set("Accept", content_type)
	return d.watsonClient.MakeRequest("GET", d.version+"/dialogs/"+id, nil, headers)
}

// Calls 'DELETE /v1/dialogs/{dialog_id}' to remove dialog, including all associated data
func (d Client) DeleteDialog(id string) error {
	_, err := d.watsonClient.MakeRequest("DELETE", d.version+"/dialogs/"+id, nil, nil)
	return err
}

type Node struct {
	Node    string `json:"node,omitempty"`
	Content string `json:"content,omitempty"`
}

type nodesStruct struct {
	Items []Node `json:"items,omitempty"`
}

func (d Client) GetNodes(id string, options map[string]string) ([]Node, error) {
	b, err := d.watsonClient.MakeRequest("GET", d.version+"/dialogs/"+id+"/content", nil, nil)
	if err != nil {
		return nil, err
	}
	var n nodesStruct
	err = json.Unmarshal(b, &n)
	if err != nil {
		return nil, err
	}
	return n.Items, nil
}

func (d Client) UpdateNodes(id string, nodes []Node) error {
	n := nodesStruct{Items: nodes}
	body, err := json.Marshal(n)
	if err != nil {
		return err
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	_, err = d.watsonClient.MakeRequest("PUT", d.version+"/dialogs/"+id+"/content", bytes.NewReader(body), headers)
	return err
}

type ConversationResponse struct {
	Response       []string `json:"response"`
	Input          string   `json:"input"`
	ConversationId uint64   `json:"conversation_id"`
	ClientId       uint64   `json:"client_id"`
	Confidence     float64  `json:"confidence"`
}

func (d Client) StartConversation(dialog_id string) (ConversationResponse, error) {
	return d.UpdateConversation(dialog_id, 0, 0, "")
}

func (d Client) UpdateConversation(dialog_id string, conversation_id uint64, client_id uint64, input string) (ConversationResponse, error) {
	values := url.Values{}
	values.Set("input", input)
	if conversation_id != 0 {
		values.Set("conversation_id", fmt.Sprintf("%v", conversation_id))
	}
	if client_id != 0 {
		values.Set("client_id", fmt.Sprintf("%v", client_id))
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/x-www-form-urlencoded")
	b, err := d.watsonClient.MakeRequest("POST", d.version+"/dialogs/"+dialog_id+"/conversation", strings.NewReader(values.Encode()), headers)
	if err != nil {
		return ConversationResponse{}, err
	}
	var r ConversationResponse
	err = json.Unmarshal(b, &r)
	return r, err
}

const dialogTimeFormat = "2006-01-02 15:04:05"

type ConversationHistory struct {
	Conversations []Conversation `json:"conversations"`
}

type Conversation struct {
	HitNodes       []HitNode     `json:"hit_nodes,omitempty"`
	ConversationId uint64        `json:"conversation_id,omitempty"`
	ClientId       uint64        `json:"client_id,omitempty"`
	Messages       []Message     `json:"messages,omitempty"`
	Profile        []ProfileItem `json:"profile,omitempty"`
}

type HitNode struct {
	Details string `json:"details,omitempty"`
	Label   string `json:"label,omitempty"`
	Type    string `json:"type,omitempty"`
	NodeId  int64  `json:"node_id,omitempty"`
}

type Message struct {
	Text       string `json:"text,omitempty"`
	DateTime   string `json:"date_time,omitempty"`
	FromClient bool   `json:"from_client,string,omitempty"`
}

type ProfileItem struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (d Client) GetConversationHistory(dialog_id string, from time.Time, to time.Time, offset int, limit int) (ConversationHistory, error) {
	query := url.Values{}
	query.Set("date_from", from.Format(dialogTimeFormat))
	query.Set("date_to", to.Format(dialogTimeFormat))
	if offset != 0 {
		query.Set("offset", fmt.Sprintf("%v", offset))
	}
	if limit != 0 {
		query.Set("limit", fmt.Sprintf("%v", limit))
	}
	headers := make(http.Header)
	headers.Set("Accept", "application/json")
	b, err := d.watsonClient.MakeRequest("GET", d.version+"/dialogs/"+dialog_id+"/conversation?"+query.Encode(), nil, headers)
	if err != nil {
		return ConversationHistory{}, err
	}
	var hist ConversationHistory
	err = json.Unmarshal(b, &hist)
	return hist, err
}

type NameValues struct {
	NameValues []NameValue `json:"name_values"`
	ClientId   uint64      `json:"client_id,omitempty"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Calls 'GET /v1/dialog/{dialog_id}/profile' to get the values for profile variables for a given client id
func (d Client) GetProfileVariables(dialog_id string, client_id uint64) (NameValues, error) {
	query := url.Values{}
	query.Set("client_id", fmt.Sprintf("%v", client_id))
	b, err := d.watsonClient.MakeRequest("GET", d.version+"/dialogs/"+dialog_id+"/profile?"+query.Encode(), nil, nil)
	if err != nil {
		return NameValues{}, err
	}
	var nameValues NameValues
	err = json.Unmarshal(b, &nameValues)
	return nameValues, err
}

// Calls 'PUT /v1/dialogs/{dialog_id}/profile' to set the values for profile variables
// Profile variables needs to be already explicitly defined in the application.
func (d Client) SetProfileVariable(dialog_id string, nv NameValues) error {
	body, err := json.Marshal(nv)
	if err != nil {
		return err
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	_, err = d.watsonClient.MakeRequest("PUT", d.version+"/dialogs/"+dialog_id+"/profile", bytes.NewReader(body), headers)
	return err
}
