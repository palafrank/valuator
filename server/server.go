package main

import (
	"errors"
	"html/template"
	"net/http"
	"path"
	"reflect"
	"time"

	"github.com/palafrank/valuator"
)

type server struct {
	valuator   valuator.Valuator
	queryTempl *template.Template
	valTempl   *template.Template
	colTempl   *template.Template
}

func (s *server) registerHandlers() {
	http.HandleFunc("/", s.queryForm)
}

func (s *server) createAndRunServer(url string, port string) error {
	hs := &http.Server{
		Addr:           url + ":" + port,
		Handler:        nil,
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.registerHandlers()
	return hs.ListenAndServe()
}

func (s *server) setupTemplates() {
	s.queryTempl = template.Must(template.ParseFiles("queryform.html"))

	var err error
	s.valTempl, err = template.New(path.Base("valuator.html")).Funcs(template.FuncMap{
		"isYoyNonNil": func(m valuator.Measures) bool {
			if yoy := m.Yoy(); !reflect.ValueOf(yoy).IsNil() {
				return true
			}
			return false
		},
		"getYoy": func(m valuator.Measures) valuator.Yoy {
			return m.Yoy()
		},
		"dcfTrend": func(ticker string, dr float64, duration int, trend float64) float64 {
			ret, _ := s.valuator.DiscountedCashFlowTrend(ticker, dr, trend, duration)
			return ret
		},
		"dcfFCFTrend": func(ticker string, dr float64, duration int, trend float64) float64 {
			ret, _ := s.valuator.DiscountedFCFTrend(ticker, dr, trend, duration)
			return ret
		},
	}).ParseFiles("valuator.html")
	if err != nil {
		panic(err.Error())
	}

	s.colTempl, err = template.New(path.Base("collector.html")).Funcs(template.FuncMap{
		"must": func(fn func() (float64, error)) float64 {
			val, _ := fn()
			return val
		},
		"date": func(f valuator.Filing) string {
			return time.Time(f.FiledOn()).Format("2006-01-02")
		},
		"equity": func(f valuator.Filing) float64 {
			val, _ := f.TotalEquity()
			return val
		},
		"sharecount": func(f valuator.Filing) float64 {
			val, _ := f.ShareCount()
			return val
		},
		"revenue": func(f valuator.Filing) float64 {
			val, _ := f.Revenue()
			return val
		},
		"costofrevenue": func(f valuator.Filing) float64 {
			val, _ := f.CostOfRevenue()
			return val
		},
		"opsincome": func(f valuator.Filing) float64 {
			val, _ := f.OperatingIncome()
			return val
		},
		"opsexpense": func(f valuator.Filing) float64 {
			val, _ := f.OperatingExpense()
			return val
		},
		"netincome": func(f valuator.Filing) float64 {
			val, _ := f.NetIncome()
			return val
		},
		"stdebt": func(f valuator.Filing) float64 {
			val, _ := f.ShortTermDebt()
			return val
		},
		"ltdebt": func(f valuator.Filing) float64 {
			val, _ := f.LongTermDebt()
			return val
		},
		"cash": func(f valuator.Filing) float64 {
			val, _ := f.Cash()
			return val
		},
	}).ParseFiles("collector.html")
	if err != nil {
		panic(err.Error())
	}
}

func newServer(url string, dbType valuator.DatabaseType) (*server, error) {
	if db, err := valuator.NewDatabase(url, dbType); err == nil {
		if v, err := valuator.NewValuator(db); err == nil {
			s := &server{
				valuator: v,
			}
			s.setupTemplates()
			return s, nil
		}
		return nil, errors.New("Failed to create a valuator")
	}
	return nil, errors.New("Failed to create a file database")
}

func startServer(v valuator.Valuator, url string, port string) error {
	s := server{
		valuator: v,
	}
	s.setupTemplates()

	return s.createAndRunServer(url, port)
}
