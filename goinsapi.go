package goinsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/yaml.v3"
)

// ####################################################
//  Type Definitions
// ####################################################

// coinList structure defines the stucture of the CMC APPI
// The structure captures all coins and assocated statistics
// the structure is used when grabbing data from CMC
// https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest
type coinList struct {
	Status struct {
		Timestamp    time.Time   `json:"timestamp"`
		ErrorCode    int         `json:"error_code"`
		ErrorMessage interface{} `json:"error_message"`
		Elapsed      int         `json:"elapsed"`
		CreditCount  int         `json:"credit_count"`
		Notice       interface{} `json:"notice"`
		TotalCount   int         `json:"total_count"`
	} `json:"status"`
	Data []struct {
		ID                int         `json:"id"`
		Name              string      `json:"name"`
		Symbol            string      `json:"symbol"`
		Slug              string      `json:"slug"`
		NumMarketPairs    int         `json:"num_market_pairs"`
		DateAdded         time.Time   `json:"date_added"`
		Tags              []string    `json:"tags"`
		MaxSupply         int         `json:"max_supply"`
		CirculatingSupply int         `json:"circulating_supply"`
		TotalSupply       int         `json:"total_supply"`
		Platform          interface{} `json:"platform"`
		CmcRank           int         `json:"cmc_rank"`
		LastUpdated       time.Time   `json:"last_updated"`
		Quote             struct {
			Usd struct {
				Price            float64   `json:"price"`
				Volume24H        float64   `json:"volume_24h"`
				PercentChange1H  float64   `json:"percent_change_1h"`
				PercentChange24H float64   `json:"percent_change_24h"`
				PercentChange7D  float64   `json:"percent_change_7d"`
				PercentChange30D float64   `json:"percent_change_30d"`
				PercentChange60D float64   `json:"percent_change_60d"`
				PercentChange90D float64   `json:"percent_change_90d"`
				MarketCap        float64   `json:"market_cap"`
				LastUpdated      time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

// coinPricing type is provided for the coin specific pricing information.
type coinPricing struct {
	Price            float64
	Volume24H        float64
	PercentChange1H  float64
	PercentChange24H float64
	PercentChange7D  float64
	PercentChange30D float64
	PercentChange60D float64
	PercentChange90D float64
	MarketCap        float64
	LastUpdated      time.Time
}

// apiexplorer struct to grab coin quantities
// given an etherum address.
// Obtaining information example:
// https://api.ethplorer.io/getAddressInfo/[ETH_ADDRESS_HERE]?apiKey=freekey
type apiExplorer struct {
	Address string `json:"address"`
	Eth     struct {
		Price struct {
			Rate            float64 `json:"rate"`
			Diff            float64 `json:"diff"`
			Diff7D          float64 `json:"diff7d"`
			Ts              int     `json:"ts"`
			MarketCapUsd    float64 `json:"marketCapUsd"`
			AvailableSupply float64 `json:"availableSupply"`
			Volume24H       float64 `json:"volume24h"`
			Diff30D         float64 `json:"diff30d"`
			VolDiff1        float64 `json:"volDiff1"`
			VolDiff7        float64 `json:"volDiff7"`
			VolDiff30       float64 `json:"volDiff30"`
		} `json:"price"`
		RawBalance float64 `json:"rawBalance"`
	} `json:"ETH"`
	CountTxs int `json:"countTxs"`
	Tokens   []struct {
		TokenInfo struct {
			Address           string `json:"address"`
			Name              string `json:"name"`
			Decimals          string `json:"decimals"`
			Symbol            string `json:"symbol"`
			TotalSupply       string `json:"totalSupply"`
			Owner             string `json:"owner"`
			LastUpdated       int    `json:"lastUpdated"`
			IssuancesCount    int    `json:"issuancesCount"`
			HoldersCount      int    `json:"holdersCount"`
			Description       string `json:"description"`
			Website           string `json:"website"`
			Twitter           string `json:"twitter"`
			Reddit            string `json:"reddit"`
			Telegram          string `json:"telegram"`
			Image             string `json:"image"`
			Coingecko         string `json:"coingecko"`
			EthTransfersCount int    `json:"ethTransfersCount"`
			Price             struct {
				Rate            float64 `json:"rate"`
				Diff            float64 `json:"diff"`
				Diff7D          float64 `json:"diff7d"`
				Ts              int     `json:"ts"`
				MarketCapUsd    float64 `json:"marketCapUsd"`
				AvailableSupply int     `json:"availableSupply"`
				Volume24H       float64 `json:"volume24h"`
				Diff30D         float64 `json:"diff30d"`
				VolDiff1        float64 `json:"volDiff1"`
				VolDiff7        float64 `json:"volDiff7"`
				VolDiff30       float64 `json:"volDiff30"`
				Currency        string  `json:"currency"`
			} `json:"price"`
			PublicTags []string `json:"publicTags"`
		} `json:"tokenInfo,omitempty"`
		Balance    float32 `json:"balance"`
		TotalIn    int64   `json:"totalIn"`
		TotalOut   int64   `json:"totalOut"`
		RawBalance string  `json:"rawBalance"`
	} `json:"tokens"`
}

// tokenInfo type contains the apiexploter coin information that
// provides advanced details beyond the basic CMC information.\
// TODO: Refactor this scruct to include the getters into the struct
type tokenInfo struct {
	Address           string  `json:"address"`
	Name              string  `json:"name"`
	Decimals          string  `json:"decimals"`
	Symbol            string  `json:"symbol"`
	TotalSupply       string  `json:"totalSupply"`
	Owner             string  `json:"owner"`
	LastUpdated       int     `json:"lastUpdated"`
	IssuancesCount    int     `json:"issuancesCount"`
	HoldersCount      int     `json:"holdersCount"`
	Description       string  `json:"description"`
	Website           string  `json:"website"`
	Twitter           string  `json:"twitter"`
	Reddit            string  `json:"reddit"`
	Telegram          string  `json:"telegram"`
	Image             string  `json:"image"`
	Coingecko         string  `json:"coingecko"`
	EthTransfersCount int     `json:"ethTransfersCount"`
	Balance           float32 `json:"balance"`
	RawBalance        string  `json:"rawbalance"`
	Price             struct {
		Rate            float64 `json:"rate"`
		Diff            float64 `json:"diff"`
		Diff7D          float64 `json:"diff7d"`
		Ts              int     `json:"ts"`
		MarketCapUsd    float64 `json:"marketCapUsd"`
		AvailableSupply int     `json:"availableSupply"`
		Volume24H       float64 `json:"volume24h"`
		Diff30D         float64 `json:"diff30d"`
		VolDiff1        float64 `json:"volDiff1"`
		VolDiff7        float64 `json:"volDiff7"`
		VolDiff30       float64 `json:"volDiff30"`
		Currency        string  `json:"currency"`
	}
}

// Config structure provides a simple structure for the CMC API
// url - The url for the API
// tokenheader		- The API token header syntax used to authenticate
// token			- The actual secret token (This should be protected!)
type Config struct {
	Api struct {
		ApiUrl         string `yaml:"url"`
		ApiTokenHeader string `yaml:"tokenheader"`
		ApiToken       string `yaml:"token"`
		ApiAddress     string `yaml:"walletaddress"`
		ApiPort        string `yaml:"port"`
	} `yaml:"api"`
}

// ####################################################
//  Internal Package Variables
// ####################################################
var cList coinList
var tokenList apiExplorer
var coinPrice = make(map[string]float64)
var coinPrices = make(map[string]coinPricing)
var activeCoins = make(map[string]tokenInfo)
var configData Config
var addressBalance float32

// ####################################################
//  Exported Package Variables
// ####################################################
var ActiveTokenList []string

// ####################################################
//  Internal Functions Block
// ####################################################

// updatePriceMap populates a coin price map for quick lookups
func updateTokenBalanceMap(addx string) {
	// Update the map
	// Delete the exiting map first and then update it.
	for i := range activeCoins {
		delete(activeCoins, i)
	}

	var tokenData tokenInfo

	// Add ETH back into the activeCoins slice
	getETHData(addx)

	// Update the map with the coin data from the API
	for _, cToken := range tokenList.Tokens {
		tokenData.Address = cToken.TokenInfo.Address
		tokenData.Name = cToken.TokenInfo.Name
		tokenData.Decimals = cToken.TokenInfo.Decimals
		tokenData.Symbol = cToken.TokenInfo.Symbol
		tokenData.TotalSupply = cToken.TokenInfo.TotalSupply
		tokenData.Owner = cToken.TokenInfo.Owner
		tokenData.LastUpdated = cToken.TokenInfo.LastUpdated
		tokenData.IssuancesCount = cToken.TokenInfo.IssuancesCount
		tokenData.HoldersCount = cToken.TokenInfo.HoldersCount
		tokenData.Description = cToken.TokenInfo.Description
		tokenData.Website = cToken.TokenInfo.Website
		tokenData.Twitter = cToken.TokenInfo.Twitter
		tokenData.Reddit = cToken.TokenInfo.Reddit
		tokenData.Telegram = cToken.TokenInfo.Telegram
		tokenData.Image = cToken.TokenInfo.Image
		tokenData.Coingecko = cToken.TokenInfo.Coingecko
		tokenData.EthTransfersCount = cToken.TokenInfo.EthTransfersCount
		tokenData.RawBalance = cToken.RawBalance
		tokenData.Balance = cToken.Balance

		if tokenData.Website != "" {
			activeCoins[tokenData.Symbol] = tokenData
		}
	}

}

// Get ETH Data utilizes a geth light node to pull data for the actual
// ETH balance information. This might be able to be expanded into
// providing all coin information in the future rather than utilizing
// API calls to third parties.
func getETHData(addx string) {
	// Get the ETH balance first
	var ethData tokenInfo

	ethData.Balance = float32(tokenList.Eth.RawBalance)
	ethData.Image = "https://ethplorer.io/images/eth.png"
	ethData.Address = addx
	ethData.Website = "https://ethereum.org/"
	ethData.Symbol = "ETH"
	ethData.Name = "Ethereum"

	// Initialize ETH client
	// ethclient, err := ethclient.Dial("https://mainnet.infura.io")
	ethclient, err := ethclient.Dial("http://192.168.1.7:8545")
	if err != nil {
		log.Fatal(err)
		return
	}
	account := common.HexToAddress(addx)

	balance, err := ethclient.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	//bigBalance := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	bigBalance := new(big.Float).Quo(fbalance, big.NewFloat(1.0))

	value, err := strconv.ParseFloat(bigBalance.String(), 32)
	if err != nil {
		// do something sensible
	}
	ethData.Balance = float32(value)
	activeCoins["ETH"] = ethData
}

// updatePriceMap populates a coin price map for quick lookups
func updatePriceMap() {
	var coinPriceData coinPricing

	// Update the coinPrice map
	// Delete the exiting map first and then update it.
	// for i := range coinPrice {
	// 	delete(coinPrice, i)
	// }

	// Update the map with the coin data from the API
	for _, cPrice := range cList.Data {
		coinPriceData.LastUpdated = cPrice.LastUpdated
		coinPriceData.Price = cPrice.Quote.Usd.Price
		coinPriceData.Volume24H = cPrice.Quote.Usd.Volume24H
		coinPriceData.PercentChange1H = cPrice.Quote.Usd.PercentChange1H
		coinPriceData.PercentChange24H = cPrice.Quote.Usd.PercentChange24H
		coinPriceData.PercentChange7D = cPrice.Quote.Usd.PercentChange7D
		coinPriceData.PercentChange30D = cPrice.Quote.Usd.PercentChange30D
		coinPriceData.PercentChange60D = cPrice.Quote.Usd.PercentChange60D
		coinPriceData.PercentChange90D = cPrice.Quote.Usd.PercentChange90D
		coinPriceData.MarketCap = cPrice.Quote.Usd.MarketCap
		coinPriceData.LastUpdated = cPrice.Quote.Usd.LastUpdated
		coinPrices[cPrice.Symbol] = coinPriceData
	}
}

// ####################################################
//  Exported Package Functions
// ####################################################

// pullCoinData accepts a basic API url/token config structure
// The coinList structure will be used to populate the cList variable
func PullCoinData(yamlConfig Config) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", yamlConfig.Api.ApiUrl, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	configData = yamlConfig
	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", yamlConfig.Api.ApiToken)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	// fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(respBody))

	json.Unmarshal([]byte(respBody), &cList)
	updatePriceMap()
}

// ####################################################
//  External Getter Functions Block
// ####################################################

// Returns the price of a coin given its ticket symbol
func GetPrice(sym string) float64 {
	return coinPrices[sym].Price

}

// Returns the 1 hour change of a coin given its ticker symbol
func GetHourChange(sym string) float64 {
	return coinPrices[sym].PercentChange1H
}

// Returns the 24 hour change of a coin given its ticker symbol
func Get24HourChange(sym string) float64 {
	return coinPrices[sym].PercentChange24H
}

// getConfigData Reads the API configuration from a config yaml file
func GetConfigData() Config {
	var getConfig Config
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&getConfig)
	if err != nil {
		fmt.Println(err)
	}
	return getConfig

}

// Call to get all ERC-20 coins assocated with a specific ERC-20 address
// provided in the config file.
func GetAddressData(addx string) {

	if len(addx) < 20 || len(addx) > 42 {
		return
	}

	// Clear out the current Active Token List
	ActiveTokenList = nil
	client := &http.Client{}
	url := "https://api.ethplorer.io/getAddressInfo/" + addx + "?apiKey=freekey"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		fmt.Println("Failed getting address information.")
		os.Exit(1)
	}

	req.Header.Set("Accepts", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	// fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBody))

	json.Unmarshal([]byte(respBody), &tokenList)
	updateTokenBalanceMap(addx)

	// Reset addressBalance
	addressBalance = 0.0
	for sym := range activeCoins {
		ActiveTokenList = append(ActiveTokenList, sym)
		//AddressBalance += activeCoins[sym].Balance * float32(activeCoins[sym].Price)
		addressBalance += float32(GetPrice(sym)) * GetTokenBalance(sym)
	}

}

// GetTokenBalance returns the balance of a coin given its ticker symbol
func GetTokenBalance(sym string) float32 {
	return activeCoins[sym].Balance / float32(math.Pow10(18))
}

// GetTokenList returns a string array of the current tokens associated
// with a provided ERC-20 address
func GetTokenList() []string {
	return ActiveTokenList
}

// GetTokenName returns the friendly name of a coin given its ticker symbol
func GetTokenName(sym string) string {
	return activeCoins[sym].Name
}

// GetTokenChange returns the 1 hour price difference of a coin given its ticker symbol
func GetTokenChange(sym string) float32 {
	return float32(activeCoins[sym].Price.Diff)
}

// GetTokenUrl returns the website of a coin given its ticker symbol
func GetTokenUrl(sym string) string {
	return activeCoins[sym].Website
}

//
// GetAddressBalance tracks the balance of coins associated with an address
// In the future a map can be created to track each balance per address
//
func GetAddressBalance() float32 {
	return addressBalance
}

//
// GetTokenImageUrl returns the url to the token image
//
func GetTokenImageUrl(sym string) string {
	if sym == "ETH" {
		return activeCoins[sym].Image
	}
	return "https://ethplorer.io" + activeCoins[sym].Image
}

// Price data refresh ticker
// This sould be called in the main function similar to the following
// go goapi.RefreshPrice()
func RefreshPrice() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		PullCoinData(configData)
	}
}
