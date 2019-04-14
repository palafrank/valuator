package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/palafrank/valuator"
)

var (
	valuatorDatabaseURL  = "./db_test/"
	valuatorDatabaseType = valuator.FileDatabaseType
	testFileDB, _        = valuator.NewDatabase(valuatorDatabaseURL, valuatorDatabaseType)
)

func TestValuatorQuery(t *testing.T) {

	s, err := newServer(valuatorDatabaseURL, valuatorDatabaseType)
	if err != nil {
		panic("Failed to create server")
	}
	ts := httptest.NewServer(http.HandlerFunc(s.queryForm))

	client := ts.Client()

	// Check HTML format
	form := make(url.Values)
	form.Add("ticker", "IBM")
	res, err := client.PostForm(ts.URL, form)
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
