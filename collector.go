package valuator

import (
	"errors"
	"io"
)

type FilingDate string
type FilingType string
type CollectorType string

var (
	collectorEdgar CollectorType = "edgar"
)

type Collector interface {
	// Name of the collector
	Name() string
	// CollectAnnualData queries the collector to get annual data for the years specified
	CollectAnnualData(ticker string, year ...int) ([]Filing, error)
	// Write the Collectors content to an IO Writer
	Write(string, io.Writer) error
}

func NewCollector(name CollectorType, store Store) (Collector, error) {

	switch name {
	case collectorEdgar:
		return NewEdgarCollector(store)
	default:
	}

	return nil, errors.New("Unsupported collector " + string(name))
}
