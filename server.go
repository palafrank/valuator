package valuator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func createAndRunServer() {
	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, from valuator")
}

func handleTickers(w http.ResponseWriter, r *http.Request) {

	var tickers, years []string

	dataOnly := false
	params := r.URL.Query()
	for key, data := range params {
		if key == "ticker" {
			// List of tickers for valuation
			tickers = data
		} else if key == "years" {
			years = data
		} else if key == "test" {
			fmt.Fprintf(w, "Hello, from valuator ticker handler "+r.URL.String())
			return
		} else if key == "data" {
			dataOnly = true
		}
	}
	if len(tickers) <= 0 {
		fmt.Fprintf(w, "No valid ticker found in query string")
		return
	}
	ys, err1 := validateYears(years)
	err2 := validateTickers(tickers)
	if err1 != nil {
		fmt.Fprintf(w, err1.Error())
	} else if err2 != nil {
		fmt.Fprintf(w, err2.Error())
	} else {
		handleWebQuery(w, tickers, ys, dataOnly)
	}
}

func registerHandlers() {
	http.HandleFunc("/valuator", handleTickers)
}

func StartServer() {
	registerHandlers()
	createAndRunServer()
}

func handleWebQuery(w http.ResponseWriter, tickers []string, years []int, dataOnly bool) {
	if dataOnly {
		db, _ := NewValuatorDB()
		store := NewStore(db)
		c, _ := NewCollector(collectorEdgar, store)
		for _, tick := range tickers {
			_, err := c.CollectAnnualData(tick, years...)
			if err != nil {
				fmt.Fprintln(w, "Failed to collect data for ", tick, err.Error())
				return
			}
			if err := c.Write(tick, w); err != nil {
				fmt.Fprintln(w, "Failed to collect data for ", tick, err.Error())
			}
		}
		return
	}
	v, err := NewValuator()
	if err != nil {
		fmt.Fprintln(w, "Failed to create a valuator")
		return
	}
	for _, tick := range tickers {
		err = v.Collect(tick)
		if err != nil {
			fmt.Fprintln(w, "Failed to collect data for ", tick, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(v.String()))
}

func validateYears(dates []string) ([]int, error) {
	var ret []int
	for _, d := range dates {
		date, err := strconv.Atoi(d)
		if err != nil {
			return ret, errors.New("Invalid year specified in query " + d)
		}
		if date >= 2012 && date > time.Now().Year() {
			return ret, errors.New("Invalid year specified: " + d)
		}
		ret = append(ret, date)
	}
	return ret, nil
}

func validateTickers(tickers []string) error {
	return nil
}
