package valuator

import (
	"log"
	"math"
)

type Valuator interface {
	/* Input:
	   ticker: ticker of the company
	   dr: Discount rate for DCF calculations
	   trend: % of the averages to factor in DCF calculations
	*/
	DiscountedCashFlow(ticker string, dr float64, trend float64, duration int) (float64, error)
	Save() error
	String() string
}

func NewValuator(ticker string) (Valuator, error) {
	v := &valuator{
		db:         NewDB(fileDBUrl, FileDatabaseType),
		collector:  make(map[string]Collector),
		Valuations: make(map[string]*valuation),
	}
	v.Valuations[ticker] = &valuation{
		Avgs: nil,
	}
	collect, err := NewCollector(collectorEdgar)
	if err != nil {
		return nil, err
	}
	v.collector[ticker] = collect
	mea, err := collect.CollectAnnualData(ticker)
	if err != nil {
		log.Println("Error collecting annual data: ", err.Error())
		return nil, err
	}
	avg, err := NewAverages(mea)
	if err != nil {
		log.Println("Error collecting averages: ", err.Error())
		return nil, err
	}

	v.Valuations[ticker].FiledData = mea
	v.Valuations[ticker].Avgs = avg

	return v, nil
}

func (v *valuator) DiscountedCashFlow(ticker string, dr float64, trend float64, duration int) (float64, error) {

	// First get the right parameters for the DCF calculations
	vals := v.Valuations[ticker]
	div := vals.Avgs.AvgDividendGrowth()
	bv := vals.Avgs.AvgBookValueGrowth()

	// Now adjust for trend
	div = div * (trend / 100)
	bv = bv * (trend / 100)

	// Start with the latest value
	outDiv := vals.FiledData[len(vals.FiledData)-1].DividendPerShare()
	outBv := vals.FiledData[len(vals.FiledData)-1].BookValue()

	sumDiv := outDiv
	sumBv := outBv / math.Pow(1+(dr/100), float64(duration))
	outBv = bv

	// Calculate discounted cash for each year
	for i := 1; i <= duration; i++ {

		// Add in the average growth
		outDiv += div
		outBv += bv

		// Discount it
		outDiv = outDiv / math.Pow(1+(dr/100), float64(i))
		outBv = outBv / math.Pow(1+(dr/100), float64(i))
		sumDiv += outDiv
		sumBv += outBv

	}

	return round(sumDiv + sumBv), nil
}
