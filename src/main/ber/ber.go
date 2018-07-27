package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	//. "main/cert"
	"os"
	"reflect"
	"strconv"
)

type OidPacket struct {
	Oid          string
	OidPacket    *Packet
	ParentPacket *Packet
}

type Packet struct {
	ClassType   uint8
	TagType     uint8
	Tag         uint8
	Value       interface{}
	Data        *bytes.Buffer
	Children    []*Packet
	Description string
	Parent      *Packet
}

const (
	TagEOC              = 0x00
	TagBoolean          = 0x01
	TagInteger          = 0x02
	TagBitString        = 0x03
	TagOctetString      = 0x04
	TagNULL             = 0x05
	TagObjectIdentifier = 0x06
	TagObjectDescriptor = 0x07
	TagExternal         = 0x08
	TagRealFloat        = 0x09
	TagEnumerated       = 0x0a
	TagEmbeddedPDV      = 0x0b
	TagUTF8String       = 0x0c
	TagRelativeOID      = 0x0d
	TagSequence         = 0x10
	TagSet              = 0x11
	TagNumericString    = 0x12
	TagPrintableString  = 0x13
	TagT61String        = 0x14
	TagVideotexString   = 0x15
	TagIA5String        = 0x16
	TagUTCTime          = 0x17
	TagGeneralizedTime  = 0x18
	TagGraphicString    = 0x19
	TagVisibleString    = 0x1a
	TagGeneralString    = 0x1b
	TagUniversalString  = 0x1c
	TagCharacterString  = 0x1d
	TagBMPString        = 0x1e
	TagBitmask          = 0x1f // xxx11111b

	//private
	TagAsNum = 0xa0
	TagRdi   = 0xa1
)

var TagMap = map[uint8]string{
	TagEOC:              "EOC (End-of-Content)",
	TagBoolean:          "Boolean",
	TagInteger:          "Integer",
	TagBitString:        "Bit String",
	TagOctetString:      "Octet String",
	TagNULL:             "NULL",
	TagObjectIdentifier: "Object Identifier",
	TagObjectDescriptor: "Object Descriptor",
	TagExternal:         "External",
	TagRealFloat:        "Real (float)",
	TagEnumerated:       "Enumerated",
	TagEmbeddedPDV:      "Embedded PDV",
	TagUTF8String:       "UTF8 String",
	TagRelativeOID:      "Relative-OID",
	TagSequence:         "Sequence and Sequence of",
	TagSet:              "Set and Set OF",
	TagNumericString:    "Numeric String",
	TagPrintableString:  "Printable String",
	TagT61String:        "T61 String",
	TagVideotexString:   "Videotex String",
	TagIA5String:        "IA5 String",
	TagUTCTime:          "UTC Time",
	TagGeneralizedTime:  "Generalized Time",
	TagGraphicString:    "Graphic String",
	TagVisibleString:    "Visible String",
	TagGeneralString:    "General String",
	TagUniversalString:  "Universal String",
	TagCharacterString:  "Character String",
	TagBMPString:        "BMP String",
	//private
	TagAsNum: "ASNum",
	TagRdi:   "Rdi",
}

const (
	ClassUniversal   = 0   // 00xxxxxxb
	ClassApplication = 64  // 01xxxxxxb
	ClassContext     = 128 // 10xxxxxxb
	ClassPrivate     = 192 // 11xxxxxxb
	ClassBitmask     = 192 // 11xxxxxxb
)

var ClassMap = map[uint8]string{
	ClassUniversal:   "Universal",
	ClassApplication: "Application",
	ClassContext:     "Context",
	ClassPrivate:     "Private",
}

const (
	TypePrimative   = 0  // xx0xxxxxb
	TypeConstructed = 32 // xx1xxxxxb
	TypeBitmask     = 32 // xx1xxxxxb
)

var TypeMap = map[uint8]string{
	TypePrimative:   "Primative",
	TypeConstructed: "Constructed",
}

var Debug bool = false

func PrintBytes(buf []byte, indent string) {
	data_lines := make([]string, (len(buf)/30)+1)
	num_lines := make([]string, (len(buf)/30)+1)

	for i, b := range buf {
		data_lines[i/30] += fmt.Sprintf("%02x ", b)
		num_lines[i/30] += fmt.Sprintf("%02d ", (i+1)%100)
	}

	for i := 0; i < len(data_lines); i++ {
		fmt.Print(indent + data_lines[i] + "\n")
		//fmt.Print(indent + num_lines[i] + "\n\n")
	}
}

func PrintPacket(p *Packet) {
	printPacket(p, 0, false)
}

func printPacket(p *Packet, indent int, printBytes bool) {
	indent_str := ""
	for len(indent_str) != indent {
		indent_str += " "
	}

	class_str := ClassMap[p.ClassType]
	tagtype_str := TypeMap[p.TagType]
	tag_str := fmt.Sprintf("0x%02X", p.Tag)

	if p.ClassType == ClassUniversal {
		tag_str = TagMap[p.Tag]
	}

	value := fmt.Sprint(p.Value)
	description := ""
	if p.Description != "" {
		description = p.Description + ": "
	}

	fmt.Printf("%s%s(%s, %s, %s) Len=%d %q\n", indent_str, description, class_str, tagtype_str, tag_str, p.Data.Len(), value)

	if printBytes {
		PrintBytes(p.Bytes(), indent_str)
	}

	for _, child := range p.Children {
		fmt.Println("[children]-->")
		printPacket(child, indent+1, printBytes)
	}
}

func printPacketString(name string, p *Packet, printBytes bool, printChild bool) {
	class_str := ClassMap[p.ClassType]
	tagtype_str := TypeMap[p.TagType]
	tag_str := fmt.Sprintf("0x%02X", p.Tag)

	if p.ClassType == ClassUniversal {
		tag_str = TagMap[p.Tag]
	}

	value := fmt.Sprint(p.Value)
	description := ""
	if p.Description != "" {
		description = p.Description + ": "
	}

	fmt.Printf("\t%s  %s(%s, %s, %s) Len=%d %q\n", name, description, class_str, tagtype_str, tag_str, p.Data.Len(), value)

	if printBytes {
		PrintBytes(p.Bytes(), "    ")
	}
	if printChild {
		for _, child := range p.Children {
			fmt.Println("[children]-->")
			printPacketString(name+" --> children ", child, printBytes, printChild)
		}
	}
}

func resizeBuffer(in []byte, new_size uint64) (out []byte) {
	out = make([]byte, new_size)
	copy(out, in)
	return
}

func readBytes(reader io.Reader, buf []byte) error {
	idx := 0
	buflen := len(buf)
	for idx < buflen {
		n, err := reader.Read(buf[idx:])
		if err != nil {
			return err
		}
		idx += n
	}
	return nil
}

func ReadPacket(reader io.Reader) (*Packet, error) {
	buf := make([]byte, 2)
	err := readBytes(reader, buf)
	if err != nil {
		return nil, err
	}
	idx := uint64(2)
	datalen := uint64(buf[1])
	if Debug {
		fmt.Printf("Read: datalen = %d len(buf) = %d ", datalen, len(buf))
		for _, b := range buf {
			fmt.Printf("%02X ", b)
		}
		fmt.Printf("\n")
	}
	if datalen&128 != 0 {
		a := datalen - 128
		idx += a
		buf = resizeBuffer(buf, 2+a)
		err := readBytes(reader, buf[2:])
		if err != nil {
			return nil, err
		}
		datalen = DecodeInteger(buf[2 : 2+a])
		if Debug {
			fmt.Printf("Read: a = %d  idx = %d  datalen = %d  len(buf) = %d", a, idx, datalen, len(buf))
			for _, b := range buf {
				fmt.Printf("%02X ", b)
			}
			fmt.Printf("\n")
		}
	}

	buf = resizeBuffer(buf, idx+datalen)
	err = readBytes(reader, buf[idx:])
	if err != nil {
		return nil, err
	}

	if Debug {
		fmt.Printf("Read: len( buf ) = %d  idx=%d datalen=%d idx+datalen=%d\n", len(buf), idx, datalen, idx+datalen)
		for _, b := range buf {
			fmt.Printf("%02X ", b)
		}
	}

	p := DecodePacket(buf)
	return p, nil
}

func DecodeString(data []byte) (ret string) {
	for _, c := range data {
		ret += fmt.Sprintf("%c", c)
	}
	return
}
func DecodeUTF8String(data []byte) (ret string) {

	return string(data)
}

func DecodeIA5String(data []byte) (ret string) {
	return string(data)
}

func DecodeInteger(data []byte) (ret uint64) {
	for _, i := range data {
		ret = ret * 256
		ret = ret + uint64(i)
	}
	return
}
func DecodeUTCTime(data []byte) (ret string) {
	year := "20" + string(data[0:2])
	month := string(data[2:4])
	day := string(data[4:6])
	hour := string(data[6:8])
	minute := string(data[8:10])
	second := string(data[10:12])
	z := string(data[12])
	return year + "-" + month + "-" + day + " " + hour + ":" + minute + ":" + second + z
}
func DecodeGeneralizedTime(data []byte) (ret string) {
	year := string(data[0:4])
	month := string(data[4:6])
	day := string(data[6:8])
	hour := string(data[8:10])
	minute := string(data[10:12])
	second := string(data[12:14])
	z := string(data[14])
	return year + "-" + month + "-" + day + " " + hour + ":" + minute + ":" + second + z
}

func DecodeOid(data []byte) (ret string) {
	//printBytes("oid", data)
	oids := make([]uint32, len(data)+2)
	//第一个八位组采用公式：first_arc* 40+second_arc
	//后面的，当高位为1， 则表示需要加入下一位，合起来算的
	// https://msdn.microsoft.com/en-us/library/windows/desktop/bb540809(v=vs.85).aspx
	f := uint32(data[0])
	if f < 80 {
		oids[0] = f / 40
		oids[1] = f % 40
	} else {
		oids[0] = 2
		oids[1] = f - 80
	}
	var tmp uint32
	for i := 2; i <= len(data); i++ {
		f = uint32(data[i-1])
		//	fmt.Printf("f:0x%x\r\n", f)
		if f >= 0x80 {
			//		fmt.Printf("tmp<<8:0x%x +   (f&0x7f)0x%x\r\n", tmp<<8, (f & 0x7f))
			tmp = tmp<<7 + (f & 0x7f)
			//		fmt.Printf("tmp:0x%x\r\n", tmp)
		} else {
			oids[i] = tmp<<7 + (f & 0x7f)
			//		fmt.Printf("oids[i]:0x%x\r\n", oids[i])
			tmp = 0
		}
	}
	var buffer bytes.Buffer
	for i := 0; i < len(oids); i++ {
		if oids[i] == 0 {
			continue
		}
		buffer.WriteString(fmt.Sprint(oids[i]) + ".")
	}
	//fmt.Println(buffer.String()[0 : len(buffer.String())-1])
	return buffer.String()[0 : len(buffer.String())-1]
}
func EncodeInteger(val uint64) []byte {
	var out bytes.Buffer
	found := false
	shift := uint(56)
	mask := uint64(0xFF00000000000000)
	for mask > 0 {
		if !found && (val&mask != 0) {
			found = true
		}
		if found || (shift == 0) {
			out.Write([]byte{byte((val & mask) >> shift)})
		}
		shift -= 8
		mask = mask >> 8
	}
	return out.Bytes()
}

func DecodePacket(data []byte) *Packet {
	p, _, _ := decodePacket(data)
	return p
}

func decodePacket(data []byte) (*Packet, []byte, error) {
	if Debug {
		fmt.Printf("decodePacket: enter %d\n", len(data))
		printBytes("decodePacket: enter", data)
	}

	if len(data) < 2 {
		return nil, nil, errors.New("data is empty")
	}

	p := new(Packet)
	p.ClassType = data[0] & ClassBitmask
	p.TagType = data[0] & TypeBitmask
	p.Tag = data[0] & TagBitmask

	datalen := DecodeInteger(data[1:2])
	datapos := uint64(2)
	if datalen&128 != 0 {
		datalen -= 128
		datapos += datalen
		if 2+datalen > uint64(len(data)) {
			fmt.Println("data is less than 2+datalen")
			return nil, nil, errors.New("data is less than datalen")
		}
		datalen = DecodeInteger(data[2 : 2+datalen])
	}

	p.Data = new(bytes.Buffer)
	p.Children = make([]*Packet, 0, 2)
	p.Value = nil
	if datapos+datalen > uint64(len(data)) {
		fmt.Println(datapos, datalen, len(data))
		printBytes("data is less than datapos+datalen", data)
		return nil, nil, errors.New("data is less than datapos+datalen")
	}
	value_data := data[datapos : datapos+datalen]
	if Debug {
		fmt.Printf("decodePacket: p.ClassType=%d, p.TagType=%d, p.Tag=%d \n",
			p.ClassType, p.TagType, p.Tag)
		fmt.Println(datapos, datapos+datalen)
		printBytes("decodePacket:value_data ", value_data)
	}
	/*
				https://blog.csdn.net/liaowenfeng/article/details/8777595
				ASN.1字头，左边第0、1位，表示类型
		左边位0	位1	类别
			0	0	通用(Universal)
			0	1	应用(Application)
			1	0	上下文特定(Context Specific)
			1	1	专用(Private)
	*/
	if p.TagType == TypeConstructed {
		if Debug {
			fmt.Println("after p.TagType == TypeConstructed ")
		}
		for len(value_data) != 0 {
			var child *Packet
			var err error
			child, value_data, err = decodePacket(value_data)
			if err != nil {
				return nil, nil, err
			}
			p.AppendChild(child)
		}
	} else if p.ClassType == ClassUniversal {
		if Debug {
			printBytes("after p.ClassType == ClassUniversal:", data[datapos:datapos+datalen])
		}
		p.Data.Write(data[datapos : datapos+datalen])
		if Debug {
			printBytes("after p.Data.Write(data[datapos : datapos+datalen]) :", p.Bytes())
		}
		switch p.Tag {
		case TagEOC:
		case TagBoolean:
			val := DecodeInteger(value_data)
			p.Value = val != 0
		case TagInteger:
			p.Value = DecodeInteger(value_data)
		case TagBitString:
		case TagOctetString:
			//p.Value = DecodeString(value_data)
			// OctetString特殊，可能有子child，需要提取子child的第一位验证
			haveChild := false
			if len(value_data) > 0 {
				value_data_saved := data[datapos : datapos+datalen]
				childTagType := value_data[0] & TypeBitmask
				if Debug {
					fmt.Printf("decodePacket: childTagType %d, (%d);   value_date[0]=%d, (%d)\n", childTagType, TypeConstructed, value_data[0], TagBitmask)
					printBytes("before childTagType is TypeConstructed:", value_data)
				}
				if int(childTagType) == TypeConstructed {
					//var child *Packet
					if Debug {
						printBytes("TagOctetString before:", value_data)
					}
					child, _, err := decodePacket(value_data)
					if Debug {
						printBytes("TagOctetString before err==nil :", value_data_saved)
					}
					// 这里如果解析错误，说明不是child，而就是字符串，因此这里err不再往上返回，仅仅表示是原始字符串
					if err == nil {
						//return nil, nil, err
						if Debug {
							printBytes("TagOctetString after decodePacket :", p.Bytes())
						}
						// 这里要清空原来是bytes，设置新child的bytes
						p.Data.Reset()
						p.AppendChild(child)
						haveChild = true
					} else {

					}

				}

			}
			//如果没有children，则需要赋值bytes
			if !haveChild {
				p.Value = value_data
				/*
					var buf bytes.Buffer
					enc := gob.NewEncoder(&buf)
					err := enc.Encode(key)
				*/
			}
			break
		case TagNULL:
			p.Value = nil
		case TagObjectIdentifier:
			p.Value = DecodeOid(value_data)
			//fmt.Println(p.Value.(string))
		case TagObjectDescriptor:
		case TagExternal:
		case TagRealFloat:
		case TagEnumerated:
			p.Value = DecodeInteger(value_data)
		case TagEmbeddedPDV:
		case TagUTF8String:
			p.Value = DecodeUTF8String(value_data)
		case TagRelativeOID:
		case TagSequence:
		case TagSet:
		case TagNumericString:
		case TagPrintableString:
			p.Value = DecodeString(value_data)
		case TagT61String:
		case TagVideotexString:
		case TagIA5String:
			p.Value = DecodeIA5String(value_data)
		case TagUTCTime:
			p.Value = DecodeUTCTime(value_data)
		case TagGeneralizedTime:
			p.Value = DecodeGeneralizedTime(value_data)
		case TagGraphicString:
		case TagVisibleString:
		case TagGeneralString:
		case TagUniversalString:
		case TagCharacterString:
		case TagBMPString:
		//private
		case TagAsNum:
			p.Value = value_data
		case TagRdi:
			p.Value = value_data
		}
	} else {
		p.Data.Write(data[datapos : datapos+datalen])
	}
	if Debug {
		printBytes("decodePacket: end switch", data[datapos+datalen:])
	}
	return p, data[datapos+datalen:], nil
}

func (p *Packet) DataLength() uint64 {
	return uint64(p.Data.Len())
}

func (p *Packet) Bytes() []byte {
	var out bytes.Buffer
	out.Write([]byte{p.ClassType | p.TagType | p.Tag})
	packet_length := EncodeInteger(p.DataLength())
	if p.DataLength() > 127 || len(packet_length) > 1 {
		out.Write([]byte{byte(len(packet_length) | 128)})
		out.Write(packet_length)
	} else {
		out.Write(packet_length)
	}
	out.Write(p.Data.Bytes())
	return out.Bytes()
}

func (p *Packet) AppendChild(child *Packet) {
	p.Data.Write(child.Bytes())
	if len(p.Children) == cap(p.Children) {
		newChildren := make([]*Packet, cap(p.Children)*2)
		copy(newChildren, p.Children)
		p.Children = newChildren[0:len(p.Children)]
	}
	p.Children = p.Children[0 : len(p.Children)+1]
	p.Children[len(p.Children)-1] = child
}

func Encode(ClassType, TagType, Tag uint8, Value interface{}, Description string) *Packet {
	p := new(Packet)
	p.ClassType = ClassType
	p.TagType = TagType
	p.Tag = Tag
	p.Data = new(bytes.Buffer)
	p.Children = make([]*Packet, 0, 2)
	p.Value = Value
	p.Description = Description

	if Value != nil {
		v := reflect.ValueOf(Value)

		if ClassType == ClassUniversal {
			switch Tag {
			case TagOctetString:
				sv, ok := v.Interface().(string)
				if ok {
					p.Data.Write([]byte(sv))
				}
			}
		}
	}

	return p
}

func NewSequence(Description string) *Packet {
	return Encode(ClassUniversal, TypePrimative, TagSequence, nil, Description)
}

func NewBoolean(ClassType, TagType, Tag uint8, Value bool, Description string) *Packet {
	intValue := 0
	if Value {
		intValue = 1
	}

	p := Encode(ClassType, TagType, Tag, nil, Description)
	p.Value = Value
	p.Data.Write(EncodeInteger(uint64(intValue)))
	return p
}

func NewInteger(ClassType, TagType, Tag uint8, Value uint64, Description string) *Packet {
	p := Encode(ClassType, TagType, Tag, nil, Description)
	p.Value = Value
	p.Data.Write(EncodeInteger(Value))
	return p
}

func NewString(ClassType, TagType, Tag uint8, Value, Description string) *Packet {
	p := Encode(ClassType, TagType, Tag, nil, Description)
	p.Value = Value
	p.Data.Write([]byte(Value))
	return p
}
func parseMft(file string) ([]OidPacket, error) {

	f, _ := os.Open(file)
	b, _ := ioutil.ReadAll(f)
	pack := DecodePacket(b)
	/*
		//children := pack.Children
		classType := pack.ClassType
		description := pack.Description
		tag := pack.Tag
		tagType := pack.TagType
		value := fmt.Sprint(pack.Value)
		fmt.Printf("classType:%v\r\n", classType)
		fmt.Printf("description:%v\r\n", description)
		fmt.Printf("tag:%v\r\n", tag)
		fmt.Printf("tagType:%v\r\n", tagType)
		fmt.Printf("value:%v\r\n", value)
		//fmt.Printf("classType:%v\r\n", children)
		printPacket(pack, 4, true)
	*/
	/*
		oidPackets := make(map[string]Packet, 20)
		addParent(pack, oidPackets)
		fmt.Println(len(oidPackets))
		for oid, packet := range oidPackets {
			fmt.Println(oid)
			PrintBytes(packet.Bytes(), "    ")

		}
	*/
	//oidPacketss := make([]OidPacket, 10)
	var oidPacketss = &[]OidPacket{}
	transformPacket(pack, oidPacketss)
	fmt.Println("all oidPacket size:", len(*oidPacketss))
	for _, oidPacket := range *oidPacketss {

		if Debug {
			fmt.Println(oidPacket.Oid)
			printBytes("oid parent bytes:", oidPacket.ParentPacket.Bytes())
			printBytes("oid self bytes:", oidPacket.OidPacket.Bytes())
			fmt.Println("")
		}
	}

	printPacketString("all packet", pack, true, true)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++")
	return *oidPacketss, nil
}

func transformPacket(p *Packet, oidPackets *[]OidPacket) {

	for i, _ := range p.Children {

		p.Children[i].Parent = p
		//fmt.Println(p.Children[i].Tag, TagObjectIdentifier)
		if p.Children[i].Tag == TagObjectIdentifier {
			oidPacket := OidPacket{}
			//fmt.Printf("%s%s(%s, %s, %s) Len=%d %q\n", indent_str, description, class_str, tagtype_str, tag_str, p.Data.Len(), value)
			//fmt.Println(p.Children[i].Value.(string))
			oid := fmt.Sprint(p.Children[i].Value)
			oidPacket.Oid = oid
			//fmt.Println("addParent():oid:", oid)
			oidPacket.ParentPacket = p
			oidPacket.OidPacket = p.Children[i]
			(*oidPackets) = append((*oidPackets), oidPacket)
		}
		transformPacket(p.Children[i], oidPackets)
	}
}

func addParent(p *Packet, oidPackets map[string]Packet) {

	for i, _ := range p.Children {

		p.Children[i].Parent = p
		//fmt.Println(p.Children[i].Tag, TagObjectIdentifier)
		if p.Children[i].Tag == TagObjectIdentifier {

			//fmt.Printf("%s%s(%s, %s, %s) Len=%d %q\n", indent_str, description, class_str, tagtype_str, tag_str, p.Data.Len(), value)
			//fmt.Println(p.Children[i].Value.(string))
			oid := fmt.Sprint(p.Children[i].Value)
			//fmt.Println("addParent():oid:", oid)
			if _, ok := oidPackets[oid]; ok {
				printBytes("double oid:"+oid, (*p).Bytes())
			} else {
				oidPackets[oid] = *p
			}

			if oid == "2.16.840.1.101.3.4.2.1" {
				printBytes("double 2.16.840.1.101.3.4.2.1 parent bytes:"+oid, (*p).Bytes())
				printBytes("double 2.16.840.1.101.3.4.2.1 slef bytes:"+oid, p.Children[i].Bytes())
			}

		}
		addParent(p.Children[i], oidPackets)
	}
}

const (
	ipv4    = 0x01
	ipv6    = 0x02
	ipv4len = 32
	ipv6len = 128
)

func decodeAddressPrefix(addressPrefixPacket *Packet, ipType int) error {
	addressPrefix := addressPrefixPacket.Bytes()
	addressShouldLen, _ := strconv.Atoi(fmt.Sprintf("%d", addressPrefix[1]))
	unusedBitLen, _ := strconv.Atoi(fmt.Sprintf("%d", addressPrefix[2]))

	address := addressPrefix[3:]
	ipAddress := ""

	if ipType == ipv4 {
		// ipv4 的CIDR 表示法
		prefix := ipv4len - 8*(addressShouldLen-1) - unusedBitLen
		if Debug {
			fmt.Println(fmt.Sprintf("prefix := ipv4len - 8*(addressShouldLen-1) - unusedBitLen:  %d := %d - 8 *(%d-1)-  %d \r\n",
				prefix, ipv4len, addressShouldLen, unusedBitLen))
		}
		//printBytes("address:", address)

		ipv4Address := ""
		for i := 0; i < len(address); i++ {
			ipv4Address += fmt.Sprintf("%d", address[i])
			if i < len(address)-1 {
				ipv4Address += "."
			}
		}
		ipv4Address += "/" + fmt.Sprintf("%d", prefix)
		ipAddress = ipv4Address

	} else if ipType == ipv6 {
		// ipv6的前缀表示法，和ipv4不一样
		prefix := 8*(addressShouldLen-1) - unusedBitLen
		if Debug {
			fmt.Println(fmt.Sprintf("prefix :=  8*(addressShouldLen-1) - unusedBitLen:  %d := 8 *(%d-1)-  %d \r\n",
				prefix, addressShouldLen, unusedBitLen))
		}

		//printBytes("address:", address)

		ipv6Address := ""
		for i := 0; i < len(address); i++ {
			ipv6Address += fmt.Sprintf("%02x", address[i])
			if i%2 == 1 && i < len(address)-1 {
				ipv6Address += ":"
			}
		}
		//补齐位数
		if len(address)%2 == 1 {
			ipv6Address += "00"
		}
		ipv6Address += "/" + fmt.Sprintf("%d", prefix)
		ipAddress = ipv6Address

	}
	addressPrefixPacket.Value = ipAddress
	return nil
}

func printAsn(name string, typ byte, ln byte, byt []byte) {
	fmt.Println(fmt.Sprintf(name+"Type:0x%02x (%d)", typ, typ))
	fmt.Println(fmt.Sprintf(name+"Len:0x%02x (%d)", ln, ln))
	printBytes(name+"Value:", byt)
}

func printBytes(name string, byt []byte) {
	fmt.Println(name)
	for _, i := range byt {
		fmt.Print(fmt.Sprintf("0x%02x ", i))
	}
	fmt.Println("")
}
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func main() {

	file := `E:\Go\go-study\src\main\cert\41870XBX5RmmOBSWl-AwgOrYdys.mft`
	//file := `E:\Go\go-study\src\main\cert\H.cer`
	oidPackets, err := parseMft(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	//var oidIpAddressStr string
	//var oidASStr string

	oidIpAddressKey := "1.3.6.1.5.5.7.1.7"
	var ipType int
	oidASKey := "1.3.6.1.5.5.7.1.8"
	manifestKey := "1.2.840.113549.1.9.16.1.26"

	for _, oidPacket := range oidPackets {
		if oidPacket.Oid == oidIpAddressKey {
			if len(oidPacket.ParentPacket.Children) > 1 {
				critical := oidPacket.ParentPacket.Children[1]
				printPacketString("critical", critical, true, false)

				extnValue := oidPacket.ParentPacket.Children[2]
				if len(extnValue.Children) > 0 {
					for _, IpAddressBlocks := range extnValue.Children {
						if len(IpAddressBlocks.Children) > 0 {
							for _, IPAddressFamily := range IpAddressBlocks.Children {
								if len(IPAddressFamily.Children) > 0 {
									addressFamily := IPAddressFamily.Children[0]
									printPacketString("addressFamily", addressFamily, true, false)

									addressFamilyBytes := addressFamily.Value.([]byte)
									if addressFamilyBytes[1] == ipv4 {
										ipType = ipv4
									} else if addressFamilyBytes[1] == ipv6 {
										ipType = ipv6
									} else {
										fmt.Println("error iptype")
										return
									}
									if Debug {
										printBytes(fmt.Sprintf("addressFamilyBytes: iptype: %d ", ipType), addressFamilyBytes)
									}
									IPAddressChoice := IPAddressFamily.Children[1]
									if Debug {
										printPacketString("IPAddressChoice", IPAddressChoice, true, false)
									}
									if len(IPAddressChoice.Children) > 0 {
										for _, addressesOrRanges := range IPAddressChoice.Children {
											if Debug {
												printPacketString("addressesOrRanges", addressesOrRanges, true, false)
												fmt.Println("addressesOrRanges: len: ", len(addressesOrRanges.Children))
											}
											if len(addressesOrRanges.Children) > 0 {

												min := addressesOrRanges.Children[0]
												max := addressesOrRanges.Children[1]
												decodeAddressPrefix(min, ipType)
												decodeAddressPrefix(max, ipType)
												printPacketString("Range min", min, true, false)
												printPacketString("Range max", max, true, false)

											} else {
												decodeAddressPrefix(addressesOrRanges, ipType)
												printPacketString("addresses", addressesOrRanges, true, false)
											}
										}
									} else {
										inherit := IPAddressChoice.Value.([]byte)
										printBytes("inherit from issuer is NULL 2 ", inherit)
									}
								}
							}
						}
					}
				}
			}
		}
		if oidPacket.Oid == oidASKey {
			if len(oidPacket.ParentPacket.Children) > 1 {
				critical := oidPacket.ParentPacket.Children[1]
				printPacketString("critical", critical, true, false)

				extnValue := oidPacket.ParentPacket.Children[2]
				if len(extnValue.Children) > 0 {
					for _, ASIdentifiers := range extnValue.Children {

						if len(ASIdentifiers.Children) > 0 {
							for _, ASIdentifier := range ASIdentifiers.Children {

								if len(ASIdentifier.Children) > 0 {
									for _, asIdsOrRanges := range ASIdentifier.Children {
										if len(asIdsOrRanges.Children) > 0 {
											for _, ASIdOrRange := range asIdsOrRanges.Children {

												//区分两种：一种children是ASRange，一种是ASId
												if len(ASIdOrRange.Children) > 1 {
													min := ASIdOrRange.Children[0]
													max := ASIdOrRange.Children[1]

													printPacketString("ASNum min", min, true, false)
													printPacketString("ASNum max", max, true, false)
												} else {
													printPacketString("ASId", ASIdOrRange, true, false)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if oidPacket.Oid == manifestKey {

			if len(oidPacket.ParentPacket.Children) > 1 {
				seq0 := oidPacket.ParentPacket.Children[1]
				if len(seq0.Children) > 0 {
					octPacket := seq0.Children[0]
					if len(octPacket.Children) > 0 {
						secPacket := octPacket.Children[0]
						if len(secPacket.Children) > 0 {
							manifestNumber := secPacket.Children[0]
							printPacketString("manifestNumber", manifestNumber, true, false)

							thisUpdate := secPacket.Children[1]
							printPacketString("thisUpdate", thisUpdate, true, false)

							nextUpdate := secPacket.Children[2]
							printPacketString("nextUpdate", nextUpdate, true, false)

							fileHashAlg := secPacket.Children[3]
							printPacketString("fileHashAlg", fileHashAlg, true, false)

							fileList := secPacket.Children[4]
							if len(fileList.Children) > 0 {
								for _, fileAndHash := range fileList.Children {
									if len(fileAndHash.Children) > 1 {
										file := fileAndHash.Children[0]
										printPacketString("file", file, true, false)

										hash := fileAndHash.Children[1]
										printPacketString("hash", hash, true, false)
									}
								}
							}

						}
					}

				}
			}
		}

	}
}
