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

// Package text_to_speech provides an interface to Watson Text to Speech service.
package text_to_speech

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

const defaultTextToSpeechVersion = "v1"

// Connects to instance of Watson Natural Language Classifier service
func NewClient(cfg watson.Config) (Client, error) {
	tts := Client{version: "/" + defaultTextToSpeechVersion}
	cfg.Credentials.ServiceName = "text_to_speech"
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	tts.watsonClient = client
	return tts, nil
}

type VoiceList struct {
	Voices []Voice `json:"voices"`
}

type Voice struct {
	// URI of the voice
	URL string `json:"url"`
	// Gender of the voice: 'male' or 'female'. ,
	Gender string `json:"gender"`
	// Name of the voice. Use this as the voice identifier in all requests. ,
	Name string `json:"name"`
	// Language and region of the voice (for example, 'en-US'). ,
	Language string `json:"language"`
	// Textual description of the voice. ,
	Description string `json:"description"`
	// If true, the voice can be customized; if false, the voice cannot be customized.
	Customizable bool `json:"customizable"`
	// Information about a specific custom voice model. Returned only when the customization_id parameter is specified with the method.
	Customization *Customization `json:"customization,omitempty"`
}

type Customization struct {
	// GUID of the custom voice model. ,
	CustomizationId string `json:"customization_id"`
	// Name of the custom voice model. ,
	Name string `json:"name"`
	// Language of the custom voice model. = ['en-US', 'en-GB', 'es-ES', 'es-US', 'fr-FR', 'it-IT', 'ja-JP', 'pt-BR'],
	Language string `json:"language"`
	// GUID that associates the owning user with the custom voice model. ,
	Owner string `json:"owner"`
	// UNIX timestamp that indicates when the custom voice model was created. The timestamp is a count of seconds since the UNIX Epoch of January 1, 1970 Coordinated Universal Time (UTC). ,
	Created int `json:"created"`
	// UNIX timestamp that indicates when the custom voice model was last modified. Equals created when the new voice model is first added but has yet to be changed. ,
	LastModified int `json:"last_modified"`
	// Description of the custom voice model.
	Description string `json:"description"`
}

// Calls 'GET /v1/voices' to retrieve all voices available for speech synthesis
func (c Client) ListVoices() (VoiceList, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/voices", nil, nil)
	if err != nil {
		return VoiceList{}, err
	}
	var voices VoiceList
	err = json.Unmarshal(body, &voices)
	return voices, err
}

// Calls 'GET /v1/voices/{voice}' to retrieve a specific voice available for speech synthesis
func (c Client) GetVoice(voice_id string, customization_id string) (Voice, error) {
	q := url.Values{}
	if len(customization_id) > 0 {
		q.Set("customization_id", customization_id)
	}
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/voices/"+voice_id+"?"+q.Encode(), nil, nil)
	if err != nil {
		return Voice{}, err
	}
	var voice Voice
	err = json.Unmarshal(body, &voice)
	return voice, err
}

// Valid accept values are: "audio/ogg; codecs=opus", "audio/wav", "audio/flac"
func (c Client) Synthesize(text string, voice string, accept string, customization_id string) ([]byte, error) {
	t := struct {
		Text string `json:"text"`
	}{Text: text}
	text_json, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if len(customization_id) > 0 {
		q.Set("customization_id", customization_id)
	}
	if len(voice) > 0 {
		q.Set("voice", voice)
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", accept)
	return c.watsonClient.MakeRequest("POST", c.version+"/synthesize?"+q.Encode(), bytes.NewReader(text_json), headers)
}

// Calls 'GET /v1/pronunciation' to get the pronunciation for a word
// format can be either 'ipa' (default) or 'spr'
func (c Client) GetPronunciation(text string, voice string, format string) (string, error) {
	q := url.Values{}
	if len(voice) > 0 {
		q.Set("voice", voice)
	}
	if len(format) > 0 {
		q.Set("format", format)
	}
	q.Set("text", text)
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/pronunciation?"+q.Encode(), nil, nil)
	if err != nil {
		return "", err
	}
	var p struct {
		Pronunciation string `json:"pronunciation"`
	}
	err = json.Unmarshal(body, &p)
	return p.Pronunciation, err
}
