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

// Package alchemy_language provides an interface to Alchemy Language services.
// More documentation available at: http://www.alchemyapi.com/api
package alchemy_language

import (
	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/alchemy"
)

type Client struct {
	alchemyClient *alchemy.Client
}

// Connects to instance of Watson Alchemy Language services
func NewClient(cfg watson.Config) (Client, error) {
	client, err := alchemy.NewClient(cfg)
	if err != nil {
		return Client{}, err
	}
	return Client{alchemyClient: &client}, nil
}

type SentimentResponse struct {
	alchemy.BaseResponse
	DocSentiment DocSentiment `json:"docSentiment"`
}
type DocSentiment struct {
	Type  string  `json:"type"`
	Score float64 `json:"score,string"`
	Mixed int     `json:"mixed,string"`
}

// Calls '*GetTextSentiment'. See documentation at http://www.alchemyapi.com/api/sentiment-analysis
func (c Client) GetSentiment(data []byte, options map[string]interface{}) (SentimentResponse, error) {
	var response SentimentResponse
	err := c.alchemyClient.Call("GetTextSentiment", data, options, &response)
	return response, err
}

type TargetedSentimentResponse struct {
	alchemy.BaseResponse
	Results []ResultSentiment `json:"results"`
}

type ResultSentiment struct {
	Sentiment DocSentiment `json:"sentiment"`
	Text      string       `json:"text"`
}

// Calls '*GetTargetedSentiment'
func (c Client) GetSentimentTargeted(data []byte, targets []string, options map[string]interface{}) (TargetedSentimentResponse, error) {
	if options == nil {
		options = make(map[string]interface{})
	}
	targets_encoded := ""
	for i := range targets {
		targets_encoded += targets[i]
		if i < len(targets)-1 {
			targets_encoded += "|"
		}
	}
	options["targets"] = targets_encoded
	var response TargetedSentimentResponse
	err := c.alchemyClient.Call("GetTargetedSentiment", data, options, &response)
	return response, err
}

type EmotionResponse struct {
	alchemy.BaseResponse
	DocEmotions struct {
		Anger   float64 `json:"anger,string,omitempty"`
		Disgust float64 `json:"disgust,string,omitempty"`
		Fear    float64 `json:"fear,string,omitempty"`
		Joy     float64 `json:"joy,string,omitempty"`
		Sadness float64 `json:"sadness,string,omitempty"`
	} `json:"docEmotions,omitempty"`
}

func (c Client) GetEmotion(data []byte, options map[string]interface{}) (EmotionResponse, error) {
	var response EmotionResponse
	err := c.alchemyClient.Call("GetEmotion", data, options, &response)
	return response, err
}

type TaxonomyResponse struct {
	alchemy.BaseResponse
	Taxonomy []Taxonomy `json:"taxonomy,omitempty"`
}

type Taxonomy struct {
	Label     string  `json:"label,omitempty"`
	Score     float64 `json:"score,string,omitempty"`
	Confident int     `json:"confident,string,omitempty"`
}

func (c Client) GetTaxonomy(data []byte, options map[string]interface{}) (TaxonomyResponse, error) {
	var response TaxonomyResponse
	err := c.alchemyClient.Call("GetRankedTaxonomy", data, options, &response)
	return response, err
}

type ConceptsResponse struct {
	alchemy.BaseResponse
	Concepts []Concept `json:"concepts,omitempty"`
}

type Concept struct {
	Text           string         `json:"text,omitempty"`
	Relevance      float64        `json:"relevance,string,omitempty"`
	KnowledgeGraph KnowledgeGraph `json:"knowledgeGraph,omitempty"`
	Website        string         `json:"website,omitempty"`
	Geo            string         `json:"geo,omitempty"`
	DBpedia        string         `json:"dbpedia,omitempty"`
	Yago           string         `json:"yago,omitempty"`
	OpenCyc        string         `json:"opencyc,omitempty"`
	Freebase       string         `json:"freebase,omitempty"`
	CiaFactbook    string         `json:"ciaFactbook,omitempty"`
	Census         string         `json:"census,omitempty"`
	GeoNames       string         `json:"geonames,omitempty"`
	MusicBrainz    string         `json:"musicBrainz,omitempty"`
	CrunchBase     string         `json:"crunchbase,omitempty"`
}

type KnowledgeGraph struct {
	TypeHierarchy string `json:"typeHierarchy"`
}

func (c Client) GetConcepts(data []byte, options map[string]interface{}) (ConceptsResponse, error) {
	var response ConceptsResponse
	err := c.alchemyClient.Call("GetRankedConcepts", data, options, &response)
	return response, err
}

type NamedEntitiesResults struct {
	alchemy.BaseResponse
	Entities []Entity `json:"entities,omitempty"`
}

type Entity struct {
	Type           string                `json:"type,omitempty"`
	Relevance      float64               `json:"relevance,string,omitempty"`
	KnowledgeGraph KnowledgeGraph        `json:"knowledgeGraph,omitempty"`
	Count          int                   `json:"count,string,omitempty"`
	Text           string                `json:"text,omitempty"`
	Disambiguated  alchemy.Disambiguated `json:"disambiguated,omitempty"`
	Quotations     []Quotation           `json:"quotations,omitempty"`
	Sentiment      DocSentiment          `json:"sentiment,omitempty"`
}

type Quotation struct {
	Quotation string `json:"quotation"`
}

func (c Client) GetNamedEntities(data []byte, options map[string]interface{}) (NamedEntitiesResults, error) {
	var response NamedEntitiesResults
	err := c.alchemyClient.Call("GetRankedNamedEntities", data, options, &response)
	return response, err
}

type KeywordsResults struct {
	alchemy.BaseResponse
	Keywords []Keyword `json:"keywords,omitempty"`
}

type Keyword struct {
	Text           string         `json:"text,omitempty"`
	Relevance      float64        `json:"relevance,string,omitempty"`
	KnowledgeGraph KnowledgeGraph `json:"knowledgeGraph,omitempty"`
	Sentiment      DocSentiment   `json:"sentiment,omitempty"`
}

func (c Client) GetKeywords(data []byte, options map[string]interface{}) (KeywordsResults, error) {
	var response KeywordsResults
	err := c.alchemyClient.Call("GetRankedKeywords", data, options, &response)
	return response, err
}

type RelationsResults struct {
	alchemy.BaseResponse
	Relations []Relation `json:"relations,omitempty"`
}

type Relation struct {
	Sentence string          `json:"sentence,omitempty"`
	Subject  RelationSubject `json:"subject,omitempty"`
	Action   RelationAction  `json:"action,omitempty"`
	Object   RelationObject  `json:"object,omitempty"`
}

type RelationSubject struct {
	Text      string       `json:"text,omitempty"`
	Sentiment DocSentiment `json:"sentiment,omitempty"`
	Entity    Entity       `json:"entity,omitempty"`
}

type RelationAction struct {
	Text       string       `json:"text,omitempty"`
	Lemmatized string       `json:"lemmatized,omitempty"`
	Verb       RelationVerb `json:"verb,omitempty"`
}

type RelationVerb struct {
	Text    string `json:"text,omitempty"`
	Tense   string `json:"tense,omitempty"`
	Negated string `json:"negated,omitempty"`
}

type RelationObject struct {
	Text                 string       `json:"text,omitempty"`
	Sentiment            DocSentiment `json:"sentiment,omitempty"`
	SentimentFromSubject DocSentiment `json:"sentimentFromSubject,omitempty"`
	Entity               Entity       `json:"entity,omitempty"`
}

func (c Client) GetRelations(data []byte, options map[string]interface{}) (RelationsResults, error) {
	var response RelationsResults
	err := c.alchemyClient.Call("GetRelations", data, options, &response)
	return response, err
}

func (c Client) GetText(data []byte, options map[string]interface{}) (alchemy.BaseResponse, error) {
	var response alchemy.BaseResponse
	err := c.alchemyClient.Call("GetText", data, options, &response)
	return response, err
}

func (c Client) GetRawText(data []byte, options map[string]interface{}) (alchemy.BaseResponse, error) {
	var response alchemy.BaseResponse
	err := c.alchemyClient.Call("GetRawText", data, options, &response)
	return response, err
}

func (c Client) GetTitle(data []byte, options map[string]interface{}) (alchemy.BaseResponse, error) {
	var response alchemy.BaseResponse
	err := c.alchemyClient.Call("GetTitle", data, options, &response)
	return response, err
}

type AuthorResponse struct {
	Status string `json:"status"`
	Url    string `json:"url,omitempty"`
	Author string `json:"author,omitempty"`
}

func (c Client) GetAuthor(data []byte, options map[string]interface{}) (AuthorResponse, error) {
	var response AuthorResponse
	err := c.alchemyClient.Call("GetAuthor", data, options, &response)
	return response, err
}

type AuthorsResponse struct {
	Status  string `json:"status"`
	Url     string `json:"url,omitempty"`
	Authors struct {
		Names []string `json:"names,omitempty"`
	} `json:"authors,omitempty"`
}

func (c Client) GetAuthors(data []byte, options map[string]interface{}) (AuthorsResponse, error) {
	var response AuthorsResponse
	err := c.alchemyClient.Call("GetAuthors", data, options, &response)
	return response, err
}

type LanguageResponse struct {
	Status         string `json:"status"`
	Usage          string `json:"usage,omitempty"`
	Url            string `json:"url,omitempty"`
	Language       string `json:"language,omitempty"`
	ISO6391        string `json:"iso-639-1,omitempty"`
	ISO6392        string `json:"iso-639-2,omitempty"`
	ISO6393        string `json:"iso-639-3,omitempty"`
	Ethnologue     string `json:"ethnologue,omitempty"`
	NativeSpeakers string `json:"native-speakers,omitempty"`
	Wikipedia      string `json:"wikipedia,omitempty"`
}

func (c Client) GetLanguage(data []byte, options map[string]interface{}) (LanguageResponse, error) {
	var response LanguageResponse
	err := c.alchemyClient.Call("GetLanguage", data, options, &response)
	return response, err
}

type FeedLinksResponse struct {
	Status string `json:"status"`
	Url    string `json:"url,omitempty"`
	Feeds  []struct {
		Feed string `json:"feed"`
	} `json:"feeds,omitempty"`
}

func (c Client) GetFeedLinks(data []byte, options map[string]interface{}) (FeedLinksResponse, error) {
	var response FeedLinksResponse
	err := c.alchemyClient.Call("GetFeedLinks", data, options, &response)
	return response, err
}

type DatesResponse struct {
	alchemy.BaseResponse
	Dates []struct {
		Date string `json:"date,omitempty"`
		Text string `json:"text,omitempty"`
	} `json:"dates,omitempty"`
}

func (c Client) ExtractDates(data []byte, options map[string]interface{}) (DatesResponse, error) {
	var response DatesResponse
	err := c.alchemyClient.Call("ExtractDates", data, options, &response)
	return response, err
}

type PubDatesResponse struct {
	alchemy.BaseResponse
	PublicationDate struct {
		Date      string `json:"date,omitempty"`
		Confident string `json:"confident,omitempty"`
	} `json:"publicationDate,omitempty"`
}

func (c Client) GetPubDate(data []byte, options map[string]interface{}) (PubDatesResponse, error) {
	var response PubDatesResponse
	err := c.alchemyClient.Call("GetPubDate", data, options, &response)
	return response, err
}
