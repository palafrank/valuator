package client

import (
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

type client struct {
	handle *http.Client
}

type Client interface {
	// GET queries the valuator for info.
	// The first parameter is the base URL of the valuator server
	// The second parameter is the query info that constructs the query
	GET(string, interface{}) (string, error)
}

func NewClient(c *http.Client) Client {
	return &client{
		handle: c,
	}
}
func (c *client) GET(baseURL string, query interface{}) (string, error) {
	url := baseURL
	switch reflect.TypeOf(query) {
	case reflect.TypeOf(&QueryAnnualData{}):
		queryData := query.(QueryAnnualData)
		url = url + "?ticker=" + queryData.Ticker
		if len(queryData.Years) > 0 {
			url = url + "&data=yes"
			for _, year := range queryData.Years {
				url = url + "&years=" + strconv.Itoa(year)
			}
		}
	default:
		log.Println("Unsupported GET operation")
	}

	res, err := c.handle.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(data), nil
}
