package valuator

import (
	"errors"
	"io"
)

// CollectorType is a type definition for different collector types
type CollectorType string

// CollectorEdgar is the collector defined in the edgar package
const CollectorEdgar CollectorType = "edgar"

// Collector interface provides access to the collector
type Collector interface {
	// Name of the collector
	Name() string
	// CollectAnnualData queries the collector to get annual data for the years specified
	CollectAnnualData(ticker string, year ...int) ([]Filing, error)
	// Write the Collectors content to an IO Writer
	Write(string, io.Writer) error
}

// NewCollector creates a collector to collect filings
func NewCollector(name CollectorType, store Store) (Collector, error) {

	switch name {
	case CollectorEdgar:
		return newEdgarCollector(store)
	default:
	}

	return nil, errors.New("Unsupported collector " + string(name))
}
