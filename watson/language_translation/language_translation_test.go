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

package language_translation

import (
	"strings"
	"testing"

	"github.com/liviosoares/go-watson-sdk/watson"
)

func TestListModels(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
	}
	models, err := c.ListModels(nil)
	if err != nil {
		t.Errorf("ListModels() failed %#v\n", err)
	}
	if len(models.Models) == 0 {
		t.Errorf("ListModels() returned 0 length account slice, wanted >= 1\n")
	}
}

func TestGetModelStatus(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	models, err := c.ListModels(nil)
	if err != nil {
		t.Errorf("ListModels() failed %#v\n", err)
		return
	}
	status, err := c.GetModelStatus(models.Models[0].ModelId)
	if err != nil {
		t.Errorf("GetModelStatus() failed %#v\n", err)
		return
	}
	if status.Status != "available" {
		t.Errorf("GetModelStatus(\"%s\") returned unexpected status. wanted %s got %s\n", models.Models[0].ModelId, "available", status.Status)
	}
}

func TestCreateModel(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}

	tmx := `<tmx version="1.4">
  <header
    creationtool="XYZTool" creationtoolversion="1.01-023"
    datatype="PlainText" segtype="sentence"
    adminlang="en-us" srclang="en"
    o-tmf="ABCTransMem"/>
  <body>
    <tu>
      <tuv xml:lang="en">
        <seg>Hello world!</seg>
      </tuv>
      <tuv xml:lang="fr">
        <seg>Bonjour tout le monde!</seg>
      </tuv>
    </tu>
  </body>
</tmx>`

	model, err := c.CreateModel("en-fr", "go-test-model", "forced_glossary", strings.NewReader(tmx))
	if err != nil {
		t.Errorf("CreateModel() failed %#v %s\n", err, err.Error())
		return
	}
	if len(model) == 0 {
		t.Errorf("CreateModel() claimed success, but returned empty model_id %#v\n", model)
	}
}

func TestDeleteModel(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}

	models, err := c.ListModels(nil)
	if err != nil {
		t.Errorf("ListModels() failed %#v\n", err)
		return
	}

	model_id := ""
	for i := range models.Models {
		if len(models.Models[i].Name) > 0 {
			model_id = models.Models[i].ModelId
			break
		}
	}

	if len(model_id) == 0 {
		t.Errorf("ListModels() could not find candidate model to delete", err)
		return
	}

	err = c.DeleteModel(model_id)
	if len(model_id) == 0 {
		t.Errorf("DeleteModel() failed to delete model %#v %s\n", err, err.Error())
		return
	}
}

func TestTranslate(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}

	translation, err := c.Translate("A sentence must have a verb", "en", "es", "")
	if err != nil {
		t.Errorf("Translate() failed %#v %s\n", err, err.Error())
		return
	}
	if len(translation.Translations) == 0 {
		t.Errorf("Translate() returned 0 translations %#v\n", translation)
		return
	}
	if translation.Translations[0].Translation != "Una sentencia debe tener un verbo" {
		t.Errorf("Translate() incorrect reply. Wanted \"%s\" got \"%s\"\n", "Una sentencia debe tener un verbo", translation.Translations[0].Translation)
		return
	}
}

func TestListIdentifiableLanguages(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}

	langs, err := c.ListIdentifiableLanguages()
	if err != nil {
		t.Errorf("ListIdentifiableLanguages() failed %#v %s\n", err, err.Error())
		return
	}
	if len(langs.Languages) == 0 {
		t.Errorf("ListIdentifiableLanguages() return 0 languages %#v\n", langs)
		return
	}
}

func TestIdentifyLanguage(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}

	langs, err := c.IdentifyLanguage("Una sentencia debe tener un verbo")
	if err != nil {
		t.Errorf("IdentifyLanguage() failed %#v %s\n", err, err.Error())
		return
	}

	if len(langs.Languages) == 0 {
		t.Errorf("IdentifyLanguage() returned 0 results %#v\n", langs)
		return
	}

	if langs.Languages[0].Language != "es" {
		t.Errorf("IdentifyLanguage() incorrect results. Wanted %s got %s\n", "es", langs.Languages[0].Language)
		return
	}
}
