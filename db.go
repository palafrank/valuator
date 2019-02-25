package valuator

import (
	"errors"
	"io"
	"log"
)

type DatabaseType string

var (
	FileDatabaseType  DatabaseType = "file"
	MongoDatabaseType DatabaseType = "mongo"
	NoneDatabaseType  DatabaseType = "none"
)

type database interface {
	Open() error
	Close()
	Read(string) (io.Reader, error)
	Write(string, []byte) error
}

func NewDatabase(url interface{}, ty DatabaseType) (database, error) {
	var db database
	switch ty {
	case FileDatabaseType:
		if _, ok := url.(string); ok {
			db = NewFileDB(url.(string))
		}
	case NoneDatabaseType:
		return &noneDB{}, nil
	default:
		log.Println("Unknown Database type ", ty)
		return nil, errors.New("Unknown database type " + string(ty))
	}
	if err := db.Open(); err != nil {
		return nil, err
	}
	return db, nil
}

type noneDB struct {
}

func (n *noneDB) Open() error {
	return nil
}

func (n *noneDB) Close() {

}

func (n *noneDB) Read(string) (io.Reader, error) {
	return nil, errors.New("None database has no read destination")
}

func (n *noneDB) Write(string, []byte) error {
	return errors.New("None database has no write destination")
}
