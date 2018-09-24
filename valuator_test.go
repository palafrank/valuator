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

func TestCollector(t *testing.T) {
	c, _ := NewCollector(collectorEdgar)
	m, _ := c.CollectAnnualData("AAPL", 2015, 2016, 2017)
	if m[2016].BookValue() != 24.05 {
		t.Error("Book value was not the expected value", m[2016].BookValue(), 24.05)
	}
	if m[2017].BookValue() != 26.10 {
		t.Error("Book value was not the expected value", m[2017].BookValue(), 26.10)
	}

	if m[2016].OperatingLeverage() != 1.40 {
		t.Error("Book value was not the expected value", m[2016].OperatingLeverage(), 1.40)
	}
	if m[2017].OperatingLeverage() != 1.43 {
		t.Error("Book value was not the expected value", m[2017].OperatingLeverage(), 1.43)
	}

	if m[2016].FinancialLeverage() != 0.61 {
		t.Error("Book value was not the expected value", m[2016].FinancialLeverage(), 0.61)
	}
	if m[2017].FinancialLeverage() != 0.77 {
		t.Error("Book value was not the expected value", m[2017].FinancialLeverage(), 0.77)
	}
	c.Save("AAPL")
}
