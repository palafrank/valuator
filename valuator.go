package valuator

// Filing interface for fetching financial data
type Filing interface {
	Ticker() string
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
}

type valuator struct {
	measures map[int]Measures
}

func NewValuator(ticker string) error {
	v := new(valuator)
	collect, err := NewCollector("edgar")
	if err != nil {
		return err
	}
	measures, err := collect.CollectAnnualData(ticker)
	if err != nil {
		return err
	}
	v.measures = measures
	return nil
}
