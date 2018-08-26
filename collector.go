package valuator

import (
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
