package valuator

import "log"

// Valuator interface allows queries into the the valuator for valuation metrics
type Valuator interface {
	/*
		 	DiscountedCashFlowTrend
			Calculated DCF based on the collected BV and DIV growth rates
			trend will adjust the growth rates based on user Input
			    - 100 will indicate 100% trend
					- 50 will indicate 50% of trend
					- 150 will indicate 150% trend
			Input:
		   ticker: ticker of the company
		   dr: Discount rate for DCF calculations
		   trend: % of the averages to factor in DCF calculations
			 duration: Time over which to discount the CF
	*/
	DiscountedCashFlowTrend(ticker string, dr float64, trend float64, duration int, endYear ...int) (float64, error)
	/*
		 	DiscountedCashFlow
			Input:
		   ticker: ticker of the company
		   dr: Discount rate for DCF calculations
		   bvIn: Book value rate of change
			 divIn: Dividend rate of change
			 duration: Time over which to discount the CF
	*/
	DiscountedCashFlow(ticker string, dr float64, bvIn float64, divIn float64, duration int, endYear ...int) (float64, error)
	/*
		 	DiscountedFCFTrend
			Calculated DCF based on the collected FCF and DIV growth rates
			FCF will be based on cash flow from opertaions adjusted for capex
			trend will adjust the growth rates based on user Input
			    - 100 will indicate 100% trend
					- 50 will indicate 50% of trend
					- 150 will indicate 150% trend
			Input:
		   ticker: ticker of the company
		   dr: Discount rate for DCF calculations
		   trend: % of the averages to factor in DCF calculations
			 duration: Time over which to discount the CF
	*/
	DiscountedFCFTrend(ticker string, dr float64, trend float64, duration int, endYear ...int) (float64, error)

	/* Utility APIs*/
	Collect(string) error
	Clean(string)
	Write() error
	String() string
	HTML(string) string
}

// NewValuator creates a new valuator to generate valuation metrics
func NewValuator() (Valuator, error) {
	db, err := newValuatorDB()
	if err != nil {
		log.Println("Error creating valuator: ", err)
	}
	v := &valuator{
		collector:  make(map[string]Collector),
		Valuations: make(map[string]*valuation),
		store:      newStore(db),
	}

	return v, nil
}
