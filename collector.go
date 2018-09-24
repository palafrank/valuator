package valuator

import (
	"errors"
)

type FilingDate string
type FilingType string

var (
	collectorEdgar = "edgar"
)

type Collector interface {
	CollectAnnualData(ticker string, year ...int) (map[int]Measures, error)
	Save(string) error
}

type Fetcher interface {
}

type collector struct {
	name    string
	db      database
	fetcher Fetcher
}

func (c *collector) Name() string {
	return c.name
}

func (c *collector) CollectAnnualData(ticker string,
	years ...int) (map[int]Measures, error) {
	switch c.Name() {
	case collectorEdgar:
		return c.CollectEdgarAnnualData(ticker, years...)
	default:
	}
	return nil, errors.New("Unknown collector type")
}

func (c *collector) Save(ticker string) error {
	switch c.Name() {
	case collectorEdgar:
		return c.SaveEdgarData(ticker)
	default:
	}
	return nil
}

func NewCollector(name string) (Collector, error) {
	c := new(collector)
	switch name {
	case collectorEdgar:
		f, err := NewEdgarCollector()
		if err == nil {
			c.fetcher = f
			c.name = name
			c.db = NewDB(fileDBUrl, FileDatabaseType)
			return c, nil
		}
	default:
	}

	return nil, errors.New("Unsupported collector")
}
