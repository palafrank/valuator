package valuator

import (
	"bytes"
	"encoding/json"
	"log"
	"time"
)

// Filing interface for fetching financial data
type Filing interface {
	Ticker() string
	FiledOn() time.Time
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
	CurrentAssets() (float64, error)
	DeferredRevenue() (float64, error)
	RetainedEarnings() (float64, error)
	OperatingCashFlow() (float64, error)
	CapitalExpenditure() (float64, error)
	Dividend() (float64, error)
	DividendPerShare() (float64, error)
	WAShares() (float64, error)
	Cash() (float64, error)
	Securities() (float64, error)
	Goodwill() (float64, error)
	Intangibles() (float64, error)
	Assets() (float64, error)
	Liabilities() (float64, error)
}

type valuation struct {
	Ticker    string            `json:"-"`
	Date      Timestamp         `json:"Recorded on"`
	FiledData []Measures        `json:"Measures"`
	Avgs      Average           `json:"Averages"`
	Pbm       PriceBasedMetrics `json:"Price Metrics"`
}

type valuator struct {
	collector  map[string]Collector
	Valuations map[string]*valuation `json:"Company"`
	store      Store
}

func (v valuator) String() string {
	return v.store.String()
}

func (v valuation) String() string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuation data: ", err)
	}
	return string(data)
}

func (v *valuator) Write() error {
	return v.store.write()
}

func (v *valuator) Store() error {
	for ticker, collector := range v.collector {
		by := bytes.NewBuffer(nil)
		if err := collector.Write(ticker, by); err == nil {
			v.store.putFinancials(ticker, by.Bytes())
		} else {
			log.Println("Error getting doc for ", ticker, err.Error())
		}
	}
	for ticker, value := range v.Valuations {
		v.store.putMeasures(ticker, []byte(value.String()))
	}
	return nil
}

func (v *valuator) Clean(ticker string) {
	delete(v.collector, ticker)
	delete(v.Valuations, ticker)

}

func (v *valuator) Filings(ticker string) []Filing {
	if v, ok := v.Valuations[ticker]; ok {
		filings := make([]Filing, len(v.FiledData))
		for i, m := range v.FiledData {
			filings[i] = m.Filing()
		}
		return filings
	}
	return []Filing{}
}

func (v *valuator) LastFiling(ticker string) Filing {
	if v, ok := v.Valuations[ticker]; ok {
		m := v.FiledData[len(v.FiledData)-1]
		return m.Filing()
	}
	return nil
}

func (v *valuator) Measures(ticker string) []Measures {
	if v, ok := v.Valuations[ticker]; ok {
		return v.FiledData
	}
	return []Measures{}
}

func (v *valuator) Averages(ticker string) Average {
	if v, ok := v.Valuations[ticker]; ok {
		return v.Avgs
	}
	return nil
}

func (v *valuator) PriceMetrics(ticker string) PriceBasedMetrics {
	if v, ok := v.Valuations[ticker]; ok {
		return v.Pbm
	}
	return nil
}

func (v *valuator) Collect(ticker string) error {
	if _, ok := v.collector[ticker]; ok {
		log.Println("Collection for ticker " + ticker + " is already done")
		return nil
	}
	v.Valuations[ticker] = &valuation{
		Ticker: ticker,
		Avgs:   nil,
	}
	v.store.read(ticker)
	collect, err := NewCollector(CollectorEdgar, v.store)
	if err != nil {
		log.Println("Error getting the collector: ", err.Error())
		v.Clean(ticker)
		return err
	}
	v.collector[ticker] = collect

	fils, err := collect.CollectAnnualData(ticker)
	if err != nil {
		log.Println("Error collecting annual data: ", err.Error())
		v.Clean(ticker)
		return err
	}
	mea := newMeasures(fils)

	if err = newYoYs(mea); err != nil {
		return err
	}

	avg, err := newAverages(mea)
	if err != nil {
		log.Println("Error collecting averages: ", err.Error())
		v.Clean(ticker)
		return err
	}
	valuation := v.Valuations[ticker]
	valuation.FiledData = mea
	valuation.Avgs = avg
	valuation.Pbm = newPriceBasedMetrics(mea[len(mea)-1])
	valuation.Date = Timestamp(time.Now())
	v.Store()

	return nil
}
