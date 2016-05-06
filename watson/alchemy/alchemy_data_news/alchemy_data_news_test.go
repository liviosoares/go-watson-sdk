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

package alchemy_data_news

import (
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestGetNews(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	_, err = c.GetNews("now-24h", "now", nil)
	if err != nil {
		t.Errorf("GetNews() failed %#v %s\n", err, err.Error())
		return
	}	
	// t.Logf("%+v\n", news)
}

func TestGetNews2(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	_, err = c.GetNews("now-24h", "now", map[string]interface{}{"q.enriched.url.title": "[baseball^soccer]", "return": "enriched.url.title,enriched.url.author,original.url"})
	if err != nil {
		t.Errorf("GetNews2() failed %#v %s\n", err, err.Error())
		return
	}
	// t.Logf("%+v\n", news)
}
