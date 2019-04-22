package main

import (
	"encoding/asn1"
	"encoding/hex"
	_ "errors"
	"fmt"
	"time"

	convert "github.com/cpusoft/goutil/convert"
	jsonutil "github.com/cpusoft/goutil/jsonutil"
)

type ManifestParseBigMftNumber struct {
	ManifestNumber asn1.RawValue         `json:"manifestNumber"`
	ThisUpdate     time.Time             `asn1:"generalized" json:"thisUpdate"`
	NextUpdate     time.Time             `asn1:"generalized" json:"nextUpdate"`
	FileHashAlg    asn1.ObjectIdentifier `json:"fileHashAlg"`
	FileList       []FileAndHashParse    `json:"fileList"`
}
type ManifestParse struct {
	ManifestNumber int64                 `json:"manifestNumber"`
	ThisUpdate     time.Time             `asn1:"generalized" json:"thisUpdate"`
	NextUpdate     time.Time             `asn1:"generalized" json:"nextUpdate"`
	FileHashAlg    asn1.ObjectIdentifier `json:"fileHashAlg"`
	FileList       []FileAndHashParse    `json:"fileList"`
}
type ManifestParse2 struct {
	FileList []FileAndHashParse `json:"fileList"`
}

type FileAndHashParse struct {
	File string         `asn1:"ia5" json:"file"`
	Hash asn1.BitString `json:"hash"`
}

func main() {

	mftHex := `3082040902020155180F32303139303332313031343130345A180F32303139303332323031343130345A0609608648016503040201308203D43044161F456542796F456C75626C7471786B5F77455676705A4C30534563342E63726C032100A3009FE52F79CD7E1CC658CC2DAC632FADEADEE2DF30D97E59893EB91CEBC0543044161F534439734A784276624177414737514E4F4673633937436649344D2E726F610321006FE255C875B1DB3F398754E53103CB07AECBBD135621DF85D17FBBA8E6AB57233044161F5A386A723949366F41785263426F7667666241433659434C7A6B6F2E726F6103210039C2CE5FA0887AE26A8B2DC54CB0D15B0BC5E4BCD23D18BF70C598FA4D8C02C33044161F637A324452566A696255394B314F794A49794137335F35374E6B512E726F610321000D689AB70940AA9FCF9D792E4E1FC428F1D534216FF7312798CA01AB55E2AE3B3044161F6A4B4B513535464B4E55625A645568555F534532616C356B626F382E726F61032100CE868B4DE404F53FB5071B8978D51A1A8A160D61C48C10AD11355643595EA5383044161F6A6B785956756B5838616D326976456336667535376B426B7136592E726F610321001ECD8B4D30A781C1632FBD159138875640F3527BA1A4E5FFB316413089D7DC283044161F6B32414C5542393431426F61426F584C2D78565043383951617A632E726F61032100758343300B538E609016C31F7822B6C0C5175A4A3F1081D8F6FF767E1187EDB43044161F6C34493038526A77736B506B50494E5F316C585A7431655A7146452E726F610321001996924AEC025E126206668A07E931B970489C220B33CF3043F718C88A1EDD053044161F6C4E7A4142364F4B427450725F30713742324A4535544141524A412E726F6103210080CA946B49D96CBCD5EFD1503A689E400D47E7CF2DC297E42202099BE066C5773044161F6F6D4A78644D4D5665465373796C5279334C447446594C444B78592E726F610321006E25DCC5D4578754EA586BA0C3A6287FD81E767EFC90031196B91B6DA22B32503044161F7446587466735963496259757750584C4E6A5A36467942547462452E726F61032100CE1683D262B6279D1547CAC3ED90FF3ADF3AA036F6703434B06AA7367D86D38E3044161F745A41696B77415976307932646F58534C54494A47476F4B5533412E726F6103210012DCFB9631B9C6287BA77EBA8760DA1E6E0745029F269E6A2DCDB2E348122A3D3044161F776B3779476A516B525F4B3766735159724A5A6C6E6D4C526D55592E726F61032100D8B07A8D0283FF64C94D8FD0778D91714BE953DD15E89E42F453F8DE60A6E5663044161F78497572665570675730657770686855515F417779356F627657772E72`
	mftByte, _ := hex.DecodeString(mftHex)

	PrintBytes("mftByte:", mftByte)

	seqType := mftByte[0]
	seqLen0 := mftByte[1]
	fmt.Println("seqType:", seqType, "      seqLen0:", seqLen0)
	var seqContent []byte
	var seqLenLen int
	if seqLen0&0x80 != 0 {
		seqLenLen = int(seqLen0 &^ 0x80)
		seqContent = mftByte[2+seqLenLen:]
	} else {
		seqLenLen = 1
		seqContent = mftByte[2:]
	}
	PrintBytes("seqContent:", seqContent)

	mftNumType := seqContent[0]
	mftNumLen0 := seqContent[1]
	fmt.Println("mftNumType:", mftNumType, "      mftNumLen0:", mftNumLen0)
	var mftNumContent []byte
	var mftNumLenLen int
	if mftNumLen0&0x80 != 0 {
		mftNumLenLen = int(mftNumLen0 &^ 0x80)
		mftNumContent = seqContent[2+mftNumLenLen:]
	} else {
		mftNumLenLen = int(mftNumLen0)
		mftNumContent = seqContent[2:]
	}
	PrintBytes("mftNumContent:", mftNumContent)

	mftNum := mftNumContent[:mftNumLenLen]
	PrintBytes("mftNum:", mftNum)

	// time类型长度固定0x18 0x0f ** ** **, 总长17
	timeStart := mftNumContent[mftNumLenLen:]
	PrintBytes("timeStart:", timeStart)

	time1 := timeStart[2:0x0f]
	PrintBytes("time1:", time1)

	time2 := timeStart[2+0x0f+2 : 2+0x0f+0x0f]
	PrintBytes("time2:", time2)

	OIDStart := timeStart[2+0x0f+2+0x0f:]
	PrintBytes("OIDStart:", OIDStart)

	oid := OIDStart[2 : 2+OIDStart[1]]
	PrintBytes("oid:", oid)

	seq1Start := OIDStart[2+OIDStart[1]:]
	PrintBytes("seq1Start:", seq1Start)

	mftBig := ManifestParse2{}
	rest, err := asn1.Unmarshal(seq1Start, &mftBig)
	fmt.Println("ManifestParse2:", mftBig, "   rest:", rest, " err:", err)

	seqLen0 = seq1Start[1]
	var seqLen uint64
	if seqLen0&0x80 != 0 {
		seqLenLen = int(seqLen0 &^ 0x80)
		seqLenB := seq1Start[2 : 2+seqLenLen]
		seqLen = convert.Bytes2Uint64(seqLenB)
		seqContent = seq1Start[2+seqLenLen:]
	} else {
		seqLenLen = 1
		seqLenB := seq1Start[2:2]
		seqLen = convert.Bytes2Uint64(seqLenB)
		seqContent = seq1Start[2:]
	}
	fmt.Println("seqLen:", seqLen, "  len(seqContent):", len(seqContent))
	PrintBytes("seqContent:", seqContent)

	fileList := FileAndHashParse{}
	for len(seqContent) > 0 {
		seqContent, err = asn1.Unmarshal(seqContent, &fileList)
		fmt.Println("rest:", rest, " err:", err, "   fileList:", jsonutil.MarshalJson(fileList))
	}
	//	var xxx uint64
	//	xxx = 0
	//	for i, _ := range mftLenByte {
	//		xxx = xxx<<8 + uint64(mftLenByte[i])
	//	}
	//	fmt.Println("xxx:", xxx)
	//	mftParse := mftByte[2+xxx:]
	//fmt.Println("mftParse:", mftParse)

}

func PrintBytes(name string, buf []byte) {
	fmt.Println(name)
	data_lines := make([]string, (len(buf)/30)+1)

	for i, b := range buf {
		data_lines[i/30] += fmt.Sprintf("%02x ", b)
	}

	for i := 0; i < len(data_lines); i++ {
		fmt.Println(data_lines[i])
	}
}
