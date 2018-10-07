package valuator

import (
	"encoding/json"
	"log"
)

type Yoy interface {
	RevenueGrowth() float64
	EarningsGrowth() float64
	OlGrowth() float64
	GrossMarginGrowth() float64
	DebtGrowth() float64
	EquityGrowth() float64
	CashFlowGrowth() float64
}

type yoy struct {
	Revenue   float64 `json:"Revenue Growth"`
	Earnings  float64 `json:"Earning Growth"`
	Oleverage float64 `json:"Operating Leverage Growth"`
	Margin    float64 `json:"Gross Margin Growth"`
	Debt      float64 `json:"Debt Growth"`
	Equity    float64 `json:"Equity Growth"`
	Cf        float64 `json:"Cash Flow Growth"`
}

func (m yoy) String() string {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling year over year data: ", err)
	}
	return string(data)
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
	ret.Revenue = yoyCalc(p, c)

	//Earnings Calculations
	p, err = past.NetIncome()
	if err != nil {
		return nil, err
	}
	c, err = current.NetIncome()
	if err != nil {
		return nil, err
	}
	ret.Earnings = yoyCalc(p, c)

	//Margin Calculations
	p, err = past.GrossMargin()
	if err != nil {
		return nil, err
	}
	c, err = current.GrossMargin()
	if err != nil {
		return nil, err
	}
	ret.Margin = yoyCalc(p, c)

	//OL calculations
	p = pastMeasure.OperatingLeverage()
	c = currentMeasure.OperatingLeverage()
	ret.Oleverage = yoyCalc(p, c)

	//Debt calculation
	p, err = past.LongTermDebt()
	if err != nil {
		return nil, err
	}
	c, err = current.LongTermDebt()
	if err != nil {
		return nil, err
	}
	ret.Debt = yoyCalc(p, c)

	//Equity calculation
	p, err = past.TotalEquity()
	if err != nil {
		return nil, err
	}
	c, err = current.TotalEquity()
	if err != nil {
		return nil, err
	}
	ret.Equity = yoyCalc(p, c)

	//Cash flow calculation
	p = pastMeasure.FreeCashFlow()
	c = currentMeasure.FreeCashFlow()
	ret.Cf = yoyCalc(p, c)

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
