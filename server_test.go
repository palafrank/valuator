package valuator

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValuatorQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleTickers))
	client := ts.Client()

	res, err := client.Get(ts.URL + "?ticker=AAPL&format=json")
	if err != nil {
		t.Error("Failed to get response from valuator server ", err.Error())
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Failed to parse greetings from valuator server ", err.Error())
		return
	}

	if !strings.Contains(string(data), `"Company": "AAPL"`) {
		t.Error("Failed to find company name")
		return
	}

	if !strings.Contains(string(data), `"Financial Data"`) {
		t.Error("Financial data not included")
		return
	}

	if !strings.Contains(string(data), `"Financial Measures"`) {
		t.Error("Financial measures not included")
		return
	}

	if !strings.Contains(string(data), `"Averages"`) {
		t.Error("Averages not included")
		return
	}

	if !strings.Contains(string(data), `"YoY"`) {
		t.Error("YoY not included")
		return
	}

	if !strings.Contains(string(data), `"Report date": "2012-10-31"`) {
		t.Error("2012 report not found")
		return
	}

	if !strings.Contains(string(data), `"Report date": "2018-11-05"`) {
		t.Error("2018 report not found")
		return
	}

	ts.Close()

}

func TestValuatorHTMLQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleTickers))
	client := ts.Client()

	// Check HTML format
	res, err := client.Get(ts.URL + "?ticker=IBM")
	if err != nil {
		t.Error("Failed to get response from valuator server ", err.Error())
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Failed to parse greetings from valuator server ", err.Error())
		return
	}

	if strings.Contains(string(data), `"Financial Measures"`) {
		t.Error("Should not be JSON format")
		return
	}

	if !strings.Contains(string(data), `<tr>`) {
		t.Error("Should contain HTML rows")
		return
	}
	if !strings.Contains(string(data), `<table`) {
		t.Error("Should contain HTML table")
		return
	}

}

func TestValuatorDataJSONQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleTickers))
	client := ts.Client()

	// Check JSON format
	res, err := client.Get(ts.URL + "?ticker=IBM&data=yes&format=json")
	if err != nil {
		t.Error("Failed to get response from valuator server ", err.Error())
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Failed to parse greetings from valuator server ", err.Error())
		return
	}

	if strings.Contains(string(data), `"Financial Measures"`) {
		t.Error("Financial measures included")
		return
	}

	if strings.Contains(string(data), `"Averages"`) {
		t.Error("Averages included")
		return
	}

	if strings.Contains(string(data), `"YoY"`) {
		t.Error("YoY included")
		return
	}

	if !strings.Contains(string(data), `"Financial Data"`) {
		t.Error("Filing data not included")
		return
	}

	if !strings.Contains(string(data), `"Operational Information"`) {
		t.Error("Operational Information not included")
		return
	}

	if !strings.Contains(string(data), `"Balance Sheet Information"`) {
		t.Error("Balance Sheet information not included")
		return
	}

	if !strings.Contains(string(data), `"Cash Flow Information"`) {
		t.Error("Cash flow information not included")
		return
	}

	if !strings.Contains(string(data), `"Entity Information"`) {
		t.Error("Entity information not included")
		return
	}

	ts.Close()

}

func TestValuatorQuerySpecific(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleTickers))
	client := ts.Client()

	res, err := client.Get(ts.URL + "?ticker=AAPL&data=yes&years=2013&years=2014&format=json")
	if err != nil {
		t.Error("Failed to get response from valuator server ", err.Error())
		return
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error("Failed to parse greetings from valuator server ", err.Error())
		return
	}

	if !strings.Contains(string(data), `"Report date": "2013-10-30"`) {
		t.Error("Returned data doesn't contain report date")
		return
	}
	if !strings.Contains(string(data), `"Report date": "2014-10-27"`) {
		t.Error("Returned data doesn't contain report date")
		return
	}

	if !strings.Contains(string(data), `"Shares Outstanding": 899738000`) {
		t.Error("Returned data doesn't contain report date")
		return
	}

	if !strings.Contains(string(data), `"Shares Outstanding": 5864840000`) {
		t.Error("Returned data doesn't contain report date")
		return
	}

	ts.Close()

}
