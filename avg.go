package valuator

import (
	"errors"
	"log"
	"reflect"
)

type Average interface {
	AvgRevenueGrowth() float64
	AvgEarningsGrowth() float64
	AvgOlGrowth() float64
	AvgGrossMarginGrowth() float64
	AvgDebtGrowth() float64
	AvgEquityGrowth() float64
	AvgCashFlowGrowth() float64
}

type averages struct {
	Revenue   float64 `json:"Average Revenue Growth"`
	Earnings  float64 `json:"Average Earning Growth"`
	Oleverage float64 `json:"Average Operating Leverage Growth"`
	Margin    float64 `json:"Average Gross Margin Growth"`
	Debt      float64 `json:"Average Debt Growth"`
	Equity    float64 `json:"Average Equity Growth"`
	Cf        float64 `json:"Average Cash Flow Growth"`
}

func NewAverages(m []Measures) (Average, error) {
	var y []Yoy

	for _, val := range m {
		lvar := val.Yoy()
		if !reflect.ValueOf(lvar).IsNil() {
			y = append(y, lvar)
		}
	}

	if len(y) == 0 {
		log.Println("No year over year calculations found. Skipping averages")
		return nil, errors.New("No YoY information found. Skipping averages")
	}

	avg := new(averages)
	for _, val := range y {
		avg.Revenue = avg.Revenue + val.RevenueGrowth()
		avg.Earnings = avg.Earnings + val.EarningsGrowth()
		avg.Oleverage = avg.Oleverage + val.OlGrowth()
		avg.Margin = avg.Margin + val.GrossMarginGrowth()
		avg.Debt = avg.Debt + val.DebtGrowth()
		avg.Equity = avg.Equity + val.EquityGrowth()
		avg.Cf = avg.Cf + val.CashFlowGrowth()
	}
	avg.Revenue = avgCalc(avg.Revenue, float64(len(y)))
	avg.Earnings = avgCalc(avg.Earnings, float64(len(y)))
	avg.Oleverage = avgCalc(avg.Oleverage, float64(len(y)))
	avg.Margin = avgCalc(avg.Margin, float64(len(y)))
	avg.Debt = avgCalc(avg.Debt, float64(len(y)))
	avg.Equity = avgCalc(avg.Equity, float64(len(y)))
	avg.Cf = avgCalc(avg.Cf, float64(len(y)))

	return avg, nil
}

func (a *averages) AvgRevenueGrowth() float64 {
	return a.Revenue
}

func (a *averages) AvgEarningsGrowth() float64 {
	return a.Earnings
}

func (a *averages) AvgOlGrowth() float64 {
	return a.Oleverage
}

func (a *averages) AvgGrossMarginGrowth() float64 {
	return a.Margin
}

func (a *averages) AvgDebtGrowth() float64 {
	return a.Debt
}

func (a *averages) AvgEquityGrowth() float64 {
	return a.Equity
}

func (a *averages) AvgCashFlowGrowth() float64 {
	return a.Cf
}
