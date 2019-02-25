package client

// Information for DCF trend query
type QueryDCFTrend struct {
	Ticker       string
	DiscountRate float32
	Trend        uint16
	Duration     uint8
	EndYear      uint32
}

// Infromation for Annual data
type QueryAnnualData struct {
	Ticker string
	Years  []int
}
