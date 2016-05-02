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

/*
Package `watson` is the Watson SDK for the Go programming language.

Usage of the driver begins with retrieving a service specific Client object:

	// assumes $VCAP_SERVICES environment variable has been set
	client, err := concept_insights.NewClient(watson.Config{})

or

	// explicitly supply credential information
	config := watson.Config{
		Credentials: watson.Credentials{
			Username: "... username ...",
			Password: "... password ...",
		},
	}
	client, err := concept_insights.NewClient(config)

The `watson` package contains generic and common code to handle making requests to services, handling errors and dealing with Cloudfoundry and Bluemix environment variables and credentials (most notably the $VCAP_SERVICES environment variable).

The sub-packages are all service specific and provide access and functionality to the service. Examples

Alchemy Language

	import "github.com/liviosoares/go-watson-sdk/watson"
	import "github.com/liviosoares/go-watson-sdk/watson/alchemy/alchemy_language"

	client, err := alchemy_language.NewClient(watson.Config{})
	if err != nil {
		return nil
	}
	sentiment, err := c.GetSentiment([]byte("..."))
	if err != nil {
		return nil
	}
	fmt.Println(sentiment)

Concept Insights
	
	import "github.com/liviosoares/go-watson-sdk/watson"
	import "github.com/liviosoares/go-watson-sdk/watson/concept_insights"

	client, err := concept_insights.NewClient(watson.Config{})
	if err != nil {
		return nil
	}
	annotations, err := c.AnnotateText("/graphs/wikipedia/en-20120601", strings.NewReader("IBM announces new Watson cloud services."), "text/plain")
	if err != nil {
		return nil
	}
	fmt.Println(annotations)


See sub-package specific documentation for more examples.

*/
package doc
