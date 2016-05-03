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

package alchemy_vision

import (
	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/alchemy"
)

type Client struct {
	alchemyClient *alchemy.Client
}

// Connects to instance of Watson Alchmey Vision services
func NewClient(cfg watson.Config) (Client, error) {
	client, err := alchemy.NewClient(cfg)
	if err != nil {
		return Client{}, err
	}
	return Client{alchemyClient: &client}, nil
}

type ImageKeywordsResponse struct {
	alchemy.BaseResponse
	ImageKeywords []struct {
		Text           string  `json:"text,omitempty"`
		Score          float64 `json:"score,string,omitempty"`
		KnowledgeGraph struct {
			TypeHierarchy string `json:"typeHierarchy"`
		} `json:"knowledgeGraph,omitempty"`
	} `json:"imageKeywords,omitempty"`
}

// Calls '*GetRankedImageKeywords'. See documentation at http://www.alchemyapi.com/api/image-tagging
func (c Client) GetImageKeywords(data []byte, options map[string]interface{}) (ImageKeywordsResponse, error) {
	var response ImageKeywordsResponse
	err := c.alchemyClient.Call("GetRankedImageKeywords", data, options, &response)
	return response, err
}

type ImageLinkResponse struct {
	alchemy.BaseResponse
	Image string `json:"image,omitempty"`
}

// Calls '*GetImage'. See documentation at http://www.alchemyapi.com/api/image-link-extraction
func (c Client) GetImageLink(data []byte, options map[string]interface{}) (ImageLinkResponse, error) {
	var response ImageLinkResponse
	err := c.alchemyClient.Call("GetImage", data, options, &response)
	return response, err
}

type ImageFaceTagsResponse struct {
	alchemy.BaseResponse
	ImageFaces []struct {
		PositionX int `json:"positionX,string,omitempty"`
		PositionY int `json:"positionY,string,omitempty"`
		Width     int `json:"width,string,omitempty"`
		Height    int `json:"height,string,omitempty"`
		Gender    struct {
			Gender string  `json:"gender,omitempty"`
			Score  float64 `json:"score,string,omitempty"`
		} `json:"gender,omitempty"`
		Age struct {
			AgeRange string  `json:"ageRange,omitempty"`
			Score    float64 `json:"score,string,omitempty"`
		} `json:"age,omitempty"`
		Identify struct {
			Name           string  `json:"name,omitempty"`
			Score          float64 `json:"score,string,omitempty"`
			KnowledgeGraph struct {
				TypeHierarchy string `json:"typeHierarchy"`
			} `json:"knowledgeGraph,omitempty"`
		} `json:"identity,omitempty"`
		Disambiguated alchemy.Disambiguated `json:"disambiguated,omitempty"`
	} `json:"imageFaces,omitempty"`
}

// Calls '*GetRankedImageFaceTags'. See documentation at http://www.alchemyapi.com/api/face-detection
func (c Client) GetImageFaceTags(data []byte, options map[string]interface{}) (ImageFaceTagsResponse, error) {
	var response ImageFaceTagsResponse
	err := c.alchemyClient.Call("GetRankedImageFaceTags", data, options, &response)
	return response, err
}
