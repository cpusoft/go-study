package main

import (
	_ "database/sql"
	"fmt"
	"time"

	//"github.com/cpusoft/goutil/dnsutil"
	"github.com/cpusoft/goutil/jsonutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guregu/null"

	"xorm.io/xorm"
)

// for mysql, zonefile
type RrModel struct {
	Id       uint64 `json:"id" xorm:"id int"`
	OriginId uint64 `json:"originId" xorm:"originId int"`

	// not have "." in the end
	Origin string `json:"origin" xorm:"origin varchar"` // lower
	// is host/subdomain, not have "." int the end
	// if no subdomain, is "", not "@"
	RrName string `json:"rrName" xorm:"rrName varchar"` // lower
	// == rrName+.+Origin
	RrFullDomain string `json:"rrFullDomain" xorm:"rrFullDomain varchar"` // lower: rrName+"."+Origin[-"."]

	RrType  string `json:"rrType" xorm:"rrType varchar"`   // upper
	RrClass string `json:"rrClass" xorm:"rrClass varchar"` // upper
	// null.NewInt(0, false) or null.NewInt(i64, true)
	RrTtl  null.Int `json:"rrTtl" xorm:"rrTtl int"`
	RrData string   `json:"rrData" xorm:"rrData varchar"`

	UpdateTime time.Time `json:"updateTime" xorm:"updateTime datetime"`
}

const (
	LICENSE_TIME_FORMAT = time.RFC3339 //"2006-01-02T15:04:05Z07:00" // RFC3339
)

type LicenseTime time.Time

func (t *LicenseTime) UnmarshalJSON(data []byte) (err error) {

	now, err := time.ParseInLocation(`"`+LICENSE_TIME_FORMAT+`"`, string(data), time.Local)
	*t = LicenseTime(now)
	return
}

func (t LicenseTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(LICENSE_TIME_FORMAT)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, LICENSE_TIME_FORMAT)
	b = append(b, '"')
	return b, nil
}

func (t LicenseTime) String() string {
	return time.Time(t).Format(LICENSE_TIME_FORMAT)
}

type LicenseInfoModel struct {
	// defined in shell
	UserName   string `json:"userName" xorm:"userName varchar(255)"`     // "zdns1", "cnic1"
	DeviceName string `json:"deviceName" xorm:"deviceName varchar(255)"` // "1","2",
	DeviceUuid string `json:"deviceUuid" xorm:"deviceUuid varchar(255)"`
	// m, _ := time.ParseDuration("-1m")
	// newTime := now.Add(m)
	DefaultDevicePeriod string `json:"defaultDevicePeriod" xorm:"defaultDevicePeriod varchar(16)"` //"1Y","1M","1D"

	// key
	KeyId uint64 `json:"keyId" xorm:"defaultDevicePeriod int"`

	InstallTime      LicenseTime `json:"installTime" xorm:"installTime datetime"`
	LicenseStartTime LicenseTime `json:"licenseStartTime" xorm:"licenseStartTime datetime"`
	LicenseEndTime   LicenseTime `json:"licenseEndTime" xorm:"licenseEndTime datetime"`

	// systeminfosha56
	SystemInfo  string `json:"systemInfo" xorm:"systemInfo json"`
	SignTimeStr string `json:"signTimeStr"`
}

func main() {
	//DB, err = sql.Open("mysql", "rpstir:Rpstir-123@tcp(202.173.9.21:13306)/rpstir")
	//user := "dns"
	//password := "Dns-123"
	//server := "202.173.14.104:13307"
	//database := "dns"

	user := "rpstir2"
	password := "Rpstir-123"
	server := "10.1.135.22:13306"
	database := "rpstir2"
	maxidleconns := 50
	maxopenconns := 50

	openSql := user + ":" + password + "@tcp(" + server + ")/" + database

	//连接数据库
	engine, err := xorm.NewEngine("mysql", openSql)
	if err != nil {
		fmt.Println(err)

		return
	}
	//连接测试
	if err := engine.Ping(); err != nil {
		fmt.Println(err)

		return
	}

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(maxidleconns)
	//设置最大打开连接数
	engine.SetMaxOpenConns(maxopenconns)
	//engine.SetTableMapper(core.SnakeMapper{})
	engine.ShowSQL(true)
	/*
		session := engine.NewSession()
		defer session.Close()

			rrType := "A"
			rrFullDomain := "dns1.example.com"
			resultRrModels := make([]*RrModel, 0)
			fmt.Println("queryDb(): rrType:", rrType, "  rrFullDomain:", rrFullDomain)
			if rrType == dnsutil.DNS_TYPE_STR_ANY {
				sql := `select o.origin, r.rrName, r.rrFullDomain, r.rrType, r.rrClass, IFNULL(r.rrTtl,o.ttl) as rrTtl, r.rrData
					from lab_dns_rr r,	lab_dns_origin o
					where r.originId = o.id and r.rrFullDomain = ?
					group by r.id `
				err = engine.SQL(sql, rrFullDomain).Find(&resultRrModels)
			} else {
				sql := `select o.origin, r.rrName, r.rrFullDomain, r.rrType, r.rrClass, IFNULL(r.rrTtl,o.ttl) as rrTtl, r.rrData
					from lab_dns_rr r, lab_dns_origin o
					where r.originId = o.id and r.rrFullDomain = ? and r.rrType = ?
					group by r.id `
				err = engine.SQL(sql, rrFullDomain, rrType).Find(&resultRrModels)
			}
			if err != nil {
				fmt.Println("queryDb(): lab_dns_rr fail:", err)
				return
			}
			fmt.Println("queryDb(): resultRrModels:", jsonutil.MarshalJson(resultRrModels))
	*/

	start := time.Now()
	licenseInfoModelDb := LicenseInfoModel{}
	keyId := 7
	deviceName := "beijing1"
	deviceUuid := "c0e56e3e-1259-45c4-a672-cd7449ae3d1c"
	sql := `select u.userName, c.keyId, c.deviceUuid, c.deviceName,c.publicInfo as systemInfo, c.installTime, c.licenseStartTime, c.licenseEndTime 
		    from lab_rpki_license_device c ,lab_rpki_license_user u
			where c.userId = u.id and  c.keyId = ? and c.deviceName = ? and c.deviceUuid = ?`
	has, err := engine.SQL(sql, keyId, deviceName, deviceUuid).Get(&licenseInfoModelDb)
	if err != nil {
		fmt.Println("getLicenseDeviceAndUpdatePublicInfoDb(): get license_device fail, keyId:", keyId, " deviceName:", deviceName, " deviceUuid:", deviceUuid, err)
		return
	}
	if !has {
		fmt.Println("getLicenseDefaultDevicePeriodDb(): not found in license_device fail, keyId:", keyId,
			" deviceName:", deviceName, " deviceUuid:", deviceUuid, "  time(s):", time.Since(start))
		return
	}
	fmt.Println("getLicenseDeviceAndUpdatePublicInfoDb(): licenseInfoModelDb:", jsonutil.MarshalJson(licenseInfoModelDb))

}
