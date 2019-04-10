package valuator

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

var (
	priceURLFormatString = `https://api.iextrading.com/1.0/stock/%s/price`
	historicalPrice      = `https://api.iextrading.com/1.0/stock/%s/chart/date/%s`
)

func round(val float64) float64 {
	return math.Floor(val*100) / 100
}

func percentage(val float64) float64 {
	return math.Floor(val * 100)
}

func yoyCalc(past float64, curr float64, pc bool) float64 {
	if past == 0 {
		return 0
	}
	if pc && past != 0 {
		return percentage((curr - past) / past)
	}
	return round(curr - past)
}

func avgCalc(total float64, length float64) float64 {
	return math.Floor(total*100/length) / 100
}

func contains(key int, db []int) bool {
	for _, d := range db {
		if d == key {
			return true
		}
	}
	return false
}

func priceFetcher(ticker string) float64 {

	url := fmt.Sprintf(priceURLFormatString, ticker)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("GET:" + err.Error())
		return -1
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("READ:" + err.Error())
		return -1
	}
	price, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0
	}
	return price

}
