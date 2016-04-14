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

package concept_insights

import (
	"strings"
	"testing"

	"github.ibm.com/lsoares/go-watson-sdk/watson"
)

func TestListAccounts(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	accounts, err := c.ListAccounts()
	if err != nil {
		t.Errorf("GetAccounts() failed %#v\n", err)
	}
	if len(accounts.Accounts) == 0 {
		t.Errorf("GetAccounts() returned 0 length account slice, wanted >= 1\n")
	}
}

func TestGetGraphs(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	graphs, err := c.ListGraphs()
	if err != nil {
		t.Errorf("GetAccounts() failed %#v\n", err)
	}
	if len(graphs.Graphs) < 2 {
		t.Errorf("GetGraphs() returned too short slice, wanted >= 2, got %d\n", len(graphs.Graphs))
	}
}

func TestGetConcept(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	concept, err := c.GetConcept("/graphs/wikipedia/en-20120601/concepts/IBM_Watson")
	if err != nil {
		t.Errorf("GetAccounts() failed %#v\n", err)
	}
	if concept.Label != "IBM Watson" {
		t.Errorf("GetConcept(\"/graphs/wikipedia/en-20120601/concepts/IBM_Watson\") returned concept with wrong label. Wanted %s, got %s\n", "IBM Watson", concept.Label)
	}
}

func TestSearchConceptByLabel(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	matches, err := c.SearchConceptByLabel("/graphs/wikipedia/en-20120601", "IBM", map[string]interface{}{"prefix": true, "concept_fields": "{\"abstract\":1}"})
	if err != nil {
		t.Errorf("SearchConceptByLabel() failed %#v\n", err)
	}
	if len(matches.Matches) != 10 {
		t.Errorf("SearchConceptByLabel(\"/graphs/wikipedia/en-20120601\", \"IBM\") returned smaller number of results. Wanted %d, got %d\n", 10, len(matches.Matches))
	}
}

func TestGetRelatedConcepts(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	concepts, err := c.GetRelatedConcepts("/graphs/wikipedia/en-20120601", []string{"/graphs/wikipedia/en-20120601/concepts/IBM_Watson"}, map[string]interface{}{"concept_fields": "{\"abstract\":1}"})
	if err != nil {
		t.Errorf("GetRelatedConcepts() failed %#v\n", err)
	}
	if len(concepts.Concepts) != 10 {
		t.Errorf("GetRelatedConcepts(\"/graphs/wikipedia/en-20120601\", \"IBM_Watson\") returned smaller number of results. Wanted %d, got %d\n", 10, len(concepts.Concepts))
	}
}

func TestAnnotateText(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	annotations, err := c.AnnotateText("/graphs/wikipedia/en-20120601", strings.NewReader("IBM announces new Watson cloud services."), "text/plain")
	if err != nil {
		t.Errorf("AnnotateText() failed %#v\n", err)
	}
	if len(annotations.Annotations) != 3 {
		t.Errorf("AnnotateText() returned smaller number of results. Wanted %d, got %d\n", 3, len(annotations.Annotations))
	}
}

func TestListCorpora(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	corpora, err := c.ListCorpora()
	if err != nil {
		t.Errorf("ListCorpora() failed %#v\n", err)
	}
	if len(corpora.Corpora) < 3 {
		t.Errorf("ListCorpora() returned smaller number of results. Wanted > %d, got %d\n", 3, len(corpora.Corpora))
	}
}

func TestGetCorpus(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	corpus, err := c.GetCorpus("/corpora/public/TEDTalks")
	if err != nil {
		t.Errorf("GetCorpus() failed %#v\n", err)
	}
	if corpus.Id != "/corpora/public/TEDTalks" {
		t.Errorf("GetCorpora() returned incorrect corpus id. Wanted %s, got %s\n", "/corpora/public/TEDTalks", corpus.Id)
	}
}

func TestGetCorpusProcessingState(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	state, err := c.GetCorpusProcessingState("/corpora/public/TEDTalks")
	if err != nil {
		t.Errorf("GetCorpusProcessingState() failed %#v\n", err)
	}
	if state.BuildStatus.Ready == 0 {
		t.Errorf("GetCorpusProcessingState() returned 0 documents in ready state. Wanted > %s, got %s\n", 0, state.BuildStatus.Ready)
	}
}

func TestGetCorpusStats(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	stats, err := c.GetCorpusStats("/corpora/public/TEDTalks")
	if err != nil {
		t.Errorf("GetCorpusStats() failed %#v\n", err)
	}
	if stats.TopTags.Documents == 0 {
		t.Errorf("GetCorpusStats() returned 0 documents. Wanted > %s, got %s\n", 0, stats.TopTags.Documents)
	}
}

func TestSearchCorpusByLabel(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	matches, err := c.SearchCorpusByLabel("/corpora/public/TEDTalks", "Al Gore", map[string]interface{}{"prefix": true})
	if err != nil {
		t.Errorf("SearchCorpusByLabel() failed %#v\n", err)
	}
	if len(matches.Matches) == 0 {
		t.Errorf("SearchCorpusByLabel(\"/corpora/public/TEDTalks\", \"Al Gore\") returned smaller number of results. Wanted >%d, got %d\n", 10, len(matches.Matches))
	}
}

func TestListDocuments(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	docList, err := c.ListDocuments("/corpora/public/TEDTalks", map[string]interface{}{"limit": 9})
	if err != nil {
		t.Errorf("ListDocuments() failed %#v\n", err)
	}
	if len(docList.Documents) != 9 {
		t.Errorf("ListDociments() returned incorrect number of ids. Wanted %s, got %s\n", 9, len(docList.Documents))
	}
}

func TestGetDocument(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	doc, err := c.GetDocument("/corpora/public/TEDTalks/documents/1")
	if err != nil {
		t.Errorf("GetDocument() failed %#v\n", err)
	}
	if len(doc.Parts) == 0 {
		t.Errorf("GetDocument() returned invalid number of parts. Wanted > %s, got %s\n", 0, len(doc.Parts))
	}
}

func TestGetDocumentProcessingState(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	state, err := c.GetDocumentProcessingState("/corpora/public/TEDTalks/documents/1")
	if err != nil {
		t.Errorf("GetDocumentProcessingState() failed %#v\n", err)
	}
	if state.Status != "ready" {
		t.Errorf("GetDocumentProcessingState() returned invalid status. Wanted %s, got %s\n", "ready", state.Status)
	}
}

func TestGetDocumentAnnotations(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	annotations, err := c.GetDocumentAnnotations("/corpora/public/TEDTalks/documents/1")
	if err != nil {
		t.Errorf("GetDocumentAnnotations() failed %#v\n", err)
	}
	if len(annotations.Annotations) == 0 {
		t.Errorf("GetDocumentAnnotations() returned invalid length. Wanted > %s, got %s\n", 0, len(annotations.Annotations))
	}
}

func TestGetDocumentRelatedConcepts(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	concepts, err := c.GetDocumentRelatedConcepts("/corpora/public/TEDTalks/documents/1", map[string]interface{}{"limit": 9})
	if err != nil {
		t.Errorf("GetDocumentRelatedConcepts() failed %#v\n", err)
	}
	if len(concepts.Concepts) != 9 {
		t.Errorf("GetDocumentRelatedConcepts() returned smaller number of results. Wanted %d, got %d\n", 9, len(concepts.Concepts))
	}
}

func TestGetRelatedDocuments(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	docs, err := c.GetRelatedDocuments("/corpora/public/ibmresearcher", []string{"/graphs/wikipedia/en-20120601/concepts/System_call"}, map[string]interface{}{"limit": 12})
	if err != nil {
		t.Errorf("GetRelatedDocuments() failed %#v\n", err)
	}
	if len(docs.Results) != 12 {
		t.Errorf("GetRelatedDocuments() returned smaller number of results. Wanted %d, got %d\n", 12, len(docs.Results))
	}
	if docs.Results[0].Id != "/corpora/public/ibmresearcher/documents/us-lsoares" {
		t.Errorf("GetRelatedDocuments() returned incorrect document id. Wanted %s, got %s\n", "/corpora/public/ibmresearcher/documents/us-lsoares", docs.Results[0].Id)
	}
}
