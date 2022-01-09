package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/guregu/null"
)

type RushNodeModel struct {
	Id       uint64 `json:"id" xorm:"id int"`
	NodeName string `json:"nodeName" xorm:"nodeName varchar(256)"`

	// if it is root, will be null
	ParentNodeId   null.Int `json:"parentNodeId" xorm:"parentNodeId int"`
	ParentNodeName string   `json:"parentNodeName" xorm:"parentNodeName varchar(256)"`
	// interface url: https://1.1.1.1:8080
	Url string `json:"url" xorm:"url varchar(256)"`
	// 'true/null: vc to identify itself. rp do not need this
	IsSelfUrl string `json:"isSelfUrl" xorm:"isSelfUrl varchar(8)"`

	Note string `json:"note" xorm:"note varchar(256)"`
	//update time
	UpdateTime time.Time `json:"updateTime" xorm:"updateTime datetime"`
}

func main() {
	str := ` [{"id":1,"nodeName":"中心缓存服务器（Center-VC）","note":"中心缓存服务器","url":"https://202.173.14.104:8081"},{"id":2,"nodeName":"二级缓存服务器","note":"二级缓存服务器","parentNodeId":1,"url":"https://202.173.14.105:8081"},{"id":3,"nodeName":"三级缓存服务器","note":"三级缓存服务器","parentNodeId":2,"url":"https://202.173.14.103:8081"}]`
	rs := make([]RushNodeModel, 0)
	sss(str, &rs)
	fmt.Println("out sss:", rs)

}
func sss(str string, f interface{}) {

	//jsonutil.UnmarshalJson(str, &f)
	json.Unmarshal([]byte(str), &f)
	fmt.Println("in sss:", f)
}
