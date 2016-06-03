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

// Package tone_analyzer provides an interface to Watson Tone Analyzer service.
package tone_analyzer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/liviosoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v3"
const defaultMinorVersion = "2016-02-11"
const defaultUrl = "https://gateway.watsonplatform.net/tone-analyzer-beta/api"

// Connects to instance of Watson Natural Language Classifier service
func NewClient(cfg watson.Config) (Client, error) {
	ta := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "tone_analyzer"
	}
	if len(cfg.Credentials.Url) == 0 {
		cfg.Credentials.Url = defaultUrl
	}
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	ta.watsonClient = client
	return ta, nil
}

type Analysis struct {
	// Tone analysis results performed on the entire document's text. This includes three tone categories: Social Tone, Emotion Tone and Writing Tone. ,
	DocumentTone DocumentAnalysis `json:"document_tone"`
	// List of sentences contained in the document, with individual Tone analysis results for each sentence.
	SentencesTone []SentenceAnalysis `json:"sentences_tone"`
}

type DocumentAnalysis struct {
	ToneCategories []ToneCategory `json:"tone_categories"`
}

type SentenceAnalysis struct {
	// A unique number identifying this sentence within this document. Reserved for future use (when sentences need to be referred from different places). ,
	SentenceId int `json:"sentence_id"`
	// Index of the character in the document where this sentence starts. ,
	InputFrom int `json:"input_from"`
	// Index of the character in the document after the end of this sentence (input_to minus input_from is the length of this sentence in characters). ,
	InputTo int `json:"input_to"`
	// The text in this sentence - as just taken from the input text from input_from to input_to. ,
	Text string `json:"text"`
	// Tone analysis results for this sentence; divided in three Tone categories: Social Tone, Emotion Tone and Writing Tone.
	ToneCategories []ToneCategory `json:"tone_categories"`
}

type ToneCategory struct {
	// Name of this tone category: one of Emotion, Social or Writing Tone. Human-readable, localized. ,
	CategoryName string `json:"category_name"`
	// Identifier of this category. It does not vary across languages or localizations. ,
	CategoryId string `json:"category_id"`
	// All individual tone results within this category. For example, the Social Tones category contains one element for each of the dimensions in Big 5 model: Agreeableness, Openness, etc.
	Tones []ToneScore `json:"tones"`
}

type ToneScore struct {
	// The name of the tone. Human-readable, localized. ,
	ToneName string `json:"tone_name"`
	// Identifier of this tone. It does not vary across languages and localizations. ,
	ToneId string `json:"tone_id"`
	// Name of the category that this tone belongs to: one of Emotion, Social or Writing Tone. Human-readable, localized. ,
	ToneCategoryName string `json:"tone_category_name"`
	// Identifier of the category that this tone belongs to. It does not vary across languages or localizations. ,
	ToneCategoryId string `json:"tone_category_id"`
	// A raw score computed by the algorithms. This can be compared to other raw scores and used to build your own normalizations.
	Score float64 `json:"score"`
}

// Calls 'POST /v3/tone' to analyze the tone of a piece of text. The message is analyzed for several tones - social, emotional, and writing. For each tone,
// various traits are derived. For example, conscientiousness, agreeableness, and openness.
func (c Client) Tone(text string, options map[string]interface{}) (Analysis, error) {
	q := url.Values{}
	for k, v := range options {
		q.Set(k, fmt.Sprintf("%v", v))
	}
	q.Set("version", defaultMinorVersion)

	headers := make(http.Header)
	headers.Set("Content-Type", "text/plain")
	headers.Set("Accept", "application/json")

	body, err := c.watsonClient.MakeRequest("POST", c.version+"/tone?"+q.Encode(), strings.NewReader(text), headers)
	if err != nil {
		return Analysis{}, err
	}
	var analysis Analysis
	err = json.Unmarshal(body, &analysis)
	return analysis, err
}
