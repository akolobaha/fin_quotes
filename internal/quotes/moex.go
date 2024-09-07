package quotes

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type Data struct {
	XMLName xml.Name `xml:"data"`
	Rows    Rows     `xml:"rows"`
}

type Rows struct {
	XMLName xml.Name `xml:"rows"`
	Row     []Ticker `xml:"row"`
}

type Ticker struct {
	Name  string    `xml:"Name,attr"`
	Price string    `xml:"Price,attr"`
	Time  time.Time `xml:"Time,attr"`
}

func Fetch(url string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	//io.Copy(os.Stdout, resp.Body)

}
