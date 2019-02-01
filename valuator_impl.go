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
	Dividend() (float64, error)
	DividendPerShare() (float64, error)
	WAShares() (float64, error)
}

type valuation struct {
	FiledData []Measures `json:"Measures"`
	Avgs      Average    `json:"Averages"`
}

type valuator struct {
	db         database
	collector  map[string]Collector  `json:"-"`
	Valuations map[string]*valuation `json:"Company"`
}

func (v valuator) String() string {

	data, err := json.MarshalIndent(v.Valuations, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuator data: ", err)
	}
	return string(data)
}

func (v valuation) String() string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuation data: ", err)
	}
	return string(data)
}

func (v *valuator) Save() error {
	for ticker, collect := range v.collector {
		err := collect.Save(ticker)
		if err != nil {
			log.Println("Error saving document for ", ticker)
			return err
		}
	}
	for ticker, valuation := range v.Valuations {
		v.db.Write(ticker+"_val", []byte(valuation.String()))
	}
	return nil
}
