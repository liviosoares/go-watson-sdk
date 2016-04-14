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

package alchemy_language

import (
	"testing"

	"github.ibm.com/lsoares/go-watson-sdk/watson"
)

var sampleText = `Today we are launching a campaign called for HeForShe. I am reaching out to you because we need your help. We want to end gender inequality, and to do this, we need everyone involved. This is the first campaign of its kind at the UN. We want to try to mobilize as many men and boys as possible to be advocates for change. And, we don’t just want to talk about it. We want to try and make sure that it’s tangible.  I was appointed as Goodwill Ambassador for UN Women six months ago. And, the more I spoke about feminism, the more I realized that fighting for women’s rights has too often become synonymous with man-hating. If there is one thing I know for certain, it is that this has to stop.  For the record, feminism by definition is the belief that men and women should have equal rights and opportunities. It is the theory of political, economic and social equality of the sexes.  I started questioning gender-based assumptions a long time ago. When I was 8, I was confused for being called bossy because I wanted to direct the plays that we would put on for our parents, but the boys were not. When at 14, I started to be sexualized by certain elements of the media. When at 15, my girlfriends started dropping out of sports teams because they didn’t want to appear muscly. When at 18, my male friends were unable to express their feelings.`

func TestSentiment(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	sentiment, err := c.GetSentiment([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("Sentiment() failed %#v %s\n", err, err.Error())
		return
	}
	if sentiment.Status != "OK" {
		t.Errorf("GetSentiment() bad status. Wanted \"%s\" got \"%s\"", "OK", sentiment.Status)
		return
	}
	// t.Logf("%v\n", sentiment)
}

func TestSentimentTargeted(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	sentiment, err := c.GetSentimentTargeted([]byte(sampleText), []string{"UN Women", "feminism", "muscly"}, nil)
	if err != nil {
		t.Errorf("SentimentTargetted() failed %#v %s\n", err, err.Error())
		return
	}
	if len(sentiment.Results) != 3 {
		t.Errorf("SentimentTargetted() wrong number of results. Wanted %d got %d\n", 3, len(sentiment.Results))
		return
	}
	// t.Logf("%+v\n", sentiment)
}

func TestEmotion(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	emotion, err := c.GetEmotion([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("Sentiment() failed %#v %s\n", err, err.Error())
	}
	if emotion.Status != "OK" {
		t.Errorf("GetEmotion() bad status. Wanted \"%s\" got \"%s\"", "OK", emotion.Status)
		return
	}
	// t.Logf("%+v\n", emotion)
}

func TestRankedTaxonomy(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	tax, err := c.GetTaxonomy([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("RankedTaxonomy() failed %#v %s\n", err, err.Error())
		return
	}
	if tax.Status != "OK" {
		t.Errorf("GetTaxonomy() bad status. Wanted \"%s\" got \"%s\"", "OK", tax.Status)
		return
	}
	// t.Logf("%+v\n", tax)
}

func TestConcepts(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	concepts, err := c.GetConcepts([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("GetConcepts() failed %#v %s\n", err, err.Error())
		return
	}
	if concepts.Status != "OK" {
		t.Errorf("GetConcepts() bad status. Wanted \"%s\" got \"%s\"", "OK", concepts.Status)
		return
	}
	// t.Logf("%+v\n", concepts)
}

func TestNamedEntities(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	entities, err := c.GetNamedEntities([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("NamedEntities() failed %#v %s\n", err, err.Error())
		return
	}
	if entities.Status != "OK" {
		t.Errorf("GetNamedEntities() bad status. Wanted \"%s\" got \"%s\"", "OK", entities.Status)
		return
	}
	// t.Logf("%+v\n", entities)
}

func TestKeywords(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	keywords, err := c.GetKeywords([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("GetKeywords() failed %#v %s\n", err, err.Error())
		return
	}
	if keywords.Status != "OK" {
		t.Errorf("GetKeywords() bad status. Wanted \"%s\" got \"%s\"", "OK", keywords.Status)
		return
	}
	// t.Logf("%+v\n", keywords)
}

func TestGetRelations(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	relations, err := c.GetRelations([]byte(sampleText), nil)
	if err != nil {
		t.Errorf("GetRelations() failed %#v %s\n", err, err.Error())
		return
	}
	if relations.Status != "OK" {
		t.Errorf("GetRelations() bad status. Wanted \"%s\" got \"%s\"", "OK", relations.Status)
		return
	}
	// t.Logf("%+v\n", relations)
}

func TestGetText(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	text, err := c.GetText([]byte("https://en.wikipedia.org/wiki/Watson_(computer)"), nil)
	if err != nil {
		t.Errorf("GetText() failed %#v %s\n", err, err.Error())
		return
	}
	if text.Status != "OK" {
		t.Errorf("GetText() bad status. Wanted \"%s\" got \"%s\"", "OK", text.Status)
		return
	}
	// t.Logf("%+v\n", text)
}

func TestGetRawText(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	text, err := c.GetRawText([]byte("https://en.wikipedia.org/wiki/Watson_(computer)"), nil)
	if err != nil {
		t.Errorf("GetRawText() failed %#v %s\n", err, err.Error())
		return
	}
	if text.Status != "OK" {
		t.Errorf("GetRawText() bad status. Wanted \"%s\" got \"%s\"", "OK", text.Status)
		return
	}
	// t.Logf("%+v\n", text)
}

func TestGetTitle(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	text, err := c.GetTitle([]byte("https://en.wikipedia.org/wiki/Watson_(computer)"), nil)
	if err != nil {
		t.Errorf("GetTitle() failed %#v %s\n", err, err.Error())
		return
	}
	if text.Status != "OK" {
		t.Errorf("GetTile() bad status. Wanted \"%s\" got \"%s\"", "OK", text.Status)
		return
	}
	// t.Logf("%+v\n", text)
}

func TestGetAuthor(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	author, err := c.GetAuthor([]byte("http://www.nytimes.com/2015/08/16/opinion/sunday/oliver-sacks-sabbath.html"), nil)
	if err != nil {
		t.Errorf("GetAuthor() failed %#v %s\n", err, err.Error())
		return
	}
	if author.Status != "OK" {
		t.Errorf("GetAuthor() bad status. Wanted \"%s\" got \"%s\"", "OK", author.Status)
		return
	}
	// t.Logf("%+v\n", author)
}

func TestGetAuthors(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	authors, err := c.GetAuthors([]byte("http://www.nytimes.com/2015/08/16/opinion/sunday/oliver-sacks-sabbath.html"), nil)
	if err != nil {
		t.Errorf("GetAuthors() failed %#v %s\n", err, err.Error())
		return
	}
	if authors.Status != "OK" {
		t.Errorf("GetAuthors() bad status. Wanted \"%s\" got \"%s\"", "OK", authors.Status)
		return
	}
	// t.Logf("%+v\n", authors)
}

func TestGetLanguage(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	lang, err := c.GetLanguage([]byte("Una sentencia debe tener un verbo"), nil)
	if err != nil {
		t.Errorf("GetLanguage() failed %#v %s\n", err, err.Error())
		return
	}
	if lang.Language != "spanish" {
		t.Errorf("GetLanguage() incorrect language. Wanted \"%s\" got \"%s\"\n", "spanish", lang.Language)
		return
	}
	// t.Logf("%+v\n", lang)
}

func TestGetFeedLinks(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	feeds, err := c.GetFeedLinks([]byte("http://www.nytimes.com"), nil)
	if err != nil {
		t.Errorf("GetFeedLinks() failed %#v %s\n", err, err.Error())
		return
	}
	if feeds.Status != "OK" {
		t.Errorf("GetFeedLinks() bad status. Wanted \"%s\" got \"%s\"", "OK", feeds.Status)
		return
	}	
	// t.Logf("%+v\n", feeds)
}

func TestExtractDates(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	dates, err := c.ExtractDates([]byte("http://www.nytimes.com/2015/08/16/opinion/sunday/oliver-sacks-sabbath.html"), nil)
	if err != nil {
		t.Errorf("ExtractDates() failed %#v %s\n", err, err.Error())
		return
	}
	if dates.Status != "OK" {
		t.Errorf("ExtractDates() bad status. Wanted \"%s\" got \"%s\"", "OK", dates.Status)
		return
	}	
	// t.Logf("%+v\n", dates)
}

func TestGetPubDate(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	dates, err := c.GetPubDate([]byte("http://www.nytimes.com/2015/08/16/opinion/sunday/oliver-sacks-sabbath.html"), nil)
	if err != nil {
		t.Errorf("GetPubDate() failed %#v %s\n", err, err.Error())
		return
	}
	if dates.Status != "OK" {
		t.Errorf("GetPubDate() bad status. Wanted \"%s\" got \"%s\"", "OK", dates.Status)
		return
	}	
	// t.Logf("%+v\n", dates)
}
