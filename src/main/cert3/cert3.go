package main

import ()

func main() {

	//https://blog.csdn.net/Zhymax/article/details/7683925
	// 比较底层的，能识别ROA各个键值
	//openssl asn1parse -in privatekey.der -inform DER

	// 比较上层，只能cer、crl
	//openssl x509 -noout -text -in 2.cer --inform der

}
