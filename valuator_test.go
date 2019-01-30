package valuator

import (
	"testing"

	"github.com/palafrank/edgar"
)

func TestGetEdgarFiling(t *testing.T) {
	fetcher, _ := NewEdgarCollector()

	folder, _ := fetcher.CompanyFolder("AAPL", edgar.FilingType10K)
	if folder.Ticker() != "AAPL" {
		t.Error("Failed to get the ticker from the fetcher")
	}
	a := folder.AvailableFilings(edgar.FilingType10K)
	if len(a) < 7 {
		t.Error("Failed to get list of filings")
		return
	}
	f, err := folder.Filing(edgar.FilingType10K, a[0])
	if err != nil || f == nil {
		t.Error(err.Error())
	}

}

func TestAAPLCollector(t *testing.T) {
	c, _ := NewCollector(collectorEdgar)
	m, _ := c.CollectAnnualData("AAPL", 2015, 2016, 2017)
	if m[1].BookValue() != 24.05 {
		t.Error("Book value was not the expected value", m[1].BookValue(), 24.05)
	}
	if m[2].BookValue() != 26.10 {
		t.Error("Book value was not the expected value", m[2].BookValue(), 26.10)
	}

	if m[1].OperatingLeverage() != 1.44 {
		t.Error("Operating Leverage was not the expected value", m[1].OperatingLeverage(), 1.44)
	}
	if m[2].OperatingLeverage() != 1.46 {
		t.Error("Operating Leverage was not the expected value", m[2].OperatingLeverage(), 1.46)
	}

	if m[1].FinancialLeverage() != 58 {
		t.Error("Financial Leverage was not the expected value", m[1].FinancialLeverage(), 0.61)
	}
	if m[2].FinancialLeverage() != 72 {
		t.Error("Financial Leverage was not the expected value", m[2].FinancialLeverage(), 0.77)
	}
	c.Save("AAPL")
}

func TestPSXCollector(t *testing.T) {
	c, _ := NewCollector(collectorEdgar)
	c.CollectAnnualData("PSX")
	c.Save("PSX")
}

func TestNewValuator(t *testing.T) {
	v, err := NewValuator("PSX")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
	}
	v.Save()
}
