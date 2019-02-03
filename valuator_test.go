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
	m, _ := c.CollectAnnualData("AAPL", 2012, 2013, 2014, 2015, 2016, 2017)
	if m[4].BookValue() != 24.05 {
		t.Error("Book value was not the expected value", m[4].BookValue(), 24.05)
	}
	if m[5].BookValue() != 26.10 {
		t.Error("Book value was not the expected value", m[5].BookValue(), 26.10)
	}

	if m[4].OperatingLeverage() != 1.44 {
		t.Error("Operating Leverage was not the expected value", m[4].OperatingLeverage(), 1.44)
	}
	if m[5].OperatingLeverage() != 1.46 {
		t.Error("Operating Leverage was not the expected value", m[5].OperatingLeverage(), 1.46)
	}

	if m[4].FinancialLeverage() != 58 {
		t.Error("Financial Leverage was not the expected value", m[4].FinancialLeverage(), 0.61)
	}
	if m[5].FinancialLeverage() != 72 {
		t.Error("Financial Leverage was not the expected value", m[5].FinancialLeverage(), 0.77)
	}
	c.Save("AAPL")
}

func TestNewPSXValuator(t *testing.T) {
	v, err := NewValuator("PSX")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
	}
	v.Save()
	ret, _ := v.DiscountedCashFlowTrend("PSX", 3, 100, 10)
	if ret != 192.95 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
	ret, _ = v.DiscountedCashFlowTrend("PSX", 3, 50, 10)
	if ret != 124.36 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
	ret, _ = v.DiscountedCashFlowTrend("PSX", 3, 20, 10)
	if ret != 83.2 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
}

func TestNewTGTValuator(t *testing.T) {
	v, err := NewValuator("TGT")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}
	v.Save()

	ret, _ := v.DiscountedCashFlowTrend("TGT", 3, 100, 10)
	if ret != 27.72 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedCashFlow("TGT", 3, 0.30, 0.23, 10)
	if ret != 51.51 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	_, err = v.DiscountedCashFlow("PSX", 3, 0.30, 0.23, 10)
	if err == nil {
		t.Error("Error: Call should have failed as data has not been collected")
	}
	_, err = v.DiscountedCashFlowTrend("PSX", 3, 100, 10)
	if err == nil {
		t.Error("Error: Call should have failed as data has not been collected")
	}

}

func TestNewIBMValuator(t *testing.T) {
	v, err := NewValuator("IBM")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}
	v.Save()

	ret, _ := v.DiscountedCashFlowTrend("IBM", 3, 100, 10)
	if ret != 81.15 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
}

func TestNewMSFTValuator(t *testing.T) {
	v, err := NewValuator("MSFT")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}
	v.Save()

	ret, _ := v.DiscountedCashFlowTrend("MSFT", 3, 100, 10)
	if ret != 41.5 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedCashFlow("MSFT", 3, 1.36, 0.29, 10)
	if ret != 79.39 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
}

func TestNewCSCOValuator(t *testing.T) {
	v, err := NewValuator("CSCO")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}

	ret, _ := v.DiscountedCashFlowTrend("CSCO", 3, 100, 10)
	if ret != 18.28 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedFCFTrend("CSCO", 3, 100, 10)
	if ret != 29.51 {
		t.Error("Error in DCFCF calculation at 100% trend ", ret)
	}
	v.Save()
}
