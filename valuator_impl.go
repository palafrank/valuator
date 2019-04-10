package valuator

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"reflect"
	"time"
)

var (
	valuatorDatabaseURL  = "./db/"
	valuatorDatabaseType = FileDatabaseType
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
}

type valuation struct {
	Date      Timestamp  `json:"Recorded on"`
	FiledData []Measures `json:"Measures"`
	Avgs      Average    `json:"Averages"`
	Price     float64    `json:"Market Price"`
	Mcap      float64    `json:"Market Capitalization"`
	EV        float64    `json:"Enterprise Value"`
}

type valuator struct {
	collector  map[string]Collector
	Valuations map[string]*valuation `json:"Company"`
	store      Store
}

func (v valuator) String() string {
	return v.store.String()
}

func (v valuator) HTML(ticker string) string {
	// Get filing data
	collector := v.collector[ticker].HTML(ticker)

	// Get Measures data
	by := bytes.NewBuffer(nil)
	t := template.New("valuator")
	t.Funcs(template.FuncMap{
		"isYoyNonNil": func(m Measures) bool {
			if yoy := m.Yoy(); !reflect.ValueOf(yoy).IsNil() {
				return true
			}
			return false
		},
		"getYoy": func(m Measures) Yoy {
			return m.Yoy()
		},
	})
	t, err := t.Parse(valuatorTemplate)
	if err != nil {
		panic(err.Error())
	}

	t.Execute(by, v.Valuations[ticker])
	measures := string(by.Bytes())

	return collector + measures
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

func (v *valuator) LastFiling(ticker string) Filing {
	if v, ok := v.Valuations[ticker]; ok {
		m := v.FiledData[len(v.FiledData)-1]
		return m.Filing()
	}
	return nil
}

func (v *valuator) Collect(ticker string) error {
	if _, ok := v.collector[ticker]; ok {
		return errors.New("Collection for ticker " + ticker + " is already done")
	}
	v.Valuations[ticker] = &valuation{
		Avgs: nil,
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
	valuation.Price = priceFetcher(ticker)
	valuation.EV = v.EnterpriseValue(ticker)
	valuation.Mcap = v.MarketCap(ticker)
	valuation.Date = Timestamp(time.Now())
	v.Store()

	return nil
}

func newValuatorDB() (Database, error) {
	return NewDatabase(valuatorDatabaseURL, valuatorDatabaseType)
}
