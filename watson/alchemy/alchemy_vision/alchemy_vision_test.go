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

package alchemy_vision

import (
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestGetImageKeywords(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	keywords, err := c.GetImageKeywords([]byte("http://www.nytimes.com/2015/08/16/opinion/sunday/oliver-sacks-sabbath.html"), nil)
	if err != nil {
		t.Errorf("GetImageKeywords() failed %#v %s\n", err, err.Error())
		return
	}
	if keywords.Status != "OK" {
		t.Errorf("GetImageKeywords() bad status. Wanted \"%s\" got \"%s\"", "OK", keywords.Status)
		return
	}
	// t.Logf("%+v\n", keywords)
}

func TestGetImageLink(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	link, err := c.GetImageLink([]byte("http://www.nytimes.com/2015/08/16/opinion/sunday/oliver-sacks-sabbath.html"), nil)
	if err != nil {
		t.Errorf("GetImageLink() failed %#v %s\n", err, err.Error())
		return
	}
	if link.Status != "OK" {
		t.Errorf("GetImageLink() bad status. Wanted \"%s\" got \"%s\"", "OK", link.Status)
		return
	}
	// t.Logf("%+v\n", link)
}

func TestGetImageFaceTags(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	tags, err := c.GetImageFaceTags([]byte("http://faculty.polytechnic.org/gfeldmeth/group_01.jpg"), nil)
	if err != nil {
		t.Errorf("GetImageLink() failed %#v %s\n", err, err.Error())
		return
	}
	if tags.Status != "OK" {
		t.Errorf("GetImageLink() bad status. Wanted \"%s\" got \"%s\"", "OK", tags.Status)
		return
	}
	// t.Logf("%+v\n", tags)
}
