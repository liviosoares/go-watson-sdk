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

// Package speech_to_text provides an interface to Watson Speech to Text service.
package speech_to_text

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"

	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/authorization"
	"golang.org/x/net/websocket"
)

type Client struct {
	version      string
	watsonClient *watson.Client
}

const defaultMajorVersion = "v1"
const defaultUrl = "https://gateway.watsonplatform.net/speech-to-text/api"

// Connects to instance of Watson Natural Language Classifier service
func NewClient(cfg watson.Config) (Client, error) {
	tts := Client{version: "/" + defaultMajorVersion}
	if len(cfg.Credentials.ServiceName) == 0 {
		cfg.Credentials.ServiceName = "speech_to_text"
	}
	if len(cfg.Credentials.Url) == 0 {
		cfg.Credentials.Url = defaultUrl
	}
	client, err := watson.NewClient(cfg.Credentials)
	if err != nil {
		return Client{}, err
	}
	tts.watsonClient = client
	return tts, nil
}

type ModelList struct {
	Models []Model `json:"models"`
}

type Model struct {
	// Name of the model for use as an identifier in calls to the service (for example, en-US_BroadbandModel). ,
	Name string `json:"name"`
	// Language identifier for the model (for example, en-US). ,
	Language string `json:"language"`
	// Sampling rate (minimum acceptable rate for audio) used by the model in Hertz. ,
	Rate int `json:"rate"`
	// URI for the model.
	URL string `json:"url"`
	// URI for the model for use with the POST /v1/sessions method. ,
	Sessions string `json:"session,omitempty"`
	// Brief description of the model.
	Description string `json:"description"`
}

// Calls 'GET /v1/models' to retrieves the models available for the servic
func (c Client) ListModels() (ModelList, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/models", nil, nil)
	if err != nil {
		return ModelList{}, err
	}
	var models ModelList
	err = json.Unmarshal(body, &models)
	return models, err
}

// Calls 'GET /v1/models/{model_id}' to retrieve information about the model
func (c Client) GetModel(model_id string) (Model, error) {
	body, err := c.watsonClient.MakeRequest("GET", c.version+"/models/"+model_id, nil, nil)
	if err != nil {
		return Model{}, err
	}
	var model Model
	err = json.Unmarshal(body, &model)
	return model, err
}

type Event struct {
	// The results array consists of 0 or more final results, followed by 0 or 1 interim result. The final results are guaranteed not to change,
	// while the interim result may be replaced by 0 or more final results, followed by 0 or 1 interim result. The service periodically sends
	// "updates" to the result list, with the result_index set to the lowest index in the array that has changed. result_index always points to
	// the slot just after the most recent final result. ,
	Results []Result `json:"results,omitempty"`
	//  Index indicating change point in the results array. See description of the results array for further information. ,
	ResultIndex int `json:"result_index,omitempty"`
	// Array of warning messages about invalid query parameters or JSON fields included with the request. Each element of the array includes a string that describes the nature of the warning followed by an array of invalid argument strings; for example, "Unknown arguments: [u'invalid_arg_1', u'invalid_arg_2']." The request succeeds despite the warnings.
	Warnings []string `json:"warnings,omitempty"`
	Error    string   `json:"error,omitempty"`
}

type Result struct {
	// If true, the result for this utterance is not updated further.
	Final bool `json:"final"`
	// List of alternative transcripts received from the service.
	Alternatives []Alternative `json:"alternatives,omitempty"`
	// Dictionary (or associative array) whose keys are the strings specified for keywords if both that parameter and keywords_threshold
	// are specified. The array is empty for any keyword for which no matches are found.
	KeywordsResult KeywordResults `json:"keyword_results,omitempty"`
	// List of word alternative hypotheses found for words of the input audio if word_alternatives_threshold is not null.
	WordAlternatives []WordAlternativeResults `json:"word_alternatives,omitempty"`
}

type Alternative struct {
	// Transcript of the utterance.
	Transcript string `json:"transcript,omitempty"`
	// Confidence score of the transcript, between 0 and 1. Available only for the best alternative and only in results marked as final.
	Confidence float64 `json:"confidence,omitempty"`
	// Time alignments for each word from transcript as list of lists. Each/ inner list consists of 3 elements: the word, start and end of the
	// word in seconds. Example: [["hello",0.0,1.2],["world",1.2,2.5]]. Available only for the best alternative.
	Timestamps [][]interface{} `json:"timestamps,omitempty"`
	// Confidence score for each word of the transcript, between 0 and 1. Each inner list consists of 2 elements: the word and the
	// confidence of the word. Example: [["hello",0.95],["world",0.866]]. Available only for the best alternative and only in results marked as final.
	WordConfidence []string `json:"word_confidence,omitempty"`
}

type KeywordResults struct {
	// List of each keyword entered via the keywords parameter and an array of its occurrences in the input audio. The keys of the list are the
	// actual keyword strings. The array that is returned for any keyword is empty if no occurrences of the keyword are spotted in the input.
	Keyword []KeywordResult `json:"keyword"`
}

type KeywordResult struct {
	// Specified keyword normalized to the spoken phrase that matched in the audio input. ,
	NormalizedText string `json:"normalized_text"`
	// Start time in seconds of the keyword match.
	StartTime int64 `json:"start_time"`
	// End time in seconds of the keyword match.
	EndTime int64 `json:"end_time"`
	// Confidence score of the keyword match.
	Confidence float64 `json:"confidence"`
}

type WordAlternativeResults struct {
	// Start time in seconds of the word that corresponds to the word alternative.
	StartTime int `json:"start_time"`
	// End time in seconds of the word that corresponds to the word alternative.
	EndTime int `json:"end_time"`
	// List of word alternative hypotheses for a word from the input audio.
	Alternatives []WordAlternativeResult
}

type WordAlternativeResult struct {
	// Confidence score of the word alternative hypothesis.
	Confidence float64 `json:"confidence"`
	// Word alternative hypothesis for a word from the input audio.
	Word string `json:"word"`
}

type StateReply struct {
	State string `json:"state"`
}

// NewStream creates a stream into the speech-to-text API. It returns an output channel of Events, and an io.WriteCloser for writing
// audio data. Upon the generation of a transcription event, an 'Event' object is pushed into the output channel.
// 'options' can contain options to be send in the stream initialization; more information available at:
// http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/doc/speech-to-text/websockets.shtml#WSstart
func (c Client) NewStream(model string, customization_id string, content_type string, options map[string]interface{}) (<-chan Event, io.WriteCloser, error) {
	token, err := authorization.GetToken(c.watsonClient.Creds)
	if err != nil {
		return nil, nil, errors.New("failed to acquire auth token: " + err.Error())
	}
	u, err := url.Parse(c.watsonClient.Creds.Url)
	if err != nil {
		return nil, nil, err
	}
	u.Scheme = "wss"
	q := url.Values{}
	q.Set("watson-token", token)
	if len(model) > 0 {
		q.Set("model", model)
	}
	if len(customization_id) > 0 {
		q.Set("customization_id", customization_id)
	}
	u.RawQuery = q.Encode()
	u.Path += c.version + "/recognize"

	origin, err := url.Parse(c.watsonClient.Creds.Url)
	if err != nil {
		return nil, nil, err
	}

	config := &websocket.Config{
		Location: u,
		Origin:   origin,
		Version:  websocket.ProtocolVersionHybi13,
	}
	ws, err := websocket.DialConfig(config)
	if err != nil {
		return nil, nil, errors.New("error dialing websocket: " + err.Error())
	}
	output := make(chan Event, 100)
	s := stream{
		input:       output,
		contentType: content_type,
		ws:          ws,
		options:     options,
	}
	return output, &s, nil
}

type stream struct {
	ws             *websocket.Conn
	contentType    string
	options        map[string]interface{}
	interimResults bool
	input          chan<- Event
	started        bool
	stopped        bool
}

func (s *stream) Write(p []byte) (n int, err error) {
	if s.stopped {
		return 0, errors.New("cannot write to stopped stream")
	}
	if !s.started {
		m := make(map[string]interface{})
		for k, v := range s.options {
			m[k] = v
		}
		m["action"] = "start"
		m["content-type"] = s.contentType
		// m["inactivity_timeout"] = -1
		if interim, present := m["interim_results"]; present {
			if interimResults, ok := interim.(bool); ok {
				s.interimResults = interimResults
			}
		}
		m["interim_results"] = true
		err := websocket.JSON.Send(s.ws, m)
		if err != nil {
			return 0, err
		}
		s.started = true
		var state StateReply
		err = websocket.JSON.Receive(s.ws, &state)
		if err != nil {
			return 0, err
		}
		go s.readReplies()
	}
	err = websocket.Message.Send(s.ws, p)
	return len(p), err

}

func (s *stream) Close() error {
	if s.stopped {
		return errors.New("stream already stopped")
	}
	s.stopped = true
	m := map[string]interface{}{"action": "stop"}
	return websocket.JSON.Send(s.ws, m)
}

func (s *stream) readReplies() {
	for {
		// read generic JSON
		var b []byte
		err := websocket.Message.Receive(s.ws, &b)
		if err != nil {
			close(s.input)
			return
		}
		// try to unmarshal it as a `state` reply
		var state StateReply
		err = json.Unmarshal(b, &state)
		if err == nil && len(state.State) > 0 {
			if s.stopped {
				close(s.input)
				return
			}
			continue
		}
		// next, try to unmarshal it as an `Event`
		var event Event
		err = json.Unmarshal(b, &event)
		if err != nil {
			continue
		}
		if s.interimResults == false && len(event.Results) > 0 && event.Results[0].Final == false {
			continue
		}
		s.input <- event
	}
}
