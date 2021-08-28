# Go Coin API

## goinsapi modue

This module provides a basic api module for for gathering coinbase
information. The initial idea is to add coinbase support and build it
up to pull from additional REST/API end points.

TODO:  
[X] - Pull data from https://api.ethplorer.io/ to grap coin balances  
[X] - Create a service that utilized goapi and refreshes prices and gets balances on demand  
[X] - Implement a API server to return portfolio balance based on POST addx  
[ ] - Create a containerized microservice  
[ ] - Setup a database backend  
[X] - Refactor the Token struct to include the getters into the struct ( Not going to
      add redundant functions for each function. Marking as done.)

It utilizes a yaml config file that constains the following:

>>>
# Configuration for Coinmarketcap
api:
  url: "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest"
  tokenheader: "X-CMC_PRO_API_KEY"
  token: "YOUR-API-TOKEN-HERE"
  walletaddress: "0xBB980a1Bf6eCaaA49A6901302E97DA2C6B9dCDbe"
>>>

  Update config.yml.sample with YOUR-API-TOKEN as needed and rename it to
  config.yml.

  ## The following functions are provided
  - As part of initialization you should set a config variable ex: **var CMCApi = goapi.GetConfigData()** it returns a yamlConfig type for using when pulling data.
  - PullCoinData(yamlConfig Config) is the call to pull data from coinmarketcap.
  - GetPrice(sym string) - Accepts a ticker symbol and returns the current price as a float64.
  - GetHourChange(sym string) Accepts a ticker symbol and returns the 1 hour percent change as a float64.
  - Get24HourChange(sym string) Accepts a ticker symbol and returns the 24 hour percent change as a float64.


## Example usage

```
package main

import (
	"fmt"
	//"goapi"
	"git.local.jnet/wendell/goapi"
)

func main() {
	// Initialize the application with configuration data
	var CMCApi = goapi.GetConfigData()

	// Call to Pull Coin Data from API servers
	goapi.PullCoinData(CMCApi)

    // Call to get tokens associated with the provided ETH addx
	goapi.GetAddressData(CMCApi)

    // Assign the array of ticker symbols assocated with the tokens in the ETH addx records
	var tokens []string = goapi.GetTokenList()
```

  [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
  ![GitHub all releases](https://img.shields.io/github/downloads/FlipTheDream/goinsapi/total)
  ![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/FlipTheDream/goinsapi?include_prereleases)
  ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/FlipTheDream/goinsapi)
  ![GitHub last commit](https://img.shields.io/github/last-commit/FlipTheDream/goinsapi)
  ![GitHub issues](https://img.shields.io/github/issues/FlipTheDream/goinsapi)
  [![Go Report Card](https://goreportcard.com/badge/github.com/flipthedream/goinsapi)](https://goreportcard.com/report/github.com/flipthedream/goinsapi)