package valuator

import (
	"bytes"
	"errors"
	"sort"

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
	years ...int) ([]Measures, error) {

	var ret []Measures

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
		m := NewMeasures(fil)
		ret = append(ret, m)
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].FiledOn() < ret[j].FiledOn()
	})

	if len(ret) > 1 {
		for i := 1; i < len(ret); i++ {
			err := ret[i].NewYoy(ret[i-1])
			if err != nil {
				return nil, err
			}
		}
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
