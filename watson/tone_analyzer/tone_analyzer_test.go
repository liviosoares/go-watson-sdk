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

package tone_analyzer

import (
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestTone(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	text := `It was the best of times. It was the worst of times. It was the age of wisdom. It was the age of foolishness. It was the epoch of belief. It was the epoch of incredulity. It was the season of Light. It was the season of Darkness. It was the spring of hope. It was the winter of despair, we had everything before us, we had nothing before us, we were all going direct to Heaven, we were all going direct the other way–in short, the period was so far like the present period, that some of its noisiest authorities insisted on its being received, for good or for evil, in the superlative degree of comparison only.`
	analysis, err := c.Tone(text, nil) // map[string]interface{}{"sentences": true})
	if err != nil {
		t.Errorf("Tone() failed %#v %s\n", err, err.Error())
		return
	}
	if len(analysis.SentencesTone) != 10 {
		t.Errorf("Tone() incorrect number of sentences returned. Wanted %d got %d\n", 10, len(analysis.SentencesTone))
		return
	}
}
