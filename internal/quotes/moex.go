package quotes

import (
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
)

type MarketData struct {
	Rows []Security `xml:"data>rows>row"`
}

type Security struct {
	Ticker         string  `xml:"SECID,attr" json:"Ticker"`
	Price          float32 `xml:"LAST,attr" json:"Price"`
	Time           string  `xml:"TIME,attr" json:"Time"`
	SeqNum         string  `xml:"SEQNUM,attr" json:"SeqNum"`
	Capitalization string  `xml:"ISSUECAPITALIZATION,attr" json:"Capitalization"`
	Pref           *Security
}

func Fetch(url string) (map[string]Security, error) {

	resp, err := http.Get(url)
	if err != nil {
		slog.Error(err.Error())
		return map[string]Security{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error(err.Error())
		return map[string]Security{}, err
	}
	//
	var data MarketData
	err = xml.Unmarshal(body, &data)
	if err != nil {
		slog.Error(err.Error())
		return map[string]Security{}, err
	}

	// Упаковать префы как вложенные структуры у обычных акий
	prefsMap := make(map[string]Security)
	for _, row := range data.Rows {
		pref, exists := PrefRev[row.Ticker]
		if exists {
			prefsMap[pref] = row
		}
	}

	// Convert the slice of Rows into a map
	rowsMap := make(map[string]Security)
	for _, row := range data.Rows {
		_, exists := PrefRev[row.Ticker]
		if exists {
			continue
		}

		val, exists := prefsMap[row.Ticker]
		if exists {
			row.Pref = &val
		}

		rowsMap[row.Ticker] = row
	}

	return rowsMap, nil
}

var PrefRev = map[string]string{
	"VJGZP": "VJGZ",
	"MISBP": "MISB",
	"RTSBP": "RTSB",
	"YKENP": "YKEN",
	"TASBP": "TASB",
	"VRSBP": "VRSB",
	"BANEP": "BANE",
	"BISVP": "BISV",
	"BSPBP": "BSPB",
	"CNTLP": "CNTL",
	"DZRDP": "DZRD",
	"IGSTP": "IGST",
	"JNOSP": "JNOS",
	"KCHEP": "KCHE",
	"KGKCP": "KGKC",
	"KRKNP": "KRKN",
	"KROTP": "KROT",
	"KRSBP": "KRSB",
	"KZOSP": "KZOS",
	"LNZLP": "LNZL",
	"LSNGP": "LSNG",
	"MAGEP": "MAGE",
	"MFGSP": "MFGS",
	"MGTSP": "MGTS",
	"MTLRP": "MTLR",
	"NKNCP": "NKNC",
	"NNSBP": "NNSB",
	"PMSBP": "PMSB",
	"RTKMP": "RTKM",
	"SAGOP": "SAGO",
	"SAREP": "SARE",
	"SBERP": "SBER",
	"SNGSP": "SNGS",
	"STSBP": "STSB",
	"SVETP": "SVET",
	"TATNP": "TATN",
	"TGKBP": "TGKB",
	"TORSP": "TORS",
	"TRNFP": "TRNF",
	"VGSBP": "VGSB",
	"VSYDP": "VSYD",
	"WTCMP": "WTCM",
	"YRSBP": "YRSB",
}

var Pref = map[string]string{
	"VJGZ": "VJGZP",
	"MISB": "MISBP",
	"RTSB": "RTSBP",
	"YKEN": "YKENP",
	"TASB": "TASBP",
	"VRSB": "VRSBP",
	"BANE": "BANEP",
	"BISV": "BISVP",
	"BSPB": "BSPBP",
	"CNTL": "CNTLP",
	"DZRD": "DZRDP",
	"IGST": "IGSTP",
	"JNOS": "JNOSP",
	"KCHE": "KCHEP",
	"KGKC": "KGKCP",
	"KRKN": "KRKNP",
	"KROT": "KROTP",
	"KRSB": "KRSBP",
	"KZOS": "KZOSP",
	"LNZL": "LNZLP",
	"LSNG": "LSNGP",
	"MAGE": "MAGEP",
	"MFGS": "MFGSP",
	"MGTS": "MGTSP",
	"MTLR": "MTLRP",
	"NKNC": "NKNCP",
	"NNSB": "NNSBP",
	"PMSB": "PMSBP",
	"RTKM": "RTKMP",
	"SAGO": "SAGOP",
	"SARE": "SAREP",
	"SBER": "SBERP",
	"SNGS": "SNGSP",
	"STSB": "STSBP",
	"SVET": "SVETP",
	"TATN": "TATNP",
	"TGKB": "TGKBP",
	"TORS": "TORSP",
	"TRNF": "TRNFP",
	"VGSB": "VGSBP",
	"VSYD": "VSYDP",
	"WTCM": "WTCMP",
	"YRSB": "YRSBP",
}
