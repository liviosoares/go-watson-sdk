# Go SDK for IBM Watson services
============================================
[![GoDoc](https://godoc.org/github.com/liviosoares/go-watson-sdk?status.svg)](https://godoc.org/github.com/liviosoares/go-watson-sdk)
[![CLA assistant](https://cla-assistant.io/readme/badge/liviosoares/go-watson-sdk)](https://cla-assistant.io/liviosoares/go-watson-sdk)

Go (golang) client library to use the [Watson Developer Cloud][wdc] services.

## Table of Contents
   * [Installation](#installation)
   * [Documentation](#documentation)
   * [Basic Usage](#basic-usage)
   * [Testing](#testing)
   * [License](#license)
   * [Contributing](#contributing)
   * [Authors and Contributors] (#authors)

## Installation
Use the go tool to install the package (and a couple of dependencies):
```
go get github.com/liviosoares/go-watson-sdk/...
```

## Documentation
Go API documentation @ godoc.org: https://godoc.org/github.com/liviosoares/go-watson-sdk
Watson Developer Cloud: http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/doc/

## Basic Usage
For complete documentation, see the references above to the [Godoc](https://godoc.org/github.com/liviosoares/go-watson-sdk) documentation.

A short example of connecting to the Watson Concept Insights service, using the `watson` package:

	import (
	       "github.com/liviosoares/go-watson-sdk/watson"
	       "github.com/liviosoares/go-watson-sdk/watson/concept_insights"
	)

	// assumes $VCAP_SERVICES environment variable is set
	client, err := concept_insights.NewClient(watson.Config{})
	if err != nil {
		return err
	}
	// Gets the accounts using the `/concept-insights/api/v2/accounts` endpoint
	accounts, err := client.ListAccounts()
	if err != nil {
		return err
	}
	fmt.Println(accounts)

If you do not want to rely on the setup of the `$VCAP_SERVICES` environment variable, it's also possible to pass credential information explicitly in the `NewClient()` call:

	import (
	       "github.com/liviosoares/go-watson-sdk/watson"
	       "github.com/liviosoares/go-watson-sdk/watson/concept_insights"
	)

	config := watson.Config{
		"url": "https://gateway.watsonplatform.net/concept-insights/api",
		"username": "... <user id> ...",
		"password": " ... <password> ..."
	}
	client, err := concept_insights.NewClient(config)

## Testing

To test the SDK, you must first obtain credentials for the specific services you
would like to test. Please see the Watson Developer Cloud documentation to
obtain credentials: https://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/doc/getting_started/gs-credentials.shtml

The `.creds.sh` file is a conveniently written version of the `VCAP_SERVICES` environment variable that can be sourced in a bash or sh shell session. Use your credentials to fill the specific service definitions in `.creds.sh` and source the file. The package will be able to parse and query the `VCAP_SERVICES` environment variable, extracting the appropriate username and password.

Then, to test all services use:
```
go test github.com/liviosoares/go-watson-sdk/watson/...
```

and to test a specific service, you can do, for example
```
go test github.com/liviosoares/go-watson-sdk/watson/concept_insights
```

## License
This library is licensed under Apache 2.0. Full license text is available in
[LICENSE](LICENSE).

## Contributing

Issues and pull requests are encouraged! Whenever applicable, please add a test in the appropriate service implementation, using one of the `_test.go` test files (or creating a new one).

Also, please review and digitally sign our Contributor License Agreement here: [[go-watson-sdk CLA](https://cla-assistant.io/readme/badge/liviosoares/go-watson-sdk)](https://cla-assistant.io/liviosoares/go-watson-sdk)

## Authors

* Livio Soares liviobs@gmail.com

[wdc]: http://www.ibm.com/smarterplanet/us/en/ibmwatson/developercloud/
