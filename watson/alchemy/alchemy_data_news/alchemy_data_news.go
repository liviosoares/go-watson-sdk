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

// Package alchemy_data_news provides an interface to AlchemyData News API.
// More documentation available at: doc.alchemyapi.com
package alchemy_data_news

import (
	"fmt"

	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/alchemy"
)

type Client struct {
	alchemyClient *alchemy.Client
}

// Connects to instance of Watson AlchemyData News
func NewClient(cfg watson.Config) (Client, error) {
	client, err := alchemy.NewClient(cfg)
	if err != nil {
		return Client{}, err
	}
	return Client{alchemyClient: &client}, nil
}

type Result map[string]interface{}

// GetNews calls the AlchemyData News endpoint to retrieve recent news articles according to the provided query.
// start and end is used to specify the time range desired for news articles. The exact format of the time string
// is documented at: https://alchemyapi.readme.io/v1.0/docs/rest-api-documentation
// The query parameter supports the AlchemyData News query language also documented here:
// https://alchemyapi.readme.io/v1.0/docs/rest-api-documentation#section-parameters-filters-
//
// One simple example of usage:
//
// 	result, err := c.GetNews("now-24h", "now",
// 					map[string]interface{}{
// 							"q.enriched.url.title": "[baseball^soccer]",
// 							"return": "enriched.url.title,enriched.url.author,original.url",
// 					})
//
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
