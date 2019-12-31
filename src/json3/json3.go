package main

import (
	"fmt"
	jsonutil "github.com/cpusoft/goutil/jsonutil"
)

type DD struct {
	Data      Data   `json:"data"`
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type Data struct {
	AcceptTime uint64 `json:"acceptTime"`
	AllotTime  uint64 `json:"allotTime"`

	CompleteTime   uint64 `json:"completeTime"`
	CreateTime     uint64 `json:"createTime"`
	CreateUserName string `json:"createUserName"`

	CusName      string `json:"cusName"`
	CusNo        string `json:"cusNo"`
	ExecutorName string `json:"executorName"`
	Level        string `json:"level"`
	LmName       string `json:"lmName"`
	LmPhone      string `json:"lmPhone"`
	PlanTime     string `json:"planTime"`
}

func main() {
	json := `{"data":{"acceptTime":1574308121000,"allotTime":1574308072000,"attribute":{"customReceipt":true,"receiptAttachment":[],"field_GNeaMMPazvb00TJq":"4","field_IjEuhP1yQ4BtTDkR":"你的孩子也不一定知道我不爱了","field_RwFFYKle2mKkSwM7":"正式用户","field_VpMNFXHJYtsTmKrO":"dadsadas","field_XZQsFaBBs5qK6G0T":"测测测测测测测","field_jxmUGahdyTOR7bun":"SHB2019112104","field_v77jZPvGNvRJvB4F":["IPAM"]},"completeTime":1574322671000,"createTime":1574308067000,"createUserName":"李昂","cusAddress":{"city":"长春市","dist":"二道区","address":"长春理工大学","latitude":43.865596,"province":"吉林省","longitude":125.374217},"cusName":"长春理工大学","cusNo":"CUSWG00756","executorName":"李昂","level":"中","lmName":"陈老师","lmPhone":"13478478909","planTime":"Fri Nov 22 11:45:00 CST 2019","productIds":[],"products":[],"startTime":1574322622000,"state":"已完成","synergyNames":[],"taskNo":"TLQ1019110173","templateId":"673e91db-9081-45a9-bbe2-809cc5542ec1","templateName":"技术咨询申请","updateTime":1574322670000},"errorCode":"0","message":"调用接口成功"}`
	dd := DD{}
	jsonutil.UnmarshalJson(json, &dd)
	fmt.Println(dd)
}
