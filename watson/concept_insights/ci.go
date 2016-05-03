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

// Package concept_insights provides an interface to Watson Concept Insights service.
package concept_insights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultConceptInsightsVersion = "v2"

// Connects to instance of Watson Concept Insights service
func NewClient(cfg watson.Config) (Client, error) {
	ci := Client{version: "/" + defaultConceptInsightsVersion}
	cfg.Credentials.ServiceName = "concept_insights"
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	ci.watsonClient = client
	return ci, nil
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}
type Account struct {
	Id string `json:"account_id"`
}

// Calls 'GET /v2/accounts' to retrieves a Concept Insights account identifier,
// to be used as top-level resource name for other APIs.
func (c Client) ListAccounts() (Accounts, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/accounts", nil, nil)
	if err != nil {
		return Accounts{}, err
	}
	var accounts Accounts
	err = json.Unmarshal(body, &accounts)
	return accounts, err
}

type Graphs struct {
	Graphs []string `json:"graphs"`
}

// Calls 'GET /v2/graphs' to retrieves the list available graphs for the authenticated user
func (c Client) ListGraphs() (Graphs, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/graphs", nil, nil)
	if err != nil {
		return Graphs{}, err
	}
	var graphs Graphs
	err = json.Unmarshal(body, &graphs)
	return graphs, err
}

type Concept struct {
	Id        string   `json:"id"`
	Label     string   `json:"label,omitempty"`
	Abstract  string   `json:"abstract,omitempty"`
	Link      string   `json:"link,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
	Type      string   `json:"type,omitempty"`
	Ontology  []string `json:"ontology,omitempty"`
}

// Calls 'GET /v2/graphs/{graph_id}/concept/{concept}' to retrieve information for a specific concept node in a graph.
// Note that concept_id is expected to be the fully path to the concept; for example: "/graphs/wikipedia/en-20120601/concepts/IBM_Watson"
func (c Client) GetConcept(concept_id string) (Concept, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+concept_id, nil, nil)
	if err != nil {
		return Concept{}, err
	}
	var concept Concept
	err = json.Unmarshal(body, &concept)
	return concept, err
}

type LabelMatches struct {
	Matches []Concept
}

// Calls 'GET /v2/graphs/{graph_id}/label_search' to search concepts in a concept graph looking for partial matches on the concept label field.
// When the 'prefix' parameter is set to true, the main use of this method is to build query boxes that offer auto-complete, to allow users
// to select valid concepts.
func (c Client) SearchConceptByLabel(graph_id string, query string, options map[string]interface{}) (LabelMatches, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("query", query)
	body, err := c.watsonClient.MakeRequest("GET", c.version+graph_id+"/label_search?"+q.Encode(), nil, nil)
	if err != nil {
		return LabelMatches{}, err
	}
	var matches LabelMatches
	err = json.Unmarshal(body, &matches)
	return matches, err
}

type ConceptMatches struct {
	Concepts []ConceptScore `json:"concepts"`
}

type ConceptScore struct {
	Concept Concept `json:"concept"`
	Score   float64 `json:"score,omitempty"`
}

// Calls 'GET /v2/graphs/{graph_id}/related_concepts' to retrieves concepts that are related to a concept
func (c Client) GetRelatedConcepts(graph_id string, concepts []string, options map[string]interface{}) (ConceptMatches, error) {
	concepts_json, err := json.Marshal(concepts)
	if err != nil {
		return ConceptMatches{}, err
	}
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("concepts", string(concepts_json))
	body, err := c.watsonClient.MakeRequest("GET", c.version+graph_id+"/related_concepts?"+q.Encode(), nil, nil)
	if err != nil {
		return ConceptMatches{}, err
	}
	var conceptMatches ConceptMatches
	err = json.Unmarshal(body, &conceptMatches)
	return conceptMatches, err
}

type Annotations struct {
	Annotations []Annotation `json:"annotations"`
}

type Annotation struct {
	Concept    Concept `json:"concept"`
	Score      float64 `json:"score,omitempty"`
	TextIndex  []int   `json:"text_index,omitempty"`
	PartsIndex int     `json:"parts_index,omitempty"`
}

// Calls 'POST /v2/graphs/{graph_id}/annotate_text' to identify concept mentions in a piece of text (concept extraction).
func (c Client) AnnotateText(graph_id string, text io.Reader, content_type string) (Annotations, error) {
	headers := make(http.Header)
	headers.Set("Content-Type", content_type)
	body, err := c.watsonClient.MakeRequest("POST", c.version+graph_id+"/annotate_text", text, headers)
	if err != nil {
		return Annotations{}, nil
	}
	var annotations Annotations
	err = json.Unmarshal(body, &annotations)
	return annotations, err
}

type ConceptScores struct {
	Scores []ConceptScore
}

// Calls 'GET /v2/graphs/{graph_id}/concepts/{concept}/relation_scores' to return a list of scores that denotes how
// related a source concept is to a list of individual concepts
func (c Client) GetRelationScore(from_concept_id string, to_concepts []string) (ConceptScores, error) {
	concepts_json, err := json.Marshal(to_concepts)
	if err != nil {
		return ConceptScores{}, err
	}
	q := url.Values{}
	q.Set("concepts", string(concepts_json))
	body, err := c.watsonClient.MakeRequest("GET", c.version+from_concept_id+"/relation_scores?"+q.Encode(), nil, nil)
	if err != nil {
		return ConceptScores{}, err
	}
	var concepts ConceptScores
	err = json.Unmarshal(body, &concepts)
	return concepts, err
}

type CorpusAuthUser struct {
	AccountId  string `json:"uid"`
	Permission string `json:"permission"`
}
type Corpus struct {
	Id     string           `json:"id,omitempty"`
	Access string           `json:"access"`
	Users  []CorpusAuthUser `json:"users"`
}

type CorporaList struct {
	Corpora []Corpus `json:"corpora"`
}

// Call 'GET /v2/corpora' to retrieve the available corpora
func (c Client) ListCorpora() (CorporaList, error) {
	return c.listCorpora("")
}

// Call 'GET /v2/corpora/{account_id}' to retrieve the available corpora
func (c Client) ListCorporaByAccountId(account_id string) (CorporaList, error) {
	return c.listCorpora("/" + account_id)
}

func (c Client) listCorpora(prefix string) (CorporaList, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/corpora"+prefix, nil, nil)
	if err != nil {
		return CorporaList{}, err
	}
	var corpora CorporaList
	err = json.Unmarshal(body, &corpora)
	return corpora, err
}

func (c Client) GetCorpus(corpus_id string) (Corpus, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id, nil, nil)
	if err != nil {
		return Corpus{}, err
	}
	var corpus Corpus
	err = json.Unmarshal(body, &corpus)
	return corpus, err
}

func (c Client) DeleteCorpus(corpus_id string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+corpus_id, nil, nil)
	return err
}

func (c Client) CreateCorpus(corpus_id string, corpus Corpus) error {
	corpus_json, err := json.Marshal(corpus)
	if err != nil {
		return err
	}
	_, err = c.watsonClient.MakeRequest("PUT", c.version+corpus_id, bytes.NewReader(corpus_json), nil)
	return err
}

func (c Client) UpdateCorpus(corpus_id string, corpus Corpus) error {
	corpus_json, err := json.Marshal(corpus)
	if err != nil {
		return err
	}
	_, err = c.watsonClient.MakeRequest("POST", c.version+corpus_id, bytes.NewReader(corpus_json), nil)
	return err
}

type CorpusProcessingState struct {
	Id          string            `json:"id"`
	Documents   int               `json:"documents,omitempty"`
	LastUpdated time.Time         `json:"last_updated,omitempty"`
	BuildStatus CorpusBuildStatus `json:"build_status,omitempty"`
}

type CorpusBuildStatus struct {
	Error      int `json:"error"`
	Processing int `json:"processing"`
	Ready      int `json:"ready"`
}

func (c Client) GetCorpusProcessingState(corpus_id string) (CorpusProcessingState, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/processing_state", nil, nil)
	if err != nil {
		return CorpusProcessingState{}, err
	}
	var state CorpusProcessingState
	err = json.Unmarshal(body, &state)
	return state, err
}

type CorpusStats struct {
	Id          string        `json:"id"`
	LastUpdated time.Time     `json:"last_updated,omitempty"`
	TopTags     CorpusTopTags `json:"top_tags,omitempty"`
}

type CorpusTopTags struct {
	Documents               int            `json:"documents,omitempty"`
	TotalTags               int            `json:"total_tags,omitempty"`
	UniqueTags              int            `json:"unique_tags,omitempty"`
	Tags                    []CorpusTag    `json:"tags,omitempty"`
	CorpusTagsHistogram     map[string]int `json:"corpus_tags_histogram,omitempty"`
	DocumentTagsHistogram   map[string]int `json:"document_tags_histogram,omitempty"`
	DocumentLengthHistogram map[string]int `json:"document_length_histogram,omitempty"`
}

type CorpusTag struct {
	Concept string `json:"concept,omitempty"`
	Count   int    `json:"count,omitempty"`
}

func (c Client) GetCorpusStats(corpus_id string) (CorpusStats, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/stats", nil, nil)
	if err != nil {
		return CorpusStats{}, err
	}
	var stats CorpusStats
	err = json.Unmarshal(body, &stats)
	return stats, err
}

// Calls 'GET /v2/corpora/{corpus_id}/label_search' to search for documents and concepts by using partial matches on the label(s) fields
func (c Client) SearchCorpusByLabel(corpus_id string, query string, options map[string]interface{}) (LabelMatches, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("query", query)
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/label_search?"+q.Encode(), nil, nil)
	if err != nil {
		return LabelMatches{}, err
	}
	var matches LabelMatches
	err = json.Unmarshal(body, &matches)
	return matches, err
}

func (c Client) GetCorpusRelatedConcepts(corpus_id string, options map[string]interface{}) (ConceptMatches, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/related_concepts?"+q.Encode(), nil, nil)
	if err != nil {
		return ConceptMatches{}, err
	}
	var conceptMatches ConceptMatches
	err = json.Unmarshal(body, &conceptMatches)
	return conceptMatches, err

}

func (c Client) GetCorpusRelationScores(corpus_id string, to_concepts []string) (ConceptScores, error) {
	concepts_json, err := json.Marshal(to_concepts)
	if err != nil {
		return ConceptScores{}, err
	}
	q := url.Values{}
	q.Set("concepts", string(concepts_json))
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/relation_scores?"+q.Encode(), nil, nil)
	if err != nil {
		return ConceptScores{}, err
	}
	var concepts ConceptScores
	err = json.Unmarshal(body, &concepts)
	return concepts, err
}

type SemanticResults struct {
	// Concepts used to create query
	QueryConcepts []Concept `json:"query_concepts"`
	// Results of search
	Results []SemanticResult `json:"results"`
}

type SemanticResult struct {
	Id     string   `json:"id"`
	Label  string   `json:"label,omitempty"`
	Labels []string `json:"labels,omitempty"`
	// Annotation from document that explains the conceptual relation of this document to a search query
	ExplanationTags []Annotation `json:"explanation_tags,omitempty"`
}

// Calls 'GET /v2/corpora/{corpus_id}/conceptual_search' to perform a conceptual search within a corpus
func (c Client) GetRelatedDocuments(corpus_id string, ids []string, options map[string]interface{}) (SemanticResults, error) {
	ids_json, err := json.Marshal(ids)
	if err != nil {
		return SemanticResults{}, err
	}
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("ids", string(ids_json))
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/conceptual_search?"+q.Encode(), nil, nil)
	if err != nil {
		return SemanticResults{}, err
	}
	var results SemanticResults
	err = json.Unmarshal(body, &results)
	return results, err
}

type DocumentList struct {
	Documents []string `json:"documents"`
}

func (c Client) ListDocuments(corpus_id string, options map[string]interface{}) (DocumentList, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	body, err := c.watsonClient.MakeRequest("GET", c.version+corpus_id+"/documents?"+q.Encode(), nil, nil)
	if err != nil {
		return DocumentList{}, err
	}
	var docList DocumentList
	err = json.Unmarshal(body, &docList)
	return docList, err
}

type Document struct {
	Id           string                 `json:"id"`
	Label        string                 `json:"label,omitempty"`
	Labels       []string               `json:"labels,omitempty"`
	LastModified time.Time              `json:"last_modified,omitempty"`
	UserFields   map[string]interface{} `json:"user_fields,omitempty"`
	TtlHours     int                    `json:"ttl_hours,omitempty"`
	ExpiresOn    time.Time              `json:"expires_on,omitempty"`
	Parts        []DocumentPart         `json:"parts"`
}

type DocumentPart struct {
	Data        string `json:"data,omitempty"`
	Name        string `json:"name,omitempty"`
	ContentType string `json:"content-type,omitempty"`
}

func (c Client) GetDocument(document_id string) (Document, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+document_id, nil, nil)
	if err != nil {
		return Document{}, err
	}
	var doc Document
	err = json.Unmarshal(body, &doc)
	return doc, err
}

func (c Client) AddDocument(document_id string, doc Document) error {
	doc_json, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = c.watsonClient.MakeRequest("PUT", c.version+document_id, bytes.NewReader(doc_json), nil)
	return err
}

func (c Client) UpdateDocument(document_id string, doc Document) error {
	doc_json, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = c.watsonClient.MakeRequest("POST", c.version+document_id, bytes.NewReader(doc_json), nil)
	return err
}

func (c Client) DeleteDocument(document_id string) error {
	_, err := c.watsonClient.MakeRequest("DELETE", c.version+document_id, nil, nil)
	return err
}

type DocumentProcessingState struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func (c Client) GetDocumentProcessingState(document_id string) (DocumentProcessingState, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+document_id+"/processing_state", nil, nil)
	if err != nil {
		return DocumentProcessingState{}, err
	}
	var state DocumentProcessingState
	err = json.Unmarshal(body, &state)
	return state, err
}

type DocumentAnnotations struct {
	Id          string         `json:"id"`
	Label       string         `json:"label,omitempty"`
	Annotations [][]Annotation `json:"annotations"`
}

func (c Client) GetDocumentAnnotations(document_id string) (DocumentAnnotations, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+document_id+"/annotations", nil, nil)
	if err != nil {
		return DocumentAnnotations{}, err
	}
	var annotations DocumentAnnotations
	err = json.Unmarshal(body, &annotations)
	return annotations, err
}

func (c Client) GetDocumentRelatedConcepts(document_id string, options map[string]interface{}) (ConceptMatches, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	body, err := c.watsonClient.MakeRequest("GET", c.version+document_id+"/related_concepts?"+q.Encode(), nil, nil)
	if err != nil {
		return ConceptMatches{}, err
	}
	var conceptMatches ConceptMatches
	err = json.Unmarshal(body, &conceptMatches)
	return conceptMatches, err

}

func (c Client) GetDocumentRelationScores(document_id string, to_concepts []string) (ConceptScores, error) {
	concepts_json, err := json.Marshal(to_concepts)
	if err != nil {
		return ConceptScores{}, err
	}
	q := url.Values{}
	q.Set("concepts", string(concepts_json))
	body, err := c.watsonClient.MakeRequest("GET", c.version+document_id+"/relation_scores?"+q.Encode(), nil, nil)
	if err != nil {
		return ConceptScores{}, err
	}
	var concepts ConceptScores
	err = json.Unmarshal(body, &concepts)
	return concepts, err
}
