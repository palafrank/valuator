package valuator

// PriceBasedMetrics provides an interface for price based stock metrics
type PriceBasedMetrics interface {
	EnterpriseValue() float64
	MarketCapitalization() float64
	PriceOverEarnings() float64
	PriceOverCashFlow() float64
	PriceOverRevenue() float64
}

type pbm struct {
	measures  Measures
	Price     float64 `json:"Market Price"`
	Ev        float64 `json:"Enterprise Value"`
	MarketCap float64 `json:"Market Capitalization"`
	PoverE    float64 `json:"Price To Earnings"`
	PoverCF   float64 `json:"Price To CashFlow"`
	PoverRev  float64 `json:"Price To Revenue"`
}

func newPriceBasedMetrics(m Measures) PriceBasedMetrics {
	pm := &pbm{
		measures: m,
		Price:    priceFetcher(m.Filing().Ticker()),
	}
	// Enterprise Value
	cash, _ := m.Filing().Cash()
	ld, _ := m.Filing().LongTermDebt()
	sd, _ := m.Filing().ShortTermDebt()
	sc, _ := m.Filing().ShareCount()
	pm.MarketCap = round(sc * pm.Price)
	pm.Ev = round(pm.MarketCap + ld + sd - cash)

	// Price over Earnings
	if ni, err := m.Filing().NetIncome(); err == nil {
		pm.PoverE = round(pm.MarketCap / ni)
	}

	// Price over Cash flow
	if cf, err := m.Filing().OperatingCashFlow(); err == nil {
		pm.PoverCF = round(pm.MarketCap / cf)
	}

	if rev, err := m.Filing().Revenue(); err == nil {
		pm.PoverRev = round(pm.MarketCap / rev)
	}

	return pm
}

func (p *pbm) collect() {

}

func (p *pbm) EnterpriseValue() float64 {
	return p.Ev
}

func (p *pbm) PriceOverEarnings() float64 {
	return p.PoverE
}

func (p *pbm) PriceOverCashFlow() float64 {
	return p.PoverCF
}

func (p *pbm) MarketCapitalization() float64 {
	return p.MarketCap
}

func (p *pbm) PriceOverRevenue() float64 {
	return p.PoverRev
}
