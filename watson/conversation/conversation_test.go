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

package conversation

import (
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestMessage(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	reply, err := c.Message("1a4119da-d9d3-4ec0-ab4f-bd8ca287e88c", "What's the store phone number?")
	if err != nil {
		t.Errorf("Message() failed %#v\n", err)
		return
	}
	t.Logf("%#v\n", reply)
	if len(reply.Intents) == 0 {
		t.Errorf("Message() failed.  0 intents returned, expected => 1 intents: %#v\n", reply)
		return
	}
}
