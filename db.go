package valuator

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type DatabaseType string

var (
	fileDBUrl        string       = "./db/"
	FileDatabaseType DatabaseType = "file"
)

type database interface {
	Read(string) (io.Reader, error)
	Write(string, []byte) error
}

type fileDB struct {
	path   string
	reader map[string]io.Reader
	writer map[string]io.Writer
}

func NewDB(url string, ty DatabaseType) database {
	switch ty {
	case FileDatabaseType:
		return &fileDB{path: url}
	default:
		log.Println("Unknown Database type ", ty)
	}
	return nil
}

func (f *fileDB) generateFilePath(filename string) string {
	return f.path + filename + ".json"
}

func (f *fileDB) Read(filename string) (io.Reader, error) {
	file := f.generateFilePath(filename)
	fd, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Cannot open file ", file)
		return nil, errors.New("No data available at " + file)
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Println("Error reading file ", file)
		return nil, errors.New("Error reading file " + file)
	}
	reader := bytes.NewReader(data)
	fd.Close()
	return reader, nil
}

func (f *fileDB) Write(ticker string, data []byte) error {
	file := f.generateFilePath(ticker)
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Cannot open file ", file)
		return errors.New("Cannot open file for write " + file)
	}
	fd.Write(data)
	fd.Close()
	return nil
}
