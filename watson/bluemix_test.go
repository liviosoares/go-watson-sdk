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

package watson

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"
)

var vcapTests = []struct {
	vcap_services string
	service_name  string
	service_plan  string
	err           error
	want          Credentials
}{
	{
		vcap_services: ``,
		service_name:  "concept_insights",
		service_plan:  "",
		err:           errors.New("VCAP_SERVICES undefined"),
		want:          Credentials{},
	},
	{
		vcap_services: `foo=bar`,
		service_name:  "concept_insights",
		service_plan:  "",
		err:           errors.New("failed to parse VCAP_SERVICES "),
		want:          Credentials{},
	},
	{
		vcap_services: `
{
   "concept_insights": [
      {
         "name": "concept-insights-service",
         "label": "concept_insights",
         "plan": "standard",
         "credentials": {
            "url": "https://gateway.watsonplatform.net/concept-insights/api",
            "username": "uuuu",
            "password": "pppp"
         }
      }
   ]
}
`,
		service_name: "concept_insights",
		service_plan: "",
		want: Credentials{
			Url:         "https://gateway.watsonplatform.net/concept-insights/api",
			Username:    "uuuu",
			Password:    "pppp",
			ServiceName: "concept_insights",
			ServicePlan: "standard",
		},
	},
}

func TestGetCredentials(t *testing.T) {
	for i := range vcapTests {
		os.Setenv("VCAP_SERVICES", vcapTests[i].vcap_services)
		creds, err := getBluemixCredentials(vcapTests[i].service_name, vcapTests[i].service_plan)
		if (err != nil && vcapTests[i].err == nil) || (err == nil && vcapTests[i].err != nil) {
			t.Errorf("credentials test %d:\nVCAP_SERVICES = %s\ngot error %#v\nwanted %#v\n", i, vcapTests[i].vcap_services, err, vcapTests[i].err)
			return
		}
		if err != nil && vcapTests[i].err != nil {
			if !strings.HasPrefix(err.Error(), vcapTests[i].err.Error()) {
				t.Errorf("credentials test %d:\nVCAP_SERVICES = %s\ngot error %#v\nwanted %#v\n", i, vcapTests[i].vcap_services, err.Error(), vcapTests[i].err)
				return
			}
		}
		if !reflect.DeepEqual(creds, vcapTests[i].want) {
			t.Errorf("credentials test %d:\nVCAP_SERVICES = %s\ngot %#v\nwanted %#v\n", i, vcapTests[i].vcap_services, creds, vcapTests[i].want)
			return
		}
	}
}
