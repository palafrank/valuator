package valuator

import (
	"errors"
	"math"
)

func (v *valuator) DiscountedCashFlowTrend(ticker string, dr float64, trend float64, duration int, endYear ...int) (float64, error) {

	if len(endYear) > 1 {
		return 0, errors.New("Specify only one end year for DCF calculation")
	}

	// First get the right parameters for the DCF calculations
	vals, ok := v.Valuations[ticker]
	if !ok {
		return 0, errors.New("Valuator has not be told to collect data on " + ticker)
	}

	div := vals.Avgs.AvgDividendGrowth()
	bv := vals.Avgs.AvgBookValueGrowth()

	if len(endYear) == 1 {
		meas := createMeasuresList(vals.FiledData, endYear[0])
		if avgs, err := NewAverages(meas); err == nil {
			div = avgs.AvgDividendGrowth()
			bv = avgs.AvgBookValueGrowth()
		} else {
			return 0, err
		}
	}

	// Now adjust for trend
	div = div * (trend / 100)
	bv = bv * (trend / 100)

	return v.DiscountedCashFlow(ticker, dr, bv, div, duration, endYear...)
}

func (v *valuator) DiscountedFCFTrend(ticker string, dr float64, trend float64, duration int, endYear ...int) (float64, error) {

	if len(endYear) > 1 {
		return 0, errors.New("Specify only one end year for DCF calculation")
	}
	// First get the right parameters for the DCF calculations
	vals, ok := v.Valuations[ticker]
	if !ok {
		return 0, errors.New("Valuator has not be told to collect data on " + ticker)
	}
	div := vals.Avgs.AvgDividendGrowth()
	fcf := vals.Avgs.AvgCashFlowGrowth()
	// Get the latest book value
	bv := vals.FiledData[len(vals.FiledData)-1].BookValue()

	if len(endYear) == 1 {
		meas := createMeasuresList(vals.FiledData, endYear[0])
		if avgs, err := NewAverages(meas); err == nil {
			div = avgs.AvgDividendGrowth()
			fcf = avgs.AvgCashFlowGrowth()
			bv = meas[len(meas)-1].BookValue()
		} else {
			return 0, err
		}
	}

	// Now adjust for trend
	div = div * (trend / 100)
	fcf = fcf * (trend / 100)

	// Get the BV growth at the rate of FCF growth
	bv = bv * (fcf / 100)

	return v.DiscountedCashFlow(ticker, dr, bv, div, duration, endYear...)
}

func (v *valuator) DiscountedCashFlow(ticker string, dr float64, bvIn float64, divIn float64, duration int, endYear ...int) (float64, error) {

	if len(endYear) > 1 {
		return 0, errors.New("Specify only one end year for DCF calculation")
	}
	// First get the right parameters for the DCF calculations
	vals, ok := v.Valuations[ticker]
	if !ok {
		return 0, errors.New("Valuator has not be told to collect data on " + ticker)
	}
	div := divIn
	bv := bvIn

	// Start with the latest value
	outDiv := vals.FiledData[len(vals.FiledData)-1].DividendPerShare()
	outBv := vals.FiledData[len(vals.FiledData)-1].BookValue()

	if len(endYear) == 1 {
		meas := createMeasuresList(vals.FiledData, endYear[0])
		outDiv = meas[len(meas)-1].DividendPerShare()
		outBv = meas[len(meas)-1].BookValue()
	}

	sumDiv := outDiv
	sumBv := outBv / math.Pow(1+(dr/100), float64(duration))
	outBv = bv

	// Calculate discounted cash for each year
	for i := 1; i <= duration; i++ {

		// Add in the average growth
		outDiv += div
		outBv += bv

		// Discount it
		outDiv = outDiv / math.Pow(1+(dr/100), float64(i))
		outBv = outBv / math.Pow(1+(dr/100), float64(i))
		sumDiv += outDiv
		sumBv += outBv

	}

	return round(sumDiv + sumBv), nil
}
