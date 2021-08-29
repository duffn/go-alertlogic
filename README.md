# go-alertlogic

![CI](https://github.com/duffn/go-alertlogic/actions/workflows/ci.yml/badge.svg) [![codecov](https://codecov.io/gh/duffn/go-alertlogic/branch/main/graph/badge.svg?token=wH2QcSPvpn)](https://codecov.io/gh/duffn/go-alertlogic) [![Go Reference](https://pkg.go.dev/badge/github.com/duffn/go-alertlogic.svg)](https://pkg.go.dev/github.com/duffn/go-alertlogic)

`go-alertlogic` is a Go client library for the Alert Logic Cloud Insight API.

This is in _very early_ development and only supports a few of the [myriad of endpoints](https://console.cloudinsight.alertlogic.com/api/#/) of the API.

## Installation

```bash
go get github.com/duffn/go-alertlogic
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/duffn/go-alertlogic"
)

func main() {
	// Create an API instance.
	api, err := alertlogic.NewWithUsernameAndPassword(
		os.Getenv("ALERTLOGIC_ACCOUNT_ID"),
		os.Getenv("ALERTLOGIC_USERNAME"),
		os.Getenv("ALERTLOGIC_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get your account details.
	accountDetails, err := api.GetAccountDetails()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", accountDetails)
}

```

## Documentation

[https://pkg.go.dev/github.com/duffn/go-alertlogic](https://pkg.go.dev/github.com/duffn/go-alertlogic)

## License

[MIT](https://opensource.org/licenses/MIT)
