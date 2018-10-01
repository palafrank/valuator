package valuator

import (
	"encoding/json"
	"log"
)

// Filing interface for fetching financial data
type Filing interface {
	Ticker() string
	FiledOn() string
	ShareCount() (float64, error)
	Revenue() (float64, error)
	CostOfRevenue() (float64, error)
	GrossMargin() (float64, error)
	OperatingIncome() (float64, error)
	OperatingExpense() (float64, error)
	NetIncome() (float64, error)
	TotalEquity() (float64, error)
	ShortTermDebt() (float64, error)
	LongTermDebt() (float64, error)
	CurrentLiabilities() (float64, error)
	DeferredRevenue() (float64, error)
	RetainedEarnings() (float64, error)
	OperatingCashFlow() (float64, error)
	CapitalExpenditure() (float64, error)
}

type Valuator interface {
	DiscountedCashFlow(ticker string) (int64, error)
	Save() error
	String() string
}

type valuation struct {
	Data map[string]Measures `json:"Date"`
}

type valuator struct {
	collector  map[string]Collector `json:"-"`
	Valuations map[string]valuation `json:"Company"`
}

func (v valuator) String() string {

	data, err := json.MarshalIndent(v.Valuations, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling financial data")
	}

	/*
		for _, val := range v.measures {
			for key1, val1 := range val {
				fmt.Println(key1, val1)
			}
		}*/
	return string(data)
}

func (v *valuator) DiscountedCashFlow(ticker string) (int64, error) {
	return 0, nil
}

func (v *valuator) Save() error {
	for ticker, collect := range v.collector {
		err := collect.Save(ticker)
		if err != nil {
			log.Println("Error saving document for ", ticker)
			return err
		}
	}
	return nil
}

func NewValuator(ticker string) (Valuator, error) {
	v := &valuator{
		collector:  make(map[string]Collector),
		Valuations: make(map[string]valuation),
	}
	v.Valuations[ticker] = valuation{
		Data: make(map[string]Measures),
	}
	collect, err := NewCollector("edgar")
	if err != nil {
		return nil, err
	}
	v.collector[ticker] = collect
	measures, err := collect.CollectAnnualData(ticker)
	if err != nil {
		return nil, err
	}
	for _, m := range measures {
		v.Valuations[ticker].Data[m.FiledOn()] = m
	}
	return v, nil
}
