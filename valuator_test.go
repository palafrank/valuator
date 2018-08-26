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
