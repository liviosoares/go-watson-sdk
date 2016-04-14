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

package personality_insights

import (
	"encoding/json"
	"io"
	"net/http"

	"github.ibm.com/lsoares/go-watson-sdk/watson"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultPiVersion = "v2"

// Connects to instance of Watson Natural Language Classifier service
func NewClient(cfg watson.Config) (Client, error) {
	pi := Client{version: "/" + defaultPiVersion}
	cfg.Credentials.ServiceName = "personality_insights"
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	pi.watsonClient = client
	return pi, nil
}

type Profile struct {
	// Detailed results for a specific characteristic of the input text.
	Tree TraitTree `json:"tree"`
	// The unique identifier for which these characteristics were computed, from the "userid" field of the input ContentItems.
	Id string `json:"id"`
	// The source for which these characteristics were computed, from the "sourceid" field of the input ContentItems.
	Source string `json:"source"`
	// The number of words found in the input.
	WordCount int `json:"word_count"`
	// A message indicating the number of words found and where that value falls in the range of required or suggested number of words when guidance is available.
	WordCountMessage string `json:"word_count_message,omitempty"`
	// The language model that was used to process the input, one of "en" or "es"
	ProcessedLang string `json:"processed_lang"`
}

type TraitTree struct {
	// The id of the characteristic, globally unique.
	Id string `json:"id"`
	// The user-displayable name of the characteristic.
	Name string `json:"name"`
	// The category of the characteristic: personality, needs, values, or behavior (for temporal data).
	Category string `json:"category,omitempty"`
	// For personality, needs, and values characterisitics, the normalized percentile score for the characteristic, from 0-1. For example, if
	// the percentage for Openness is 0.25, the author scored in the 25th percentile; the author is more open than 24% of the population and
	// less open than 74% of the population. For temporal behavior characteristics, the percentage of timestamped data that occurred
	// during that day or hour.
	Percentage float64 `json:"percentage,omitempty"`
	// Indicates the sampling error of the percentage based on the number of words in the input, from 0-1. The number defines a 95% confidence
	// interval around the percentage. For example, if the sampling error is 4% and the percentage is 61%, it is 95% likely that the actual
	// percentage value is between 57% and 65% if more words are given.
	SamplingError float64 `json:"sampling_error,omitempty"`
	// For personality, needs, and values characterisitics, the raw score for the characteristic. A positive or negative score indicates more
	// or less of the characteristic; zero indicates neutrality or no evidence for a score. The raw score is computed based on the input
	// and the service model; it is not normalized or compared with a sample population. The raw score enables comparison of the results against a
	// different sampling population and with a custom normalization approach.
	RawScore float64 `json:"raw_score,omitempty"`
	// Indicates the sampling error of the raw score based on the number of words in the input; the practical range is 0-1. The number defines a
	// 95% confidence interval around the raw score. For example, if the raw sampling error is 5% and the raw score is 65%, it is 95% likely that
	// the actual raw score is between 60% and 70% if more words are given.
	RawSamplingError float64 `json:"raw_sampling_error,omitempty"`
	// Recursive array of characteristics inferred from the input text
	Children []TraitTree `json:"children,omitempty"`
}

// Calls /v2/profile to send data for analysis (up to 20MB)

// The content type of the request can be: plain text (the default), HTML, or JSON. When specifying a content type of plain text or
// HTML, include the charset parameter to indicate the character encoding of the input text, for example, "text/plain; charset=utf-8".
func (c Client) GetProfile(data io.Reader, content_type string, language string) (Profile, error) {
	headers := make(http.Header)
	headers.Set("Content-Type", content_type)
	if len(language) > 0 {
		headers.Set("Content-Language", language)
	}
	b, err := c.watsonClient.MakeRequest("POST", c.version+"/profile", data, headers)
	if err != nil {
		return Profile{}, err
	}
	var profile Profile
	err = json.Unmarshal(b, &profile)
	return profile, err
}
