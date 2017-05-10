// Display Bitcoin prices of 4 major exchange in BTC/JPY on Bitbar

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

func getJSON(url string) []byte {
	resp, herr := http.Get(url)
	if herr != nil {
		log.Print(herr)
		return nil
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil
	}
	return byteArray
}

func zaif() (result string) {
	name := "Zaif"
	byteArray := getJSON("https://api.zaif.jp/api/1/ticker/btc_jpy")
	var j struct {
		Bid float32 `json:"bid"`
		Ask float32 `json:"ask"`
	}
	if err := json.Unmarshal(byteArray, &j); err != nil {
		result = ""
	} else {
		result = fmt.Sprintf("%s\t\t\tbid: %d, ask: %d | href=https://zaif.jp", name, int(j.Bid), int(j.Ask))
	}
	return
}

func bitflyer(symbol string) (result string) {
	byteArray := getJSON("https://api.bitflyer.jp/v1/ticker?product_code=" + symbol)
	var j struct {
		Bid float32 `json:"best_bid"`
		Ask float32 `json:"best_ask"`
	}

	if err := json.Unmarshal(byteArray, &j); err != nil {
		log.Print(err)
		result = ""
	} else {
		result = fmt.Sprintf("%sbid: %d, ask: %d | href=https://bitflyer.jp", symbol, int(j.Bid), int(j.Ask))
	}
	return
}

func btcbox() (result string) {
	name := "BTCBOX"
	byteArray := getJSON("https://www.btcbox.co.jp/api/v1/ticker/")
	var j struct {
		Bid float32 `json:"buy"`
		Ask float32 `json:"sell"`
	}
	if err := json.Unmarshal(byteArray, &j); err != nil {
		log.Print(err)
		result = ""
	} else {
		result = fmt.Sprintf("%s\t\tbid: %d, ask: %d | href=https://www.btcbox.co.jp", name, int(j.Bid), int(j.Ask))
	}
	return
}

func quoine() (result string) {
	name := "Quoine"
	byteArray := getJSON("https://api.quoine.com/products/5")
	var j struct {
		Bid float32 `json:"market_bid"`
		Ask float32 `json:"market_ask"`
	}
	if err := json.Unmarshal(byteArray, &j); err != nil {
		log.Print(err)
		result = ""
	} else {
		result = fmt.Sprintf("%s\t\tbid: %d, ask: %d | href=https://trade.quoinex.com", name, int(j.Bid), int(j.Ask))
	}
	return
}

func coincheck(name string, symbol string) (result string) {
	byteArray := getJSON("https://coincheck.com/api/rate/" + symbol)
	var j struct {
		Rate string `json:"rate"`
	}
	if err := json.Unmarshal(byteArray, &j); err != nil {
		log.Print(err)
		result = ""
	} else {
		result = fmt.Sprintf("%srate: %s | href=https://coincheck.com/ja/buys?pair=eth_jpy", name, j.Rate)
	}
	return
}

func main() {
	fmt.Println("💲")

	ch := make(chan string)
	go func() { ch <- zaif() }()
	go func() { ch <- btcbox() }()
	go func() { ch <- quoine() }()

	bfch := make(chan string)
	go func() { bfch <- bitflyer("BTC_JPY\t") }()
	go func() { bfch <- bitflyer("FX_BTC_JPY\t") }()

	cch := make(chan string)
	go func() { cch <- coincheck("Ethereum\t", "eth_jpy") }()
	go func() { cch <- coincheck("EtheClassic\t", "etc_jpy") }()
	go func() { cch <- coincheck("Lisk\t\t\t", "lsk_jpy") }()
	go func() { cch <- coincheck("Factom\t\t", "fct_jpy") }()
	go func() { cch <- coincheck("Monero\t\t", "xmr_jpy") }()
	go func() { cch <- coincheck("Augur\t\t", "rep_jpy") }()
	go func() { cch <- coincheck("Ripple\t\t", "xrp_jpy") }()
	go func() { cch <- coincheck("Zcach\t\t", "zec_jpy") }()
	go func() { cch <- coincheck("NEM\t\t", "xem_jpy") }()
	go func() { cch <- coincheck("Litecoin\t\t", "ltc_jpy") }()
	go func() { cch <- coincheck("DASH\t\t", "dash_jpy") }()

	fmt.Println("---")
	fmt.Println("Bitcoin | size=16 href=http://jpbitcoin.com/markets")
	barr := []string{}
	for i := 0; i < 3; i++ {
		barr = append(barr, <-ch)
	}
	sort.Strings(barr)
	for i := 0; i < 3; i++ {
		fmt.Println(barr[i])
	}

	fmt.Println("---")
	fmt.Println("Bitfliyer | size=16 href=https://bitflyer.jp")
	bfarr := []string{}
	for i := 0; i < 2; i++ {
		bfarr = append(bfarr, <-bfch)
	}
	sort.Strings(bfarr)
	for i := 0; i < 2; i++ {
		fmt.Println(bfarr[i])
	}

	fmt.Println("---")
	fmt.Println("Coincheck | size=16 href=https://coincheck.com/ja/exchange")
	carr := []string{}
	for i := 0; i < 11; i++ {
		carr = append(carr, <-cch)
	}
	sort.Strings(carr)
	for i := 0; i < 11; i++ {
		fmt.Println(carr[i])
	}
}
