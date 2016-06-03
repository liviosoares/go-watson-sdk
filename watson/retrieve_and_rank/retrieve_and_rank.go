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

// Package retrieve_and_rank provides an interface to Watson Retrieve and Rank service.
package retrieve_and_rank

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

const defaultMajorVersion = "v1"
const defaultUrl = "https://gateway.watsonplatform.net/retrieve-and-rank/api"

// Connects to instance of Watson Concept Insights service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "retrieve_and_rank"
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

type ClusterList struct {
	Clusters []Cluster `json:"clusters"`
}

type Cluster struct {
	// Unique identifier for this cluster
	Id string `json:"solr_cluster_id,omitempty"`
	// Name that identifies the cluster
	Name string `json:"cluster_name,omitempty"`
	// Size of the cluster to create
	Size string `json:"cluster_size,omitempty"`
	// The state of the cluster = ['NOT_AVAILABLE', 'READY']}
	Status string `json:"solr_cluster_status,omitempty"`
}

// Calls 'GET /v1/solr_clusters' to get a list Solr clusters
func (c Client) ListClusters() (ClusterList, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/solr_clusters", nil, nil)
	if err != nil {
		return ClusterList{}, err
	}
	var response ClusterList
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'POST /v1/solr_clusters' to create Solr cluster
// 'size' is corresponds to the cluster to create; ranges from 1 to 7. Use a zero value to create a small free-size cluster for testing. You can create
//        only one free-size cluster for each service instance.
func (c Client) CreateCluster(name string, size int) (Cluster, error) {
	var def struct {
		ClusterName string `json:"cluster_name"`
		ClusterSize string `json:"cluster_size,omitempty"`
	}
	def.ClusterName = name
	if size > 0 {
		def.ClusterSize = fmt.Sprintf("%d", size)
	}
	def_json, err := json.Marshal(def)
	if err != nil {
		return Cluster{}, err
	}

	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")
	body, err := c.watsonClient.MakeRequest("POST", c.version+"/solr_clusters", bytes.NewReader(def_json), headers)
	if err != nil {
		return Cluster{}, err
	}
	var response Cluster
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'DELETE /v1/solr_clusters/{solr_cluster_id}' to delete a Solr cluster
func (c Client) DeleteCluster(id string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+"/solr_clusters/"+id, nil, nil)
	return err
}

// Calls 'GET /v1/solr_clusters/{solr_cluster_id}' to retrieve information about a Solr cluster
func (c Client) GetCluster(id string) (Cluster, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/solr_clusters/"+id, nil, nil)
	if err != nil {
		return Cluster{}, err
	}
	var response Cluster
	err = json.Unmarshal(body, &response)
	return response, err
}

// type Configs struct {
// 	Configs []struct {
//		ConfigName string `json:"config_name"`
//	} `json:"solr_configs"`
// }

type Configs struct {
	Configs []string `json:"solr_configs"`
}

// Calls 'GET /v1/solr_clusters/{solr_cluster_id}/config' to list Solr configurations
func (c Client) ListConfigs(id string) (Configs, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/solr_clusters/"+id+"/config", nil, nil)
	if err != nil {
		return Configs{}, err
	}
	var response Configs
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'POST /v1/solr_clusters/{solr_cluster_id}/config/{config_name}' to upload Solr configuration
func (c Client) UploadConfig(solr_id string, config_name string, zipReader io.Reader) error {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/zip")
	headers.Set("Accept", "application/json")
	_, err := c.watsonClient.MakeRequest("POST", c.version+"/solr_clusters/"+solr_id+"/config/"+config_name, zipReader, headers)
	return err
}

// Calls 'DELETE /v1/solr_clusters/{solr_cluster_id}/config/{config_name}' to delete Solr configuration
func (c Client) DeleteConfig(solr_id string, config_name string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+"/solr_clusters/"+solr_id+"/config/"+config_name, nil, nil)
	return err
}

// Calls 'GET /v1/solr_clusters/{solr_cluster_id}/config/{config_name}' to get Solr configuration
func (c Client) GetConfig(solr_id string, config_name string) ([]byte, error) {
	return c.watsonClient.MakeRequest("GET", c.version+"/solr_clusters/"+solr_id+"/config/"+config_name, nil, nil)
}

// Calls 'POST /v1/solr_clusters/{solr_cluster_id}/solr/admin/collections' to forward collection requests to Solr (CREATE, DELETE, LIST)
// 'action' : Operation to carry out. CREATE creates a Solr collection, DELETE removes a collection, and LIST returns the names of the collections in the cluster.
func (c Client) doCollectionRequest(solr_id string, user_options map[string]interface{}, default_options map[string]interface{}) ([]byte, error) {
	q := url.Values{}
	for k, v := range user_options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	for k, v := range default_options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	return c.watsonClient.MakeRequest("POST", c.version+"/solr_clusters/"+solr_id+"/solr/admin/collections?"+q.Encode(), nil, nil)
}

// Calls 'POST /v1/solr_clusters/{solr_cluster_id}/solr/admin/collections' to with 'CREATE' action to create Solr collection.
func (c Client) CreateCollection(solr_id string, collection_name string, config_name string, options map[string]interface{}) ([]byte, error) {
	default_options := map[string]interface{}{
		"action":                "CREATE",
		"name":                  collection_name,
		"collection.configName": config_name,
	}
	return c.doCollectionRequest(solr_id, options, default_options)
}

// Calls 'POST /v1/solr_clusters/{solr_cluster_id}/solr/admin/collections' to with 'DELETE' action to create Solr collection.
func (c Client) DeleteCollection(solr_id string, collection_name string, options map[string]interface{}) ([]byte, error) {
	default_options := map[string]interface{}{
		"action": "DELETE",
		"name":   collection_name,
	}
	return c.doCollectionRequest(solr_id, options, default_options)
}

// Calls 'POST /v1/solr_clusters/{solr_cluster_id}/solr/admin/collections' to with 'LIST' action to create Solr collection.
func (c Client) ListCollections(solr_id string, options map[string]interface{}) ([]byte, error) {
	default_options := map[string]interface{}{
		"action": "LIST",
	}
	return c.doCollectionRequest(solr_id, options, default_options)
}

// Calls 'POST /v1/solr_clusters/{solr_cluster_id}/solr/{collection_name}/update' to index document
func (c Client) Update(solr_id string, collection_name string, content_type string, reader io.Reader, options map[string]interface{}) ([]byte, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	headers := make(http.Header)
	headers.Set("Content-Type", content_type)
	fmt.Println(c.version + "/solr_clusters/" + solr_id + "/solr/" + collection_name + "/update")
	return c.watsonClient.MakeRequest("POST", c.version+"/solr_clusters/"+solr_id+"/solr/"+collection_name+"/update?"+q.Encode(), reader, headers)
}

// Calls 'GET /v1/solr_clusters/{solr_cluster_id}/solr/{collection_name}/select' to execute Solr standard search query
func (c Client) Search(solr_id string, collection_name string, query string, options map[string]interface{}) ([]byte, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("q", query)
	return c.watsonClient.MakeRequest("GET", c.version+"/solr_clusters/"+solr_id+"/solr/"+collection_name+"/select?"+q.Encode(), nil, nil)
}

type RankerList struct {
	// InfoPayload], optional): The rankers available to the user. Returns an empty array if no rankers are available.
	Rankers []Ranker `json:"rankers"`
}

type Ranker struct {
	// Unique identifier for this ranker ,
	RankerId string `json:"ranker_id,omitempty"`
	// Link to the ranker ,
	Url string `json:"url,omitempty"`
	// User-supplied name for the ranker ,
	Name string `json:"name,omitempty"`
	// Date and time (UTC) the ranker was created
	Created string `json:"created,omitempty"`
	// The state of the ranker = ['Non Existent', 'Training', 'Failed', 'Available', 'Unavailable']
	Status string `json:"status,omitempty`
	// Additional detail about the status
	StatusDescription string `json:"status_description,omitempty"`
}

// Calls 'GET /v1/rankers' to  list rankers
func (c Client) ListRankers() (RankerList, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/rankers", nil, nil)
	if err != nil {
		return RankerList{}, err
	}
	var response RankerList
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'POST /v1/rankers' to create a ranker
func (c Client) CreateRanker(name string, trainingData io.Reader) (Ranker, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	part, err := w.CreateFormFile("file", "training_data")
	if err != nil {
		return Ranker{}, err
	}
	_, err = io.Copy(part, trainingData)
	if err != nil {
		return Ranker{}, err
	}
	var meta struct {
		Name string `json:"name"`
	}
	meta.Name = name
	meta_json, err := json.Marshal(meta)
	if err != nil {
		return Ranker{}, err
	}
	w.WriteField("training_metadata", string(meta_json))
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())
	headers.Set("Accept", "application/json")

	body, err := c.watsonClient.MakeRequest("POST", c.version+"/rankers", buf, headers)
	if err != nil {
		return Ranker{}, err
	}
	var response Ranker
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'GET /v1/rankers/{ranker_id}' to get information about a ranker
func (c Client) GetRanker(ranker_id string) (Ranker, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/rankers/"+ranker_id, nil, nil)
	if err != nil {
		return Ranker{}, err
	}
	var response Ranker
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'DELETE /v1/rankers/{ranker_id}' to delete ranker
func (c Client) DeleteRanker(ranker_id string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+"/rankers/"+ranker_id, nil, nil)
	return err
}

type RankerOutput struct {
	//  Unique identifier for this ranker
	RankerId string `json:"ranker_id,omitempty"`
	// Name of this ranker
	Name string `json:"name"`
	// Link to the ranker
	Url string `json:"url,omitempty"`
	// The class with the highest confidence ,
	TopAnswer string `json:"top_answer,omitempty"`
	// An array of up to ten class-confidence pairs sorted in descending order of confidence
	Answers []Answer `json:"answers,omitempty"`
}

type Answer struct {
	// Answer label ,
	AnswerId string `json:"answer_id,omitempty"`
	// A decimal percentage that represents the confidence that Watson has in this class. Higher values represent higher confidences. ,
	Score float64 `json:"score,omitempty"`
	// A decimal percentage that represents the confidence that Watson has in this class. Higher values represent higher confidences.
	Confidence float64 `json:"confidence,omitempty"`
}

// Calls 'POST /v1/rankers/{ranker_id}/rank' to rank a series of search queries
// Returns the top answer and a list of ranked answers with their ranked scores and confidence values.
// Use this method to return answers when you train the ranker with custom features. However, in most cases, you can use the Search and rank method.
func (c Client) Rank(ranker_id string, answerData io.Reader) (RankerOutput, error) {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	part, err := w.CreateFormFile("file", "answer_data")
	if err != nil {
		return RankerOutput{}, err
	}
	_, err = io.Copy(part, answerData)
	if err != nil {
		return RankerOutput{}, err
	}
	w.Close()

	headers := make(http.Header)
	headers.Set("Content-Type", w.FormDataContentType())
	headers.Set("Accept", "application/json")

	body, err := c.watsonClient.MakeRequest("POST", c.version+"/rankers/"+ranker_id+"/rank", buf, headers)
	if err != nil {
		return RankerOutput{}, err
	}
	var response RankerOutput
	err = json.Unmarshal(body, &response)
	return response, err
}

// Calls 'GET /v1/solr_clusters/{solr_cluster_id}/solr/{collection_name}/fcselect' to execute Solr search query with reranking.
// See documentation at: https://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/retrieve-and-rank/api/v1/#query_ranker 
func (c Client) RankAndSearch(solr_id string, collection_name string, ranker_id string, query string, options map[string]interface{}) ([]byte, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("ranker_id", ranker_id)
	q.Set("q", query)
	return c.watsonClient.MakeRequest("GET", c.version+"/solr_clusters/"+solr_id+"/solr/"+collection_name+"/fcselect?"+q.Encode(), nil, nil)
}
