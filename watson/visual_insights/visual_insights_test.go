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

package visual_insights

import (
	"os"
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestListClassifiers(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	cl, err := c.ListClassifiers()
	if err != nil {
		t.Errorf("ListClassifiers() failed %#v %s\n", err, err.Error())
		return
	}
	if len(cl.Classifiers) == 0 {
		t.Errorf("ListClassifiers() returned 0 length account slice, wanted >= 1\n")
		return
	}
	// t.Logf("%#v\n", cl)
}

func TestSummarize(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}

	f, err := os.Open("test_data/images.zip")
	if err != nil {
		t.Errorf("TestSummarize() could not open \"test_data/images.zip\" %#v\n", err)
		return
	}
	defer f.Close()

	s, err := c.Summarize(f)
	if err != nil {
		t.Errorf("Summarize() failed %#v %s\n", err, err.Error())
		return
	}

	features := 0
	for i := range s.Summary {
		if s.Summary[i].Score > 0 {
			// t.Log(s.Summary[i])
			features++
		}
	}
	// t.Logf("%#v\n", s)
	if features == 0 {
		t.Errorf("Summarize() returned only all classifiers with 0 score.\n")
		return
	}
}
