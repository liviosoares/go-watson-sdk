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

package speech_to_text

import (
	"io"
	"os"
	"testing"

	"github.ibm.com/lsoares/go-watson-sdk/watson"
)

func TestListModels(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	models, err := c.ListModels()
	if err != nil {
		t.Errorf("ListModels() failed %#v\n", err)
	}
	if len(models.Models) == 0 {
		t.Errorf("ListModels() returned 0 length account slice, wanted >= 1\n")
	}
	found_US := false
	for i := range models.Models {
		if models.Models[i].Language == "en-US" {
			found_US = true
			break
		}
	}
	if !found_US {
		t.Errorf("ListModels() returned no model with \"en-US\"\n")
	}
}

func TestStream(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	output, stream, err := c.NewStream("", "audio/wav", map[string]interface{}{"continuous": true, "interim_results": false, "timestamps": false})
	if err != nil {
		t.Errorf("NewStream() failed %#v %s\n", err, err.Error())
	}

	f, err := os.Open("test_data/speech.wav")
	if err != nil {
		t.Errorf("NewStream() failed to open audio file %s %s\n", "test_data/speech.wav", err)
	}
	go func() {
		_, err = io.Copy(stream, f)
		if err != nil {
			t.Errorf("io.Copy() failed to copy audio file to API %s\n", err.Error())
		}
		stream.Close()
	}()

	select {
	case event, ok := <-output:
		if !ok || len(event.Error) > 0 {
			t.Errorf("TestStream() failed to transcribe %#v %s\n", ok, event.Error)
			return
		}
		if len(event.Results) == 0 {
			t.Errorf("TestStream() failed to transcribe, empty results %#v\n", event)
			return
		}
		if len(event.Results[0].Alternatives) == 0 {
			t.Errorf("TestStream() failed to transcribe, empty alternatives %#v\n", event.Results[0])
			return
		}
		if event.Results[0].Alternatives[0].Transcript != "thunderstorms could produce large hail isolated tornadoes and heavy rain " {
			t.Errorf("TestStream() incorrect transcription, wanted \"%s\" got \"%s\"\n",
				"thunderstorms could produce large hail isolated tornadoes and heavy rain",
				event.Results[0].Alternatives[0].Transcript)
			return
		}
	}
}
