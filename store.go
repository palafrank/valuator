package valuator

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

type MarshaledData []byte

func (d MarshaledData) MarshalJSON() ([]byte, error) {
	return []byte(d), nil
}

func (d *MarshaledData) UnmarshalJSON(b []byte) error {
	*d = b
	return nil
}

type StoreEntry struct {
	Company     string        `json:"Company"`
	FinData     MarshaledData `json:"Financial Data"`
	FinMeasures MarshaledData `json:"Financial Measures"`
}

type StoreCollection struct {
	db      database
	Entries map[string]StoreEntry `json:"Ticker"`
}

func (s StoreEntry) String() string {

	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuator data: ", err)
	}
	return string(data)
}

func (s StoreCollection) String() string {
	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuator data: ", err)
	}
	return string(data)
}

/*
  The store is an interface to the in-memory store used by the Valuator
  The collector operates on the store data and manipulates it based on
  measurements and calculations
  The valuator populates the store initially based on the DB that it has and
  saves data from the store to the DB
*/
type Store interface {
	GetFinancials(string) io.Reader
	PutFinancials(string, []byte)
	GetMeasures(string) io.Reader
	PutMeasures(string, []byte)
	Write() error
	Read(string) error
	String() string
}

func NewStore(d database) Store {
	return &StoreCollection{
		db:      d,
		Entries: make(map[string]StoreEntry),
	}
}

func (s *StoreCollection) GetFinancials(ticker string) io.Reader {
	if entry, ok := s.Entries[ticker]; ok {
		return bytes.NewReader(entry.FinData)
	}
	return nil
}

func (s *StoreCollection) PutFinancials(ticker string, data []byte) {
	if entry, ok := s.Entries[ticker]; ok {
		entry.FinData = data
		entry.Company = ticker
		s.Entries[ticker] = entry
		return
	}
	entry := StoreEntry{
		FinData: data,
	}
	s.Entries[ticker] = entry
	return
}

func (s *StoreCollection) GetMeasures(ticker string) io.Reader {
	if entry, ok := s.Entries[ticker]; ok {
		return bytes.NewReader(entry.FinMeasures)
	}
	return nil
}

func (s *StoreCollection) PutMeasures(ticker string, data []byte) {
	if entry, ok := s.Entries[ticker]; ok {
		entry.FinMeasures = data
		entry.Company = ticker
		s.Entries[ticker] = entry
		return
	}
	entry := StoreEntry{
		FinMeasures: data,
	}
	s.Entries[ticker] = entry
	return
}

func (s *StoreCollection) Write() error {
	for ticker, entry := range s.Entries {

		if err := s.db.Write(ticker, []byte(entry.String())); err != nil {
			return err
		}

	}
	return nil
}

func (s *StoreCollection) Read(ticker string) error {
	if _, ok := s.Entries[ticker]; ok {
		return nil
	}
	r, err := s.db.Read(ticker)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	var se StoreEntry

	err = json.Unmarshal(b, &se)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return err
	}
	se.Company = ticker
	s.Entries[ticker] = se
	return nil
}
