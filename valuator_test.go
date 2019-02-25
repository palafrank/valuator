package valuator

import (
	"testing"
)

func TestAAPLCollector(t *testing.T) {
	db, _ := NewValuatorDB()
	s := NewStore(db)
	c, _ := NewCollector(collectorEdgar, s)
	fs, _ := c.CollectAnnualData("AAPL", 2012, 2013, 2014, 2015, 2016, 2017)
	m := NewMeasures(fs)
	if m[4].BookValue() != 24.05 {
		t.Error("Book value was not the expected value", m[4].BookValue(), 24.05)
	}
	if m[4].WorkingCapital() != 27863000000 {
		t.Error("Working Capital was not the expected value", m[4].WorkingCapital(), 27863000000)
	}
	if m[4].CurrentRatio() != 1.35 {
		t.Error("Current Ratio was not the expected value", m[4].CurrentRatio(), 1.35)
	}
	if m[5].BookValue() != 26.10 {
		t.Error("Book value was not the expected value", m[5].BookValue(), 26.10)
	}
	if m[5].WorkingCapital() != 27831000000 {
		t.Error("Working Capital was not the expected value", m[5].WorkingCapital(), 27831000000)
	}
	if m[5].CurrentRatio() != 1.27 {
		t.Error("Current Ratio was not the expected value", m[5].CurrentRatio(), 1.27)
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
	s.Write()
}

func TestNewPSXValuator(t *testing.T) {
	v, err := NewValuator()
	if err != nil {
		t.Error("Failed to create valuator: ", err.Error())
		return
	}
	err = v.Collect("PSX")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
	}

	ret, _ := v.DiscountedCashFlowTrend("PSX", 3, 100, 10, 2018)
	if ret != 192.95 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedCashFlowTrend("PSX", 3, 50, 10, 2018)
	if ret != 124.36 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
	ret, _ = v.DiscountedCashFlowTrend("PSX", 3, 20, 10, 2018)
	if ret != 83.2 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	v.Write()

}

func TestNewTGTValuator(t *testing.T) {
	v, err := NewValuator()
	if err != nil {
		t.Error("Failed to create valuator: ", err.Error())
		return
	}
	err = v.Collect("TGT")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}

	ret, _ := v.DiscountedCashFlowTrend("TGT", 3, 100, 10, 2018)
	if ret != 27.72 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedCashFlow("TGT", 3, 0.30, 0.23, 10, 2018)
	if ret != 51.51 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	_, err = v.DiscountedCashFlow("PSX", 3, 0.30, 0.23, 10, 2018)
	if err == nil {
		t.Error("Error: Call should have failed as data has not been collected")
	}
	_, err = v.DiscountedCashFlowTrend("PSX", 3, 100, 10, 2018)
	if err == nil {
		t.Error("Error: Call should have failed as data has not been collected")
	}
	if err = v.Write(); err != nil {
		t.Error("Error: Failed to save valuations")
	}
}

func TestNewIBMValuator(t *testing.T) {
	v, err := NewValuator()
	if err != nil {
		t.Error("Failed to create valuator: ", err.Error())
		return
	}
	err = v.Collect("IBM")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}

	ret, _ := v.DiscountedCashFlowTrend("IBM", 3, 100, 10, 2018)
	if ret != 81.15 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
	if err = v.Write(); err != nil {
		t.Error("Error: Failed to save valuations")
	}
}

func TestNewMSFTValuator(t *testing.T) {
	v, err := NewValuator()
	if err != nil {
		t.Error("Failed to create valuator: ", err.Error())
		return
	}
	err = v.Collect("MSFT")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}

	ret, _ := v.DiscountedCashFlowTrend("MSFT", 3, 100, 10, 2018)
	if ret != 41.5 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedCashFlow("MSFT", 3, 1.36, 0.29, 10, 2018)
	if ret != 79.39 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}
}

func TestNewCSCOValuator(t *testing.T) {
	v, err := NewValuator()
	if err != nil {
		t.Error("Failed to create valuator: ", err.Error())
		return
	}
	err = v.Collect("CSCO")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}

	ret, _ := v.DiscountedCashFlowTrend("CSCO", 3, 100, 10, 2018)
	if ret != 18.28 {
		t.Error("Error in DCF calculation at 100% trend ", ret)
	}

	ret, _ = v.DiscountedFCFTrend("CSCO", 3, 100, 10, 2018)
	if ret != 29.51 {
		t.Error("Error in DCFCF calculation at 100% trend ", ret)
	}
	if err = v.Write(); err != nil {
		t.Error("Error: Failed to save valuations")
	}

}

func TestNewWValuator(t *testing.T) {
	v, err := NewValuator()
	if err != nil {
		t.Error("Failed to create valuator: ", err.Error())
		return
	}
	err = v.Collect("W")
	if err != nil {
		t.Error("Failed to create a valuator: ", err.Error())
		return
	}

}

func TestValuatorStore(t *testing.T) {
	db, _ := NewValuatorDB()
	s := NewStore(db)
	s.Read("IBM")
	collect, err := NewCollector(collectorEdgar, s)
	if err != nil {
		t.Error("Failed to create a collector with a store")
		return
	}
	m, _ := collect.CollectAnnualData("IBM")
	if len(m) == 0 {
		t.Error("Failed to get measures from data read from database")
	}

}
