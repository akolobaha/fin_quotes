package quotes

import (
	"encoding/xml"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type MarketData struct {
	Rows []Row `xml:"data>rows>row"`
}

type Row struct {
	SECID string  `xml:"SECID,attr"`
	LAST  float32 `xml:"LAST,attr"`
	TIME  string  `xml:"TIME,attr"`
}

func Fetch(url string) (MarketData, error) {

	resp, err := http.Get(url)
	if err != nil {
		slog.Error(err.Error())
		return MarketData{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		return MarketData{}, err
	}
	//
	var data MarketData
	err = xml.Unmarshal(body, &data)
	if err != nil {
		slog.Error(err.Error())
		return MarketData{}, err
	}

	return data, nil
}
