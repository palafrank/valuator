package valuator

import (
	"errors"
	"os"

	"github.com/palafrank/edgar"
)

func NewEdgarCollector() (edgar.FilingFetcher, error) {
	//Get the Fetcher
	fetcher := edgar.NewFilingFetcher()
	if fetcher == nil {
		return nil, errors.New("Could not get edgar data")
	}
	return fetcher, nil
}

func contains(key int, db []int) bool {
	for _, d := range db {
		if d == key {
			return true
		}
	}
	return false
}

func (c *collector) CollectEdgarAnnualData(ticker string,
	years ...int) (map[int]Measures, error) {

	ret := make(map[int]Measures)

	fetcher := c.fetcher.(edgar.FilingFetcher)

	fp, err := os.Open("./db/" + ticker + ".json")

	var cf edgar.CompanyFolder
	//If there is no historical data. Get it from Edgar.
	if err != nil {
		cf, err = fetcher.CompanyFolder(ticker, edgar.FilingType10K)
	} else {
		cf, err = fetcher.CreateFolder(fp, edgar.FilingType10K)
	}

	af := cf.AvailableFilings(edgar.FilingType10K)

	for _, f := range af {
		if len(years) > 0 && !contains(f.Year(), years) {
			continue
		}
		fil, err := cf.Filing(edgar.FilingType10K, f)
		if err != nil {
			return nil, err
		}
		m := getMeasures(fil)
		ret[f.Year()] = m
	}

	return ret, err
}
