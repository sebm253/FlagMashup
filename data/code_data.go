package data

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	codesCDN = "https://flagcdn.com/en/codes.json"
	flagCDN  = "https://flagcdn.com/w2560/%s.jpg"
)

type CodeData struct {
	codes map[string]string // code:name
}

func (d *CodeData) Populate() {
	rs, err := http.Get(codesCDN)
	if err != nil {
		panic(err)
	}
	defer rs.Body.Close()
	if err := json.NewDecoder(rs.Body).Decode(&d.codes); err != nil {
		panic(err)
	}
	// append country prefix to subdivisions, i.e. "California" -> "United States - California"
	for code, name := range d.codes {
		split := strings.Split(code, "-") // us-ca
		if len(split) == 1 {              // no "-" separator
			continue
		}
		d.codes[code] = fmt.Sprintf("%s - %s", d.codes[split[0]], name) // country - subdivision
	}
	slog.Debug("populated code data", slog.Int("data.length", len(d.codes)))
}

func (d *CodeData) FetchFlag(code string) (io.ReadCloser, error) {
	rs, err := http.Get(fmt.Sprintf(flagCDN, code))
	if err != nil {
		return nil, err
	}
	return rs.Body, nil
}

func (d *CodeData) Map() map[string]string {
	return d.codes
}
