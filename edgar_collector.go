package valuator

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/palafrank/edgar"
)

type edgarCollector struct {
	name    string
	fetcher edgar.FilingFetcher
	store   Store
}

func NewEdgarCollector(store Store) (Collector, error) {
	//Get the Fetcher
	fetcher := edgar.NewFilingFetcher()
	if fetcher == nil {
		return nil, errors.New("Could not create edgar collector")
	}
	ret := &edgarCollector{
		fetcher: fetcher,
		store:   store,
	}
	return ret, nil
}

func (c *edgarCollector) MapEdgarFilingToValuatorFiling(fs []edgar.Filing) []Filing {
	var ret []Filing
	for _, f := range fs {
		ret = append(ret, f.(Filing))
	}
	return ret
}

func (c *edgarCollector) Name() string {
	return "Edgar Collector"
}

func (c *edgarCollector) HTML(ticker string) string {
	comp, err := c.fetcher.CompanyFolder(ticker)
	if err == nil {
		return comp.HTML(edgar.FilingType10K)
	}
	return ""

}

func (c *edgarCollector) CollectAnnualData(ticker string,
	years ...int) ([]Filing, error) {

	var err error
	var cf edgar.CompanyFolder
	var fp io.Reader = nil

	if c.store != nil {
		fp = c.store.GetFinancials(ticker)
	}

	if fp == nil {
		//If there is no historical data. Get it from Edgar.
		log.Println("No data found. Fetching from Edgar")
		cf, err = c.fetcher.CompanyFolder(ticker, edgar.FilingType10K)
	} else {
		// If data available in store use that you create folder
		cf, err = c.fetcher.CreateFolder(fp, edgar.FilingType10K)
	}
	if err != nil {
		return nil, err
	}

	// Get all the Available filings
	af := cf.AvailableFilings(edgar.FilingType10K)

	//Filter out based on what was requested
	var filteredAF []time.Time
	for _, f := range af {
		if len(years) > 0 && !contains(f.Year(), years) {
			continue
		}
		filteredAF = append(filteredAF, f)
	}

	// Get all the financial data
	if len(filteredAF) > 0 {
		fils, err := cf.Filings(edgar.FilingType10K, filteredAF...)
		if err != nil {
			return nil, err
		}

		return c.MapEdgarFilingToValuatorFiling(fils), nil
	}

	return nil, errors.New("No filings collected")
}

func (c *edgarCollector) Write(ticker string, writer io.Writer) error {

	comp, err := c.fetcher.CompanyFolder(ticker)
	if err == nil {
		return comp.SaveFolder(writer)
	}
	return err

}
