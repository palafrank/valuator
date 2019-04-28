package valuator

// Valuator interface allows queries into the the valuator for valuation metrics
type Valuator interface {

	// Interfaces to calculate valuation of a Company

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

	/* Per Ticker interface */

	// Collect collects filing information for a specific ticker
	Collect(string) error

	// Filings returns all the collected filings for a specific ticker
	Filings(string) []Filing

	// LastFiling returns the cronologically ordered last filing collected
	LastFiling(string) Filing

	// Measures returns all the computed measures for a specific ticker
	Measures(string) []Measures

	// Averages returns the computed average metrics for a specific ticker
	Averages(string) Average

	// PriceMetrics gets the price based metrics computed based on current price
	PriceMetrics(string) PriceBasedMetrics

	// Clean clears all the filing data collected for a specific ticker
	Clean(string)

	/* Overall Valuator interface */

	// Write saves the entire data in the valuator to the underlying database
	Write() error

	// String creates a JSON output of the entire data in the valuator
	String() string
}

// NewValuator creates a new valuator to generate valuation metrics
func NewValuator(db Database) (Valuator, error) {
	if db == nil {
		var err error
		if db, err = NewDatabase(nil, NoneDatabaseType); err != nil {
			panic("Could not create a NoneDatabaseType")
		}
	}
	v := &valuator{
		collector:  make(map[string]Collector),
		Valuations: make(map[string]*valuation),
		store:      newStore(db),
	}

	return v, nil
}
