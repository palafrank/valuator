package main

import (
	"fmt"
	"net/http"

	"github.com/palafrank/valuator"
)

type valuatorStruct struct {
	Ticker   string
	Filings  []valuator.Filing
	Measures []valuator.Measures
	Avg      valuator.Average
}

func (s *server) queryForm(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		s.queryTempl.Execute(w, nil)
		return
	}
	s.queryTempl.Execute(w, nil)
	ticker := r.FormValue("ticker")
	save := len(r.Form["save"])
	if err := s.valuator.Collect(ticker); err != nil {
		w.Write([]byte(err.Error()))
		if save == 1 {
			s.valuator.Write()
		}
		return
	}
	s.writeData(w, ticker)
	//w.Write([]byte(s.valuator.HTML(ticker)))
}

func (s *server) writeData(w http.ResponseWriter, ticker string) {
	vals := valuatorStruct{}
	vals.Filings = s.valuator.Filings(ticker)
	vals.Measures = s.valuator.Measures(ticker)
	vals.Avg = s.valuator.Averages(ticker)
	vals.Ticker = ticker

	err := s.colTempl.Execute(w, vals)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = s.valTempl.Execute(w, vals)
	if err != nil {
		fmt.Println(err.Error())
	}
}
