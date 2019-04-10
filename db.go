package valuator

import (
	"errors"
	"io"
	"log"
)

// DatabaseType is a type definition of the different databases supported by valuator
type DatabaseType string

// FileDatabaseType is a file database type
const FileDatabaseType DatabaseType = "file"

// MongoDatabaseType is a mongo database type
const MongoDatabaseType DatabaseType = "mongo"

// NoneDatabaseType is a no database type
const NoneDatabaseType DatabaseType = "none"

// Database is an interface to the database used by the valuator
type Database interface {
	Open() error
	Close()
	Read(string) (io.Reader, error)
	Write(string, []byte) error
}

// NewDatabase creates a new database to be used by valuator
func NewDatabase(url interface{}, ty DatabaseType) (Database, error) {
	var db Database
	switch ty {
	case FileDatabaseType:
		if _, ok := url.(string); ok {
			db = newFileDB(url.(string))
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
