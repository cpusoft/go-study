package main

import (
	"fmt"
	"time"

	"github.com/cpusoft/goutil/jsonutil"
)

type MyTime time.Time

const (
	timeFormart = "2006-01-02T15:04:05Z07:00"
)

func (t *MyTime) UnmarshalJSON(data []byte) (err error) {

	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = MyTime(now)
	return
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}
func (t MyTime) String() string {
	return time.Time(t).Format(timeFormart)
}

type Pfxes map[string]PfxeData
type PfxeData struct {
	Pfx         string                 `json:"pfx"`
	Competition map[string]Competition `json:"competition"`
	Key         string                 `json:"key"`
	Len         uint64                 `json:"len"`
	MaxLen      uint64                 `json:"maxlen"`
}
type Competition struct {
	RelativePosition uint64 `json:"relative_position"`
	RoaName          string `json:"roa_name"`
}

type RoaComp struct {
	CompFlag      bool     `json:"comp_flag"`
	CertChainName []string `json:"certchain_name"`
	State         string   `json:"state"`
	Asn           uint64   `json:"asn"`
	ValFrom       MyTime   `json:"valfrom"`
	ValTo         MyTime   `json:"valto"`
	PfxesData     Pfxes    `json:"pfxes"`
}

func main() {

	str := `
{
  "comp_flag": true,
  "certchain_name": [
    "apnic-rpki-root-iana-origin.cer",
    "mBQsnQtBo7n7YD12mEgjb9HzGSQ.cer",
    "DmWk9f02tb1o6zySNAiXjJB6p58.cer",
    "fWXr4UwwRuc-OYaVcwKibggfGvg.cer"
  ],
  "state": "{\"state\": \"valid\", \"errors\": [], \"warnings\": []}",
  "valfrom": "2019-09-20T18:11:26+08:00",
  "pfxes": {
    "20431": {
      "pfx": "125.5.108.0/24",
      "competition": {
        "62648": {
          "roa_name": "E4E0C13C5EBF11EAB066F050C4F9AE02.roa",
          "relative_position": 3
        },
        "62626": {
          "roa_name": "E4E0C13C5EBF11EAB066F050C4F9AE02.roa",
          "relative_position": 3
        },
        "62643": {
          "roa_name": "E4E0C13C5EBF11EAB066F050C4F9AE02.roa",
          "relative_position": 3
        }
      },
      "len": 24,
      "key": "4$011111010000010101101100",
      "maxlen": 24
    }
  },
  "valto": "2020-12-01T08:00:00+08:00",
  "asn": 45791
}
`
	fmt.Println(str)

	roaComp := RoaComp{}
	err := jsonutil.UnmarshalJson(str, &roaComp)
	fmt.Println(jsonutil.MarshalJson(roaComp))
	fmt.Println(err)
}
