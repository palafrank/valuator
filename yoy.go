package valuator

import (
	"encoding/json"
	"log"
	"strconv"
)

type Yoy interface {
	RevenueGrowth() float64
	EarningsGrowth() float64
	OlGrowth() float64
	GrossMarginGrowth() float64
	DebtGrowth() float64
	EquityGrowth() float64
	CashFlowGrowth() float64
	DividendGrowth() float64
	BookValueGrowth() float64
	HTML(string) string
}

type yoy struct {
	Revenue   float64 `json:"Revenue Growth"`
	Earnings  float64 `json:"Earning Growth"`
	Oleverage float64 `json:"Operating Leverage Growth"`
	Margin    float64 `json:"Gross Margin Growth"`
	Debt      float64 `json:"Debt Growth"`
	Equity    float64 `json:"Equity Growth"`
	Cf        float64 `json:"Cash Flow Growth"`
	Div       float64 `json:"Dividend Growth"`
	Bv        float64 `json:"Book Value Growth"`
}

func YoYHTMLHeader() string {

	trOpen := `<tr>`
	trClose := `</tr>`
	trDataOpen := `<th>`
	trDataClose := `</th>`

	trDoc := trOpen

	trDoc += trDataOpen + "Filed" + trDataClose
	trDoc += trDataOpen + "Revenue(%)" + trDataClose
	trDoc += trDataOpen + "Earnings(%)" + trDataClose
	trDoc += trDataOpen + "FCF(%)" + trDataClose
	trDoc += trDataOpen + "Margin(%)" + trDataClose
	trDoc += trDataOpen + "Debt(%)" + trDataClose
	trDoc += trDataOpen + "Equity(%)" + trDataClose
	trDoc += trDataOpen + "BV" + trDataClose
	trDoc += trDataOpen + "Div" + trDataClose
	trDoc += trClose

	return trDoc
}

func (y yoy) HTML(date string) string {
	trOpen := `<tr>`
	trClose := `</tr>`
	trDataOpen := `<th>`
	trDataClose := `</th>`

	// Start a new row
	trData := trOpen

	// Add all the Columns
	trData += trDataOpen
	trData += date
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.RevenueGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.EarningsGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.CashFlowGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.GrossMarginGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.DebtGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.EquityGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.BookValueGrowth(), 'f', 2, 64)
	trData += trDataClose

	trData += trDataOpen
	trData += strconv.FormatFloat(y.DividendGrowth(), 'f', 2, 64)
	trData += trDataClose

	// End the Row
	trData += trClose

	return trData

}

func (m yoy) String() string {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling year over year data: ", err)
	}
	return string(data)
}

func NewYoYs(mea []Measures) error {
	// Calculate YoY
	if len(mea) > 1 {
		for i := 1; i < len(mea); i++ {
			err := mea[i].NewYoy(mea[i-1])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewYoy(pastMeasure Measures, currentMeasure Measures) (*yoy, error) {

	ret := new(yoy)
	past := pastMeasure.Filing()
	current := currentMeasure.Filing()

	//Revenue Calculations
	p, err := past.Revenue()
	if err != nil {
		return nil, err
	}
	c, err := current.Revenue()
	if err != nil {
		return nil, err
	}
	ret.Revenue = yoyCalc(p, c, true)

	//Earnings Calculations
	p, err = past.NetIncome()
	if err != nil {
		return nil, err
	}
	c, err = current.NetIncome()
	if err != nil {
		return nil, err
	}
	ret.Earnings = yoyCalc(p, c, true)

	//Margin Calculations
	p, err = past.GrossMargin()
	if err != nil {
		return nil, err
	}
	c, err = current.GrossMargin()
	if err != nil {
		return nil, err
	}
	ret.Margin = yoyCalc(p, c, true)

	//OL calculations
	p = pastMeasure.OperatingLeverage()
	c = currentMeasure.OperatingLeverage()
	ret.Oleverage = yoyCalc(p, c, true)

	//Debt calculation.
	//Ignore error as LTD could be 0 or not collected
	p, _ = past.LongTermDebt()
	c, _ = current.LongTermDebt()
	ret.Debt = yoyCalc(p, c, true)

	//Equity calculation
	p, err = past.TotalEquity()
	if err != nil {
		return nil, err
	}
	c, err = current.TotalEquity()
	if err != nil {
		return nil, err
	}
	ret.Equity = yoyCalc(p, c, true)

	//Cash flow calculation
	p = pastMeasure.FreeCashFlow()
	c = currentMeasure.FreeCashFlow()
	ret.Cf = yoyCalc(p, c, true)

	p = pastMeasure.DividendPerShare()
	c = currentMeasure.DividendPerShare()
	ret.Div = yoyCalc(p, c, false)

	p = pastMeasure.BookValue()
	c = currentMeasure.BookValue()
	ret.Bv = yoyCalc(p, c, false)

	return ret, nil

}

func (y *yoy) RevenueGrowth() float64 {
	return y.Revenue
}

func (y *yoy) EarningsGrowth() float64 {
	return y.Earnings
}

func (y *yoy) OlGrowth() float64 {
	return y.Oleverage
}

func (y *yoy) GrossMarginGrowth() float64 {
	return y.Margin
}

func (y *yoy) DebtGrowth() float64 {
	return y.Debt
}

func (y *yoy) EquityGrowth() float64 {
	return y.Equity
}

func (y *yoy) CashFlowGrowth() float64 {
	return y.Cf
}

func (y *yoy) DividendGrowth() float64 {
	return y.Div
}

func (y *yoy) BookValueGrowth() float64 {
	return y.Bv
}
