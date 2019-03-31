package valuator

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"
)

var (
	valuatorDatabaseURL  string       = "./db/"
	valuatorDatabaseType DatabaseType = FileDatabaseType
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
	collector := v.collector[ticker].HTML(ticker)
	measures := v.Valuations[ticker].HTML()
	return collector + measures
}

func (v valuation) String() string {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuation data: ", err)
	}
	return string(data)
}

func (v valuation) HTML() string {
	header := `<html><body><table border="1">`
	footer := `</table></body></html>`

	title := `<h2>Valuation Metrics</h2>`
	trDoc := header + title
	trDoc += MeasuresHTMLHeader()
	for _, m := range v.FiledData {
		trDoc += m.HTML()
	}
	trDoc += footer

	title = `<h2>YoY Metrics</h2>`
	trDoc += header + title
	trDoc += YoYHTMLHeader()
	for _, m := range v.FiledData {
		if yoy := m.Yoy(); !reflect.ValueOf(yoy).IsNil() {
			trDoc += yoy.HTML(m.FiledOn())
		}
	}

	trDoc += footer

	return trDoc
}

func (v *valuator) Write() error {
	return v.store.Write()
}

func (v *valuator) Store() error {
	for ticker, collector := range v.collector {
		by := bytes.NewBuffer(nil)
		if err := collector.Write(ticker, by); err == nil {
			v.store.PutFinancials(ticker, by.Bytes())
		} else {
			log.Println("Error getting doc for ", ticker, err.Error())
		}
	}
	for ticker, value := range v.Valuations {
		v.store.PutMeasures(ticker, []byte(value.String()))
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
	v.store.Read(ticker)
	collect, err := NewCollector(collectorEdgar, v.store)
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
	mea := NewMeasures(fils)

	if err := NewYoYs(mea); err != nil {
		return err
	}

	avg, err := NewAverages(mea)
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

func NewValuatorDB() (database, error) {
	return NewDatabase(valuatorDatabaseURL, valuatorDatabaseType)
}
