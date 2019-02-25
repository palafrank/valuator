package valuator

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type fileDB struct {
	path   string
	writer map[string]*os.File
}

func NewFileDB(url string) *fileDB {
	return &fileDB{
		path:   url,
		writer: make(map[string]*os.File),
	}
}

func (f *fileDB) generateFilePath(filename string) string {
	return f.path + filename + ".json"
}

func (f *fileDB) Open() error {
	_, err := os.Stat(f.path)
	return err
}

func (f *fileDB) Close() {
	for _, writer := range f.writer {
		writer.Close()
	}
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
	if _, ok := f.writer[ticker]; !ok {
		file := f.generateFilePath(ticker)
		fd, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			log.Println("Cannot open file ", file)
			return errors.New("Cannot open file for write " + file)
		}
		f.writer[ticker] = fd
	}
	f.writer[ticker].Write(data)
	return nil
}
