package main

import (
	"github.com/cpusoft/goutil/fileutil"
)

var b = []byte{
	0x30, 0x82, 0x02, 0x18, 0x30,
	0x82, 0x01, 0xc2, 0x02, 0x09, 0x00, 0x8c, 0xc3, 0x37, 0x92, 0x10, 0xec, 0x2c,
	0x98, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01,
	0x05, 0x05, 0x00, 0x30, 0x81, 0x92, 0x31, 0x0b, 0x30, 0x09, 0x06, 0x03, 0x55,
	0x04, 0x06, 0x13, 0x02, 0x58, 0x58, 0x31, 0x13, 0x30, 0x11, 0x06, 0x03, 0x55,
	0x04, 0x08, 0x13, 0x0a, 0x53, 0x6f, 0x6d, 0x65, 0x2d, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x31, 0x0d, 0x30, 0x0b, 0x06, 0x03, 0x55, 0x04, 0x07, 0x13, 0x04, 0x43,
	0x69, 0x74, 0x79, 0x31, 0x21, 0x30, 0x1f, 0x06, 0x03, 0x55, 0x04, 0x0a, 0x13,
	0x18, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x20, 0x57, 0x69, 0x64,
	0x67, 0x69, 0x74, 0x73, 0x20, 0x50, 0x74, 0x79, 0x20, 0x4c, 0x74, 0x64, 0x31,
	0x1a, 0x30, 0x18, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13, 0x11, 0x66, 0x61, 0x6c,
	0x73, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f,
	0x6d, 0x31, 0x20, 0x30, 0x1e, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d,
	0x01, 0x09, 0x01, 0x16, 0x11, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x40, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x30, 0x1e, 0x17, 0x0d,
	0x30, 0x39, 0x31, 0x30, 0x30, 0x38, 0x30, 0x30, 0x32, 0x35, 0x35, 0x33, 0x5a,
	0x17, 0x0d, 0x31, 0x30, 0x31, 0x30, 0x30, 0x38, 0x30, 0x30, 0x32, 0x35, 0x35,
	0x33, 0x5a, 0x30, 0x81, 0x92, 0x31, 0x0b, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04,
	0x06, 0x13, 0x02, 0x58, 0x58, 0x31, 0x13, 0x30, 0x11, 0x06, 0x03, 0x55, 0x04,
	0x08, 0x13, 0x0a, 0x53, 0x6f, 0x6d, 0x65, 0x2d, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x31, 0x0d, 0x30, 0x0b, 0x06, 0x03, 0x55, 0x04, 0x07, 0x13, 0x04, 0x43, 0x69,
	0x74, 0x79, 0x31, 0x21, 0x30, 0x1f, 0x06, 0x03, 0x55, 0x04, 0x0a, 0x13, 0x18,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x20, 0x57, 0x69, 0x64, 0x67,
	0x69, 0x74, 0x73, 0x20, 0x50, 0x74, 0x79, 0x20, 0x4c, 0x74, 0x64, 0x31, 0x1a,
	0x30, 0x18, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13, 0x11, 0x66, 0x61, 0x6c, 0x73,
	0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	0x31, 0x20, 0x30, 0x1e, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01,
	0x09, 0x01, 0x16, 0x11, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x40, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x30, 0x5c, 0x30, 0x0d, 0x06,
	0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x01, 0x05, 0x00, 0x03,
	0x4b, 0x00, 0x30, 0x48, 0x02, 0x41, 0x00, 0xcd, 0xb7, 0x63, 0x9c, 0x32, 0x78,
	0xf0, 0x06, 0xaa, 0x27, 0x7f, 0x6e, 0xaf, 0x42, 0x90, 0x2b, 0x59, 0x2d, 0x8c,
	0xbc, 0xbe, 0x38, 0xa1, 0xc9, 0x2b, 0xa4, 0x69, 0x5a, 0x33, 0x1b, 0x1d, 0xea,
	0xde, 0xad, 0xd8, 0xe9, 0xa5, 0xc2, 0x7e, 0x8c, 0x4c, 0x2f, 0xd0, 0xa8, 0x88,
	0x96, 0x57, 0x72, 0x2a, 0x4f, 0x2a, 0xf7, 0x58, 0x9c, 0xf2, 0xc7, 0x70, 0x45,
	0xdc, 0x8f, 0xde, 0xec, 0x35, 0x7d, 0x02, 0x03, 0x01, 0x00, 0x01, 0x30, 0x0d,
	0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x05, 0x05, 0x00,
	0x03, 0x41, 0x00, 0xa6, 0x7b, 0x06, 0xec, 0x5e, 0xce, 0x92, 0x77, 0x2c, 0xa4,
	0x13, 0xcb, 0xa3, 0xca, 0x12, 0x56, 0x8f, 0xdc, 0x6c, 0x7b, 0x45, 0x11, 0xcd,
	0x40, 0xa7, 0xf6, 0x59, 0x98, 0x04, 0x02, 0xdf, 0x2b, 0x99, 0x8b, 0xb9, 0xa4,
	0xa8, 0xcb, 0xeb, 0x34, 0xc0, 0xf0, 0xa7, 0x8c, 0xf8, 0xd9, 0x1e, 0xde, 0x14,
	0xa5, 0xed, 0x76, 0xbf, 0x11, 0x6f, 0xe3, 0x60, 0xaa, 0xfa, 0x88, 0x21, 0x49,
	0x04, 0x35,
}

func main() {
	file := `E:\Go\go-study\src\file2\1.cer`
	fileutil.WriteBytesToFile(file, b)
}
