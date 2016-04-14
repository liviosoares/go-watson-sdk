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

package document_conversion

import (
	"strings"
	"testing"

	"github.ibm.com/lsoares/go-watson-sdk/watson"
)

func TestConvert(t *testing.T) {
	c, err := NewClient(watson.Config{})
	if err != nil {
		t.Errorf("NewClient() failed %#v\n", err)
		return
	}
	doc, err := c.Convert(NormalizedText, nil, strings.NewReader(html_doc), "text/html")
	if err != nil {
		t.Errorf("TextConvert() failed %#v %s\n", err, err.Error())
		return
	}
	if !strings.Contains(string(doc), "Document Conversion") {
		t.Errorf("TextConvert() incorrect conversion; did not find \"Document Conversion\" substring")
		return
	}
	// t.Log(string(doc))
}

const html_doc = `<!DOCTYPE html>
<html>
<head>
<title>Document Conversion Demo</title>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">

<link rel="stylesheet" href="css/style.css">
</head>
<body>
  <header class="_demo--heading">
	<div class="_demo--container">
		<a class="wordmark" href="http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/">
			<span class="wordmark--left">IBM</span>
			<span class="wordmark--right">Watson Developer Cloud</span>
		</a>
		<nav class="heading-nav" role="menubar">
			<li class="base--li heading-nav--li" role="presentation">
				<a class="heading-nav--item" href="http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/services-catalog.html" role="menuitem">
					Services
				</a>
			</li>
			<li class="base--li heading-nav--li" role="presentation">
				<a class="heading-nav--item" href="http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/doc/" role="menuitem">
					Docs
				</a>
			</li>
			<li class="base--li heading-nav--li" role="presentation">
				<a class="heading-nav--item" href="http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/gallery.html" role="menuitem">
					App Gallery
				</a>
			</li>
			<li class="base--li heading-nav--li" role="presentation">
				<a class="heading-nav--item" href="https://developer.ibm.com/watson/" role="menuitem">
					Community
				</a>
			</li>
		</nav>
	</div>
</header>

  <div class="_demo--banner">
	<div class="_demo--container">
		<div class="banner--service-icon-container">
			<img class="banner--service-icon" src="images/document-conversion.svg" alt="Document Conversion API Icon">
		</div>
		<div class="banner--service-info">
			<h1 class="banner--service-title base--h1">
				<img class="banner--service-icon_INLINE" src="images/document-conversion.svg" alt="Document Conversion API Icon">
				Document Conversion
			</h1>
			<div class="banner--service-description">
			The Document Conversion service allows you to transform one or many HTML, PDF, and Microsoft Word documents into a single and well formatted HTML document or Answer units (JSON) file that can be used to train the Retrieve and Rank service.
			</div>
			<div class="banner--service-resource">
				<span class="icon icon-link"></span>
				<strong>Resources:</strong>
			</div>
			<div class="banner--service-links">
				<li class="base--li banner--service-link-item">
					<a href="[NEED LINK HERE]" class="base--a">API Overview</a>
				</li>
				<li class="base--li banner--service-link-item">
					<a href="[NEED LINK HERE]" class="base--a">Documentation</a>
				</li>
				<li class="base--li banner--service-link-item">
					<a href="[NEED LINK HERE]" class="base--a">Fork and Deploy on Bluemix</a>
				</li>
				<li class="base--li banner--service-link-item">
					<a href="[NEED LINK HERE]" class="base--a">Fork on Github</a>
				</li>
			</div>
		</div>
	</div>
</div>
  <div class="_demo--container">
	<article class="_content base--article">
		<div class="_content--choose-input-file">
			<h2 class="base--h2">Upload Your Document</h2>
			<span class="icon-hyperlink">
			    <span class="icon icon-reset"></span>
			    <button class="base--a reset-button" href="" type="reset">
				    Reset
			    </button>
		    </span>
			<div class="_content--upload">
				<div class="upload--description">
					Upload a pdf, word(.doc, .docx) or html document or drag your document here
				</div>

			</div>
			<div class="_content--sample">
				<div class="_content--sample-title">Or use sample documents</div>
				<div class="_content--radio-group">
					<div class="_content--radio-group-item">
						<input role="radio" class="base--radio" type="radio" id="html-sample-input" name="rb" value="">
						<label class="base--inline-label" for="html-sample-input">Sample.html</label>
					</div>
					<div class="_content--radio-group-item">
						<input role="radio" class="base--radio" type="radio" id="docx-sample-input" name="rb" value="">
						<label class="base--inline-label" for="docx-sample-input">Sample.docx</label>
					</div>
					<div class="_content--radio-group-item">
						<input role="radio" class="base--radio" type="radio" id="pdf-sample-input" name="rb" value="">
						<label class="base--inline-label" for="pdf-sample-input">Sample.pdf</label>
					</div>
				</div>
			</div>
		</div>

		<div class="_content--choose-output-format">
			<h2 class="base--h2">Choose Output Format</h2>
			<div class="_content--radio-group">
				<div class="_content--radio-group-item">
					<input role="radio" class="base--radio" type="radio" id="html-sample-input" name="rb" value="">
					<label class="base--inline-label" for="html-sample-input">Answer Units JSON</label>
				</div>
				<div class="_content--radio-group-item">
					<input role="radio" class="base--radio" type="radio" id="docx-sample-input" name="rb" value="">
					<label class="base--inline-label" for="docx-sample-input">Normalized HTML</label>
				</div>
				<div class="_content--radio-group-item">
					<input role="radio" class="base--radio" type="radio" id="pdf-sample-input" name="rb" value="">
					<label class="base--inline-label" for="pdf-sample-input">Normalized Plain Text</label>
				</div>
			</div>
		</div>

		<div class="_content--choose-output">
			<div class="tab-panels" role="tabpanel">
				<ul class="tab-panels--tab-list" role="tablist">
					<li class="tab-panels--tab-list-item base--li" role="presentation">
						<a class="tab-panels--tab base--a active" href="#panel1" aria-controls="panel1" role="tab">Your Document</a>
					</li>
					<li class="tab-panels--tab-list-item base--li" role="presentation">
						<a class="tab-panels--tab base--a" href="#panel2" aria-controls="panel2" role="tab">Rest API</a>
					</li>
					<span class="icon icon-download "></span>
				</ul>
				<div class="tab-panels--tab-content">
					<div id="panel1" class="tab-panels--tab-pane active" role="tab-panel">
						<textarea class="base--textarea">this is a textarea</textarea>
					</div>
					<div id="panel2" class="tab-panels--tab-pane" role="tab-panel">
						<div class="base--textarea">this is a div</div>
					</div>
				</div>
			</div>

			<div class="tab-panels" role="tabpanel">
				<ul class="tab-panels--tab-list" role="tablist">
					<li class="tab-panels--tab-list-item base--li" role="presentation">
						<a class="tab-panels--tab base--a active" href="#panel1" aria-controls="panel1" role="tab">Output Document</a>
					</li>
					<li class="tab-panels--tab-list-item base--li" role="presentation">
						<a class="tab-panels--tab base--a" href="#panel2" aria-controls="panel2" role="tab">JSON</a>
					</li>
					<span class="icon icon-download"></span>
				</ul>
				<div class="tab-panels--tab-content">
					<div id="panel1" class="tab-panels--tab-pane active" role="tab-panel">
						<textarea class="base--textarea">this is a textarea</textarea>
					</div>
					<div id="panel2" class="tab-panels--tab-pane" role="tab-panel">
						<div class="base--textarea">this is a div</div>
					</div>
				</div>
			</div>
		</div>

		<div id="display-word" style="white-space: pre-line;">
			
		</div>

		<span class="icon-hyperlink">
			<span class="icon icon-back2top"></span>
			<a href="" class="base--a">
				Back to top
			</a>
		</span>

	</article>

</div>


  <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
  <script type="text/javascript" src="../../scss/patterns/components/tab-panels/tab-panels.js"></script>
  <script type="text/javascript" src="minjs/dist.js"></script>
</body>
</html>
`
