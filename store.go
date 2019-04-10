package valuator

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
)

// MarshaledData is a wrapper of already marshalled data
type MarshaledData []byte

// MarshalJSON is a wrapper to return already marshalled data from store
func (d MarshaledData) MarshalJSON() ([]byte, error) {
	return []byte(d), nil
}

// UnmarshalJSON is a wrapper to return already unmarshalled data to the store
func (d *MarshaledData) UnmarshalJSON(b []byte) error {
	*d = b
	return nil
}

type storeEntry struct {
	Company     string        `json:"Company"`
	FinData     MarshaledData `json:"Financial Data"`
	FinMeasures MarshaledData `json:"Financial Measures"`
}

type storeCollection struct {
	db      Database
	Entries map[string]storeEntry `json:"Ticker"`
}

func (s storeEntry) String() string {

	data, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling valuator data: ", err)
	}
	return string(data)
}

func (s storeCollection) String() string {
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

// Store interface provides a read & write interface with the in-memory store
type Store interface {
	getFinancials(string) io.Reader
	putFinancials(string, []byte)
	getMeasures(string) io.Reader
	putMeasures(string, []byte)
	write() error
	read(string) error
	String() string
}

func newStore(d Database) Store {
	return &storeCollection{
		db:      d,
		Entries: make(map[string]storeEntry),
	}
}

func (s *storeCollection) getFinancials(ticker string) io.Reader {
	if entry, ok := s.Entries[ticker]; ok {
		return bytes.NewReader(entry.FinData)
	}
	return nil
}

func (s *storeCollection) putFinancials(ticker string, data []byte) {
	if entry, ok := s.Entries[ticker]; ok {
		entry.FinData = data
		entry.Company = ticker
		s.Entries[ticker] = entry
		return
	}
	entry := storeEntry{
		FinData: data,
	}
	s.Entries[ticker] = entry
	return
}

func (s *storeCollection) getMeasures(ticker string) io.Reader {
	if entry, ok := s.Entries[ticker]; ok {
		return bytes.NewReader(entry.FinMeasures)
	}
	return nil
}

func (s *storeCollection) putMeasures(ticker string, data []byte) {
	if entry, ok := s.Entries[ticker]; ok {
		entry.FinMeasures = data
		entry.Company = ticker
		s.Entries[ticker] = entry
		return
	}
	entry := storeEntry{
		FinMeasures: data,
	}
	s.Entries[ticker] = entry
	return
}

func (s *storeCollection) write() error {
	for ticker, entry := range s.Entries {

		if err := s.db.Write(ticker, []byte(entry.String())); err != nil {
			return err
		}

	}
	return nil
}

func (s *storeCollection) read(ticker string) error {
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
	var se storeEntry

	err = json.Unmarshal(b, &se)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return err
	}
	se.Company = ticker
	s.Entries[ticker] = se
	return nil
}
