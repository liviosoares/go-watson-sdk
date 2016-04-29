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

package alchemy_data_news

import (
	"fmt"

	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/alchemy"
)

type Client struct {
	alchemyClient *alchemy.Client
}

// Connects to instance of Watson Concept Insights service
func NewClient(cfg watson.Config) (Client, error) {
	client, err := alchemy.NewClient(cfg)
	if err != nil {
		return Client{}, err
	}
	return Client{alchemyClient: &client}, nil
}

type Result map[string]interface{}
// 	alchemy.BaseResponse
// 	Result struct {
// 		Count  int    `json:"count,omitempty"`
// 		Status string `json:"status,omitempty"`
// 		Docs   []struct {
// 			Id     string `json:"id,omitempty"`
// 			Source struct {
// 				Original struct {
// 					Url string `json:"url,omitempty"`
// 				} `json:"original,omitempty"`
// 				Enriched struct {
// 					Url struct {
// 						Author        string       `json:"author,omitempty"`
// 						CleanedTitle  string       `json:cleanedTitle,omitempty"`
// 						Concepts      []Concept    `json:"concepts,omitempty"`
// 						DocSentiment  DocSentiment `json:"docSentiment,omitempty"`
// 						EnrichedTitle struct {
// 							Concepts     []Concept    `json:"concepts,omitempty"`
// 							DocSentiment DocSentiment `json:"docSentiment,omitempty"`
// 							Entities []string
// 						} `json:"enrichedTitle,omitempty"`
// 					} `json:"url,omitempty"`
// 				}
// 			} `json:"source,omitempty"`
// 			Timestamp int `json:"timestamp,omitempty"`
// 		} `json:"docs,omitempty"`
// 		Next string `json:"next,omitempty"`
// 	} `json:"result,omitempty"`
// }

// type Concept struct {
// 	KnowledgeGraph struct {
// 		TypeHierarchy string `json:"typeHierarchy"`
// 	} `json:"knowledgeGraph,omitempty"`
// 	Relevance float64 `json:relevance,omitempty"`
// 	Text      string  `json:"text,omitempty"`
// }

// type DocSentiment struct {
// 	Mixed int     `json:"mixed,omitempty"`
// 	Score float64 `json:"score,omitempty"`
// 	Type  string  `json:"type,omitempty"`
// }

func (c Client) GetNews(start string, end string, query map[string]interface{}) (Result, error) {
	q := make(map[string]interface{})
	for k, v := range query {
		q[k] = fmt.Sprintf("%v", v)
	}
	q["start"] = start
	q["end"] = end

	var result Result
	err := c.alchemyClient.Get("/data/GetNews", q, &result)
	return result, err
}
