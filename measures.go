package valuator

import (
	"fmt"
	"math"
)

type Measures interface {
	BookValue() float64
	OperatingLeverage() float64
	FinancialLeverage() float64
	ReturnOnEquity() float64
	ReturnOnAssets() float64
	DividendPerShare() float64
	FreeCashFlow() float64
	PayOutToFcf() float64
	String() string
}

type measures struct {
	filing Filing
}

func (m *measures) String() string {
	bv := fmt.Sprintf("%s%.2f\n", "Book Value: ", m.BookValue())
	ol := fmt.Sprintf("%s%.2f\n", "Operating Leverage: ", m.OperatingLeverage())
	fl := fmt.Sprintf("%s%.2f\n", "Financial Leverage: ", m.FinancialLeverage())
	return bv + ol + fl
}

func getMeasures(filing Filing) Measures {
	m := new(measures)
	m.filing = filing
	return m
}

/*
 BookValue:
    Value of the company retained within the equity portion of the BS
		BV = TotalEquity on balance sheet/Total share count

*/
func (m *measures) BookValue() float64 {
	eq, err := m.filing.TotalEquity()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	sc, err := m.filing.ShareCount()
	if err != nil {
		return 0
	}
	ret := math.Floor((eq/sc)*100) / 100
	return ret

}

/*
 Operating leverage:
     ratio of contribution margin to operating margin
 The ratio captures the relation between material cost of revenue vs the
 running cost of revenue
 contribution margin (CM) = Margin of profit against materials cost
 Operating margin (OM) = Margin of operating income against revenue
 Operating leverage = CM/OM
*/
func (m *measures) OperatingLeverage() float64 {
	oi, err := m.filing.OperatingIncome()
	if err != nil {
		return 0
	}

	rev, err := m.filing.Revenue()
	if err != nil {
		return 0
	}

	cr, err := m.filing.CostOfRevenue()
	if err != nil {
		return 0
	}

	cm := ((float64(rev) - float64(cr)) * 100) / float64(rev)
	om := (float64(oi) * 100) / float64(rev)

	ol := cm / om
	ret := math.Floor(ol*100) / 100
	return ret
}

func (m *measures) FinancialLeverage() float64 {
	eq, err := m.filing.TotalEquity()
	if err != nil {
		return 0
	}
	ld, err := m.filing.LongTermDebt()
	if err != nil {
		return 0
	}
	sd, err := m.filing.ShortTermDebt()
	if err != nil {
		return 0
	}
	ret := math.Floor(((ld+sd)/eq)*100) / 100
	return ret
}

func (m *measures) ReturnOnEquity() float64 {
	ni, err := m.filing.NetIncome()
	if err != nil {
		return 0
	}
	eq, err := m.filing.TotalEquity()
	if err != nil {
		return 0
	}
	return ni / eq
}

func (m *measures) ReturnOnAssets() float64 {
	return 0
}

func (m *measures) DividendPerShare() float64 {
	return 0
}

func (m *measures) FreeCashFlow() float64 {
	fcf, err := m.filing.OperatingCashFlow()
	if err != nil {
		return 0
	}
	return fcf
}

func (m *measures) PayOutToFcf() float64 {
	return 0
}
