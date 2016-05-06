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

package text_to_speech

import (
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestListVoices(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	voices, err := c.ListVoices()
	if err != nil {
		t.Errorf("ListVoices() failed %#v\n", err)
		return
	}
	if len(voices.Voices) == 0 {
		t.Errorf("ListVoices() returned 0 length account slice, wanted >= 1\n")
		return
	}
}

func TestGetVoice(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	voice, err := c.GetVoice("en-US_MichaelVoice", "")
	if err != nil {
		t.Errorf("GetVoices() failed %#v\n", err)
		return
	}
	if voice.Language != "en-US" {
		t.Errorf("GetVoice() returned wrong language. Wanted %s, got %s\n", "en-US", voice.Language)
		return
	}
}

func TestGetPronuncation(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	p, err := c.GetPronunciation("Watson", "", "")
	if err != nil {
		t.Errorf("GetPronunciation() failed %#v\n", err)
		return
	}
	if p != ".ˈwɑt.sən" {
		t.Errorf("GetVoice() returned wrong language. Wanted %s, got %s\n", ".ˈwɑt.sən", p)
		return
	}
}
