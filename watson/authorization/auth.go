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

package authorization

import (
	"net/url"

	"github.com/liviosoares/go-watson-sdk/watson"
)

// Uses /authorization endpoint to obtain a token for the service described in creds
func GetToken(creds watson.Credentials) (string, error) {
	serviceClient, err := watson.NewClient(creds)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(serviceClient.Creds.Url)
	if err != nil {
		return "", err
	}
	u.Path = ""

	// set up authClient with simular credential information as original service
	authClient := new(watson.Client)
	*authClient = *serviceClient
	authClient.Creds.Url = u.String()

	b, err := authClient.MakeRequest("GET", "/authorization/api/v1/token?url="+serviceClient.Creds.Url, nil, nil)
	if err != nil {
		return "", err
	}
	return url.QueryUnescape(string(b))
}
