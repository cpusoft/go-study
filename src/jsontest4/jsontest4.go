package main

import (
	"fmt"

	"github.com/cpusoft/goutil/jsonutil"
	"github.com/guregu/null"
)

type AsaToRtrFullLog struct {
	AsaId         uint64   `json:"roaId" xorm:"roaId int"`
	CustomerAsn   uint64   `json:"customerAsn" xorm:"customerAsn int"`
	ProviderAsn   uint64   `json:"providerAsns" xorm:"providerAsns int"`
	AddressFamily null.Int `json:"addressFamily" xorm:"addressFamily int"`
	SyncLogId     uint64   `json:"syncLogId" xorm:"syncLogId int"`
	SyncLogFileId uint64   `json:"syncLogFileId" xorm:"syncLogFileId int"`
}
type AsaStrToRtrFullLog struct {
	AsaId         uint64 `json:"roaId" xorm:"roaId int"`
	CustomerAsns  string `json:"customerAsns" xorm:"customerAsns varchar"`
	SyncLogId     uint64 `json:"syncLogId" xorm:"syncLogId int"`
	SyncLogFileId uint64 `json:"syncLogFileId" xorm:"syncLogFileId int"`
}
type CustomerAsn struct {
	CustomerAsn  uint64        `json:"customerAsn"`
	ProviderAsns []ProviderAsn `json:"ProviderAsns"`
}

type ProviderAsn struct {
	AddressFamily null.Int `json:"addressFamily"`
	ProviderAsn   uint64   `json:"providerAsn"`
}

func main() {
	// [{"customerAsn": 50555, "ProviderAsns": [{"providerAsn": 970, "addressFamily": null}], "addressFamily": null}]
	t := `[{"customerAsn": 50555, "ProviderAsns": [{"providerAsn": 970, "addressFamily": null}], "addressFamily": null}]`
	a := make([]CustomerAsn, 0)
	err := jsonutil.UnmarshalJson(t, &a)
	fmt.Println(a, err)
	fmt.Println(jsonutil.MarshalJson(a), err)
}
