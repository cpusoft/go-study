package main

import (
	xml "encoding/xml"
	"fmt"
)

//https://stackoverflow.com/questions/23126133/golang-xml-attribute-and-value
/*
https://blog.csdn.net/caoyujiao520/article/details/81710242
标签使用介绍：
“-“：不会输出
“name,attr”：以name作为属性名
“,attr”：以这个struct的字段名作为属性名输出为XML元素的属性
“,chardata”：输出为xml的 character data而非element
“,innerxml”：被原样输出
“,comment”：将被当作xml注释来输出，字段值中不能含有”–”字符串
“omitempty”：如果该字段的值为空值那么该字段就不会被输出到XML，空值包括：false、0、nil指针或nil接口，任何长度为0的array, slice, map或者string

https://www.cnblogs.com/jkko123/p/8325813.html
    XMLName xml.Name `xml:"book"`;
    Name    string   `xml:"name,attr"`;
    Author  string   `xml:"author"`;
    Time    string   `xml:"time"`;
    //字段定义如a>b>c，这样，解析时会从xml当前节点向下寻找元素并将值赋给该字段
    Types []string `xml:"types>type"`;
    //字段定义有,any，则解析时如果xml元素没有与任何字段匹配，那么这个元素就会映射到该字段
    Test string `xml:",any"`;

https://rrdp.ripe.net/notification.xml
<notification xmlns="http://www.ripe.net/rpki/rrdp" version="1" session_id="26918a47-eb1c-45bc-bba8-05c4d6ea4b83" serial="200">
<snapshot uri="https://rrdp.ripe.net/26918a47-eb1c-45bc-bba8-05c4d6ea4b83/200/snapshot.xml" hash="CE0D3B09976633B7ADD656790B26B24491964A3AD65632199EA60381198FEF6C"/>
<delta serial="200" uri="https://rrdp.ripe.net/26918a47-eb1c-45bc-bba8-05c4d6ea4b83/200/delta.xml" hash="DE5152AEE57FE62DB9BFEA834266C83F4686D9DEE6FE5C49C109C2AEA72DDCA1"/>
<delta serial="199" uri="https://rrdp.ripe.net/26918a47-eb1c-45bc-bba8-05c4d6ea4b83/199/delta.xml" hash="3AE3684C0646D01548F9E73E4185135614F43816B4291AFB161A73D6771DEC39"/>
</notification>
*/
type Notification struct {
	XMLName    xml.Name             `xml:"notification"`
	Xmlns      string               `xml:"xmlns,attr"`
	Version    string               `xml:"version,attr"`
	Session_id string               `xml:"session_id,attr"`
	Serial     string               `xml:"serial,attr"`
	Snapshot   NotificationSnapshot `xml:"snapshot"`
	Deltas     []NotificationDelta  `xml:"delta"`
}

type NotificationSnapshot struct {
	XMLName xml.Name `xml:"snapshot"`
	Uri     string   `xml:"uri,attr"`
	Hash    string   `xml:"hash,attr"`
}
type NotificationDelta struct {
	XMLName xml.Name `xml:"delta"`
	Serial  string   `xml:"serial,attr"`
	Uri     string   `xml:"uri,attr"`
	Hash    string   `xml:"hash,attr"`
}

type Snapshot struct {
	XMLName          xml.Name          `xml:"snapshot"`
	Xmlns            string            `xml:"xmlns,attr"`
	Version          string            `xml:"version,attr"`
	Session_id       string            `xml:"session_id,attr"`
	Serial           string            `xml:"serial,attr"`
	SnapshotPublishs []SnapshotPublish `xml:"publish"`
}

type SnapshotPublish struct {
	XMLName xml.Name `xml:"publish"`
	Uri     string   `xml:"uri,attr"`
	Base64  string   `xml:",chardata"`
}

type Delta struct {
	XMLName        xml.Name        `xml:"delta"`
	Xmlns          string          `xml:"xmlns,attr"`
	Version        string          `xml:"version,attr"`
	Session_id     string          `xml:"session_id,attr"`
	Serial         string          `xml:"serial,attr"`
	DeltaPublishs  []DeltaPublish  `xml:"publish"`
	DeltaWithdraws []DeltaWithdraw `xml:"withdraw"`
}
type DeltaPublish struct {
	XMLName xml.Name `xml:"publish"`
	Uri     string   `xml:"uri,attr"`
	Hash    string   `xml:"hash,attr"`
	Base64  string   `xml:",chardata"`
}
type DeltaWithdraw struct {
	XMLName xml.Name `xml:"withdraw"`
	Uri     string   `xml:"uri,attr"`
	Hash    string   `xml:"hash,attr"`
}

func main() {
	notficationXml := `
	<notification xmlns="http://www.ripe.net/rpki/rrdp" version="1" session_id="26918a47-eb1c-45bc-bba8-05c4d6ea4b83" serial="200">
		<snapshot uri="https://rrdp.ripe.net/26918a47-eb1c-45bc-bba8-05c4d6ea4b83/200/snapshot.xml" hash="CE0D3B09976633B7ADD656790B26B24491964A3AD65632199EA60381198FEF6C"/>
		<delta serial="200" uri="https://rrdp.ripe.net/26918a47-eb1c-45bc-bba8-05c4d6ea4b83/200/delta.xml" hash="DE5152AEE57FE62DB9BFEA834266C83F4686D9DEE6FE5C49C109C2AEA72DDCA1"/>
		<delta serial="199" uri="https://rrdp.ripe.net/26918a47-eb1c-45bc-bba8-05c4d6ea4b83/199/delta.xml" hash="3AE3684C0646D01548F9E73E4185135614F43816B4291AFB161A73D6771DEC39"/>
	</notification>`
	var notification Notification
	err := xml.Unmarshal([]byte(notficationXml), &notification)
	if err != nil {
		fmt.Println("读文件内容错误信息：", err)
	}
	s := fmt.Sprintf("%+v", notification)
	fmt.Println(s)

	snapshotXml := `     
	   <snapshot version="1" session_id="26918a47-eb1c-45bc-bba8-05c4d6ea4b83" serial="200" xmlns="http://www.ripe.net/rpki/rrdp">
	       <publish uri="rsync://rpki.ripe.net/repository/DEFAULT/d9/6b8415-c876-46c6-a9f6-e6bbef5d5e3e/1/Oz1goI0QiX5nqx9IHWA-5WR7Uwk.roa">
		        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
	       </publish>
	       <publish uri="rsync://rpki.ripe.net/repository/DEFAULT/90/a28e5b-2eb2-4dc4-89e2-fcebb7e57412/1/2mUtDw3JD5zgtz_eoPK96r_UGfw.roa">
		        MIAGCSqGSIb3DQEHAqCAMIACAQMxDzANBglghkgBZQMEAgEFADCABgsqhkiG9w0BCRABGKCAJIAEHDAaAgMA74UwEzARBAIAATALMAkDBAJtagACARgAAAAAAACggDCCBO8wggPXoAMCAQICBAFftNIwDQYJKoZIhvcNAQELBQAwMzExMC8GA1UEAxMoMjEzYzQ3NmQ0YmU2NjRkOTM4ODdkYjcxZmRiYjE2YWY1Y2MzMDFmYzAeFw0xOTA5MDUxMjA1MDlaFw0yMDA3MDEwMDAwMDBaMDMxMTAvBgNVBAMTKGRhNjUyZDBmMGRjOTBmOWNlMGI3M2ZkZWEwZjJiZGVhYmZkNDE5ZmMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQD1yzOlP6ER0Ebwwzf7xvjtaz2avGTo28xcxDRlzxORJgUOH8BE/ihccoUqOpJaWeJ/PfkmuROMMNGlxB95vlw3dt4KWIXccu12dxMR+q40Ln9mpY6TFBB5hSg65CISDYw4HPSz+hJzIAwwWNqvJ5x3y401THUrM3WiUOs/Qr1VZ+GLwkM/jh3R6SfpSeAVj2qZ/KFjpD5iIPMsv3Wr7Z6/aHUvSbXXZJa7IFHH+BAlq6H3dNxs6K0RNDg+Zj85WGcPob+NmNwIc7MxZLV1Ow1fEyqP/jqBv6Jf4/D9L5IH1vI1seg2PCf2+WMkg8GnzuswK++Gce1tHkFZ49ColkuHAgMBAAGjggIJMIICBTAdBgNVHQ4EFgQU2mUtDw3JD5zgtz/eoPK96r/UGfwwHwYDVR0jBBgwFoAUITxHbUvmZNk4h9tx/bsWr1zDAfwwDgYDVR0PAQH/BAQDAgeAMGQGCCsGAQUFBwEBBFgwVjBUBggrBgEFBQcwAoZIcnN5bmM6Ly9ycGtpLnJpcGUubmV0L3JlcG9zaXRvcnkvREVGQVVMVC9JVHhIYlV2bVpOazRoOXR4X2JzV3IxekRBZncuY2VyMIGNBggrBgEFBQcBCwSBgDB+MHwGCCsGAQUFBzALhnByc3luYzovL3Jwa2kucmlwZS5uZXQvcmVwb3NpdG9yeS9ERUZBVUxULzkwL2EyOGU1Yi0yZWIyLTRkYzQtODllMi1mY2ViYjdlNTc0MTIvMS8ybVV0RHczSkQ1emd0el9lb1BLOTZyX1VHZncucm9hMIGBBgNVHR8EejB4MHagdKByhnByc3luYzovL3Jwa2kucmlwZS5uZXQvcmVwb3NpdG9yeS9ERUZBVUxULzkwL2EyOGU1Yi0yZWIyLTRkYzQtODllMi1mY2ViYjdlNTc0MTIvMS9JVHhIYlV2bVpOazRoOXR4X2JzV3IxekRBZncuY3JsMBgGA1UdIAEB/wQOMAwwCgYIKwYBBQUHDgIwHwYIKwYBBQUHAQcBAf8EEDAOMAwEAgABMAYDBAJtagAwDQYJKoZIhvcNAQELBQADggEBABvbBnqOhw+RqzXogdH0T+wZE7mjdAZy0Qlx+s479lMCTcU0SP31TRH6PCPnNmiIFKuURM/6CakqJGKEtvSoGSOyIJbsN/lQtYyE0VSBDHe77oLpaGHqMzNOVOn5LXZw8vT1rdAM//OdKo5vLmiLLBEQe9TxiKTPCNKPTf6gvb0d7ZQYg3PiNETmJ3G2VoDuUlpdHfIGr+omN1EYg2hOgOSw5jv1MjpgwSqh0QkNqRrzj7BbXHuTxfFQJOQCxjBArKVsqwcTkIjzm7JOxfENN5xLW4c32sZ8zuiufnzM+qvLTgYj4zlh2EttvNI69rBh8U/N8FLo+qsAsvapQzdlWbgAADGCAawwggGoAgEDgBTaZS0PDckPnOC3P96g8r3qv9QZ/DANBglghkgBZQMEAgEFAKBrMBoGCSqGSIb3DQEJAzENBgsqhkiG9w0BCRABGDAcBgkqhkiG9w0BCQUxDxcNMTkwOTA1MTIwNTA5WjAvBgkqhkiG9w0BCQQxIgQgMVTJ069JBQ6iO48LwrGnbQ+xXYTbBi32NPb+4k7cSoowDQYJKoZIhvcNAQELBQAEggEACuGjqkCQEMmKp/4VHU6qp6+5Fcc8P6cbQNwJJBEY0GqTUSBjCBcExlITXgCubAtm2+ss9wcFtii4ce0xtNeRwfBAdVV/XOB4YtIHfLQXHzwIupgslOT5rDqApiMzcJK6fxvVY1uG0Hido/Y+mkWgHS4GUUL9snDMiusBFDMgB0NaYPGtuH0IHomVz5OacWaGBrRfz0TR5ou0Hr5zmxzY8v9p1O8bHDv1SB6DzTDvexc+ZrV93QV1jLrPBPU9pelxt0YS9YcFfcajM60C0K9/EhgjWMFCyLpG2vtiEMPE7s/s2cqMMAjMGoQFI+x3zrJduBGDhflFWAj1f+irVna7YgAAAAAAAA==
	       </publish>
       </snapshot>`
	snapshot := Snapshot{}
	err = xml.Unmarshal([]byte(snapshotXml), &snapshot)
	if err != nil {
		fmt.Println("读文件内容错误信息：", err)
	}
	s = fmt.Sprintf("%+v", snapshot)
	fmt.Println(s)

	deltaXml := `
	   <delta xmlns="http://www.ripe.net/rpki/rrdp"
            version="1"
            session_id="26918a47-eb1c-45bc-bba8-05c4d6ea4b83"
            serial="200">
	       <publish uri="rsync://rpki.ripe.net/repository/DEFAULT/64/53d56a-096e-458a-9a94-08a68632a63a/1/8B0aTEXpIj0ULfyTJ8IeHltEGh8.mft"
                hash="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       
	       </publish>
	       <publish uri="rrsync://rpki.ripe.net/repository/DEFAULT/64/53d56a-096e-458a-9a94-08a68632a63a/1/8B0aTEXpIj0ULfyTJ8IeHltEGh8.crl"
                hash="2218FE34EC50253DD4D084F146CC02D6B9A1DF051C4F764421B8545555DFA80C">
		       MIIBrjCBlwIBATANBgkqhkiG9w0BAQsFADAzMTEwLwYDVQQDEyhmMDFkMWE0YzQ1ZTkyMjNkMTQyZGZjOTMyN2MyMWUxZTViNDQxYTFmFw0xOTExMTgwNjA2NTJaFw0xOTExMTkwNjA2NTJaoDAwLjAfBgNVHSMEGDAWgBTwHRpMRekiPRQt/JMnwh4eW0QaHzALBgNVHRQEBAICANwwDQYJKoZIhvcNAQELBQADggEBAINsztMVkeF1vEXun02jC1vrlwyeSfbelM/gnkYFaCaoAI/Yhc0K6ymbyUvJyXtScojumNotjU25oxBYRz7PJ3A0DiJ0U4IJBjFPaPhlYGC6gSMdS66oPCdmVYhWPwnNTYcbucHncHt+BJKETKm+vNc0LdYDF17fAkxC7tKY3TinZgh29rHaj9qNPC+wkW3FDpHMovO4klNCpwJ8xpSk3RglxmPFAXkacEXewk0eUVNVwKgXlGW6d19taUJp0kZG3DIADylLiicmHbx6dnGX3912NigJd7Jv6wK2TIk5hhLdK9ukiuUzSdd2xUBh4yb8Tnja2vo0tjC5SNWwPEkdC8I=
	       </publish>
	       <withdraw uri="rsync://rpki.ripe.net/repo/Alice/Bob.cer"
                 hash="caeb...15c1"/>
      </delta>`
	delta := Delta{}
	err = xml.Unmarshal([]byte(deltaXml), &delta)
	if err != nil {
		fmt.Println("读文件内容错误信息：", err)
	}
	s = fmt.Sprintf("%+v", delta)
	fmt.Println(s)
}
