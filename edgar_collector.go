package valuator

import (
	"bytes"
	"errors"

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

	fp, err := c.db.Read(ticker)

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

func (c *collector) SaveEdgarData(ticker string) error {
	fetcher := c.fetcher.(edgar.FilingFetcher)
	comp, err := fetcher.CompanyFolder(ticker)
	if err == nil {
		data := bytes.NewBuffer(nil)
		comp.SaveFolder(data)
		if data.Len() > 0 {
			c.db.Write(ticker, data.Bytes())
		}
	}
	return nil

}
