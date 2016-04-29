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
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestGetAuth(t *testing.T) {
	creds := watson.Credentials{ServiceName: "dialog"}
	token, err := GetToken(creds)
	if err != nil {
		t.Errorf("GetToken() failed for %#v: %#v\n", creds, err)
	}
	if len(token) < 512 {
		t.Errorf("GetToken() returned token with short length. Wanted >%d, got %d\n", 512, len(token))
	}
}
