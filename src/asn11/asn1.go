package main

import (
	"fmt"
	"reflect"

	"github.com/PromonLogicalis/asn1"
)

func main() {
	ctx := asn1.NewContext()

	// Use BER for encoding and decoding.
	ctx.SetDer(false, false)

	// Add a CHOICE
	ctx.AddChoice("value", []asn1.Choice{
		{
			Type:    reflect.TypeOf(""),
			Options: "tag:0",
		},
		{
			Type:    reflect.TypeOf(int(0)),
			Options: "tag:1",
		},
	})

	type Message struct {
		Id    int
		Value interface{} `asn1:"choice:value"`
	}

	// Encode
	/*
		msg := Message{
			Id:    1000,
			Value: "this is a value",
		}
	*/
	msg := Message{
		Id:    1000,
		Value: 999,
	}
	data, err := ctx.Encode(msg)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range data {
		fmt.Print(fmt.Sprintf("0x%02x ", d))
	}
	// Decode
	decodedMsg := Message{}
	_, err = ctx.Decode(data, &decodedMsg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("%+v\n", decodedMsg)
}
