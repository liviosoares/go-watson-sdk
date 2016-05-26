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

package watson

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"
)

const goSdkVersion = "0.1.0"

// Client contains credential information for making requests to the Watson Developer Cloud gateway.
// It is used by service-specific clients to make requests and unmarshal replies and error codes.
type Client struct {
	Creds Credentials
}

// Config contains versioning and credential information to a specific Watson service.
type Config struct {
	// Version of API to use; defaults to "v1"
	Version string
	// Credentials to use. If empty, VCAP_SERVICES environment variable will be used
	Credentials Credentials
}

// Credentials contains information necessary to connect to a specific Watson service.
type Credentials struct {
	// Users can either set Url, Username and Password; or leave these blank
	// so that can get auto-filled through VCAP_SERVICES environment variable
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	// ApiKey is currently only used for Alchemy based services
	ApiKey string `json:"apikey"`
	// ServiceName needs to be filled in order to extract credentials from
	// VCAP_SERVICES environment variable
	ServiceName string
	ServicePlan string
}

// NewClient creates a generic Watson client object, using creds as Credential information.
// If creds is a Credential object with empty Url, Username or Password, then $VCAP_SERVICES
// environment variable is used to retrieve service specific credentials.
// The returned client can be used by multiple go routines concurrently (thread safe).
func NewClient(creds Credentials) (*Client, error) {
	if len(creds.Url) == 0 || len(creds.Username) == 0 || len(creds.Password) == 0 {
		var err error
		creds, err = getBluemixCredentials(creds.ServiceName, creds.ServicePlan)
		if err != nil {
			return nil, err
		}
	}
	return &Client{Creds: creds}, nil
}

// WatsonError stores error information from a Watson API endpoint
type WatsonError struct {
	Code    int
	Message string
}

func (e *WatsonError) Error() string {
	return "code: " + strconv.Itoa(e.Code) + "; message: " + e.Message
}

type watsonError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e *watsonError) Error() string {
	return "code: " + strconv.Itoa(e.Code) + "; message: " + e.Message
}

type alternativeError struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
}

func (ae *alternativeError) Error() string {
	return "code: " + strconv.Itoa(ae.Code) + "; message: " + ae.Message
}

type alternativeError1 struct {
	Code    int     `json:"code"`
	Message *string `json:"msg,omitempty"`
}

func (ae *alternativeError1) Error() string {
	s := "code: " + strconv.Itoa(ae.Code) + "; message: "
	if ae.Message != nil {
		s += *ae.Message
	}
	return s
}

// MakeRequest issues an HTTP request to one of the Watson API endpoints. Authentication information from the Client
// object is used.
// If the endpoint replies with a non-20x reply, an error of WatsonError type is returned, otherwise
// the body of the reply is returned.
func (c *Client) MakeRequest(method string, path string, body io.Reader, header http.Header) ([]byte, error) {
	req, err := http.NewRequest(method, c.Creds.Url+path, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Creds.Username, c.Creds.Password)
	for key := range header {
		req.Header.Set(key, header[key][0])
	}
	req.Header.Set("User-Agent", "watson-developer-cloud-go-"+goSdkVersion)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 300 {
		return b, nil
	}

	if len(b) == 0 {
		return nil, &WatsonError{Code: resp.StatusCode, Message: ""}
	}

	var e watsonError
	err = json.Unmarshal(b, &e)
	if err == nil && e.Code != 0 && len(e.Message) > 0 {
		return nil, &WatsonError{Code: e.Code, Message: e.Message}
	}
	var ae alternativeError
	err = json.Unmarshal(b, &ae)
	if err == nil && ae.Code != 0 && len(ae.Message) > 0 {
		return nil, &WatsonError{Code: ae.Code, Message: ae.Message}
	}
	var ae1 alternativeError1
	err = json.Unmarshal(b, &ae1)
	if err == nil && ae1.Code != 0 && ae1.Message != nil && len(*ae1.Message) > 0 {
		return nil, &WatsonError{Code: ae1.Code, Message: *ae1.Message}
	}

	if utf8.ValidString(string(b)) {
		return nil, &WatsonError{Code: resp.StatusCode, Message: string(b)}
	}

	return nil, errors.New("received non-20x status code and body contained invalid error JSON")
}
