package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cpusoft/goutil/belogs"
	"github.com/cpusoft/goutil/convert"
	"github.com/wenzhenxi/gorsa"
)

var Pubkey = `-----BEGIN Public key-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAk+89V7vpOj1rG6bTAKYM
56qmFLwNCBVDJ3MltVVtxVUUByqc5b6u909MmmrLBqS//PWC6zc3wZzU1+ayh8xb
UAEZuA3EjlPHIaFIVIz04RaW10+1xnby/RQE23tDqsv9a2jv/axjE/27b62nzvCW
eItu1kNQ3MGdcuqKjke+LKhQ7nWPRCOd/ffVqSuRvG0YfUEkOz/6UpsPr6vrI331
hWRB4DlYy8qFUmDsyvvExe4NjZWblXCqkEXRRAhi2SQRCl3teGuIHtDUxCskRIDi
aMD+Qt2Yp+Vvbz6hUiqIWSIH1BoHJer/JOq2/O6X3cmuppU4AdVNgy8Bq236iXvr
MQIDAQAB
-----END Public key-----
`

var Pirvatekey = `-----BEGIN Private key-----
MIIEpAIBAAKCAQEAk+89V7vpOj1rG6bTAKYM56qmFLwNCBVDJ3MltVVtxVUUByqc
5b6u909MmmrLBqS//PWC6zc3wZzU1+ayh8xbUAEZuA3EjlPHIaFIVIz04RaW10+1
xnby/RQE23tDqsv9a2jv/axjE/27b62nzvCWeItu1kNQ3MGdcuqKjke+LKhQ7nWP
RCOd/ffVqSuRvG0YfUEkOz/6UpsPr6vrI331hWRB4DlYy8qFUmDsyvvExe4NjZWb
lXCqkEXRRAhi2SQRCl3teGuIHtDUxCskRIDiaMD+Qt2Yp+Vvbz6hUiqIWSIH1BoH
Jer/JOq2/O6X3cmuppU4AdVNgy8Bq236iXvrMQIDAQABAoIBAQCCbxZvHMfvCeg+
YUD5+W63dMcq0QPMdLLZPbWpxMEclH8sMm5UQ2SRueGY5UBNg0WkC/R64BzRIS6p
jkcrZQu95rp+heUgeM3C4SmdIwtmyzwEa8uiSY7Fhbkiq/Rly6aN5eB0kmJpZfa1
6S9kTszdTFNVp9TMUAo7IIE6IheT1x0WcX7aOWVqp9MDXBHV5T0Tvt8vFrPTldFg
IuK45t3tr83tDcx53uC8cL5Ui8leWQjPh4BgdhJ3/MGTDWg+LW2vlAb4x+aLcDJM
CH6Rcb1b8hs9iLTDkdVw9KirYQH5mbACXZyDEaqj1I2KamJIU2qDuTnKxNoc96HY
2XMuSndhAoGBAMPwJuPuZqioJfNyS99x++ZTcVVwGRAbEvTvh6jPSGA0k3cYKgWR
NnssMkHBzZa0p3/NmSwWc7LiL8whEFUDAp2ntvfPVJ19Xvm71gNUyCQ/hojqIAXy
tsNT1gBUTCMtFZmAkUsjqdM/hUnJMM9zH+w4lt5QM2y/YkCThoI65BVbAoGBAMFI
GsIbnJDNhVap7HfWcYmGOlWgEEEchG6Uq6Lbai9T8c7xMSFc6DQiNMmQUAlgDaMV
b6izPK4KGQaXMFt5h7hekZgkbxCKBd9xsLM72bWhM/nd/HkZdHQqrNAPFhY6/S8C
IjRnRfdhsjBIA8K73yiUCsQlHAauGfPzdHET8ktjAoGAQdxeZi1DapuirhMUN9Zr
kr8nkE1uz0AafiRpmC+cp2Hk05pWvapTAtIXTo0jWu38g3QLcYtWdqGa6WWPxNOP
NIkkcmXJjmqO2yjtRg9gevazdSAlhXpRPpTWkSPEt+o2oXNa40PomK54UhYDhyeu
akuXQsD4mCw4jXZJN0suUZMCgYAgzpBcKjulCH19fFI69RdIdJQqPIUFyEViT7Hi
bsPTTLham+3u78oqLzQukmRDcx5ddCIDzIicMfKVf8whertivAqSfHytnf/pMW8A
vUPy5G3iF5/nHj76CNRUbHsfQtv+wqnzoyPpHZgVQeQBhcoXJSm+qV3cdGjLU6OM
HgqeaQKBgQCnmL5SX7GSAeB0rSNugPp2GezAQj0H4OCc8kNrHK8RUvXIU9B2zKA2
z/QUKFb1gIGcKxYr+LqQ25/+TGvINjuf6P3fVkHL0U8jOG0IqpPJXO3Vl9B8ewWL
cFQVB/nQfmaMa4ChK0QEUe+Mqi++MwgYbRHx1lIOXEfUJO+PXrMekw==
-----END Private key-----
`
var PrivateKey2 string = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1oUFD7aYYWh7p9HdNVYUM+R7zdkzsUMCM5EXkwxIV1FDR/3B
GQ5T1uAzSphdfqjgBxbfv9x5gmieTvKXW84f126VZBNXtIGXjfUDakOuGiqDdKl2
He9KBf9lfitwkLih9zYsWTuiBh8o7V3FFtM1i3DuzP4HCo57aUxtUy5b7ij8LQxZ
U0XwOXZHaszVl+nAOCL+0/+VmLJkyio3kiU9C0hJDd+j4tvJC3jyLipOYX50ID4G
PfZxIi7FCvCHKHBadfehlNcW9EtDpxPhGcyVpP7AsgvStQvyU0TZz7vhLHDfUtY/
nl7dlZ+5PWUmvuw4myHGzDmvQ6vcUVW7DRMGAQIDAQABAoIBAQCXklFrMtckLFEC
2LP2FaYcrFoVrlxp6TDLAr+ndMxAdfiWC2O+snLmpm9XS6Tz85qnJ7BcvglU7Vq9
6YaspU22SDpiBZC4x8Av22jYUo3XiyZq7bm5mPOynSw3I7Zbazl1lN9tBUeMD8Q5
Q0IYyI9SwS7ZxLtw6A+m7Qtp9J2b/ipURjVLyyW7UvIyq2qHfTA0YdN7FU8R3Y3X
MLHB0QGLUUTmnABmlYSbYh1hFvmXS6vDrkVUz1zgbbmOXF8bMI1fCyD1t45F15yD
uQUQ5dx9m1Z6rLw1Gjr/RgpngAX65oOlGtkYp+EYk1hGHAKmpGEr0vsskJc4zYI4
+/KAcW+BAoGBAOidUrKr/v/mOU2HINzHgVI+255nUgaGRezh6RMQnzwW6nz2KdZT
VkE/YTWCKh3zo5+i/uHddT/miT8jlPhBl2Te4fYPhBBZPXsOYcX3p2rLbsx96AD+
YxWmuQaIaRzk9zLgfBjOXTwL7lHpy8scTmW6gzskqRrELqZ1eMcUm2wpAoGBAOwV
/zv82AehZUJNkM9B4A5BtlgjTNUZinejRIJ8k6maeC2BIe8Gzb9Xj2xs6IBm7Ry9
DkPcalE+aySO+CYS6TtiKeSdJigzJyg7Kwq0C/M0GNZwP92ZvlKnV/4Patcl7+E2
VSnxwfSu3/We32vUd15yg5gpkqjyDnV/4USeCIYZAoGBALcoYRxkh5XRHl+oPbz5
rh8ndWAFtLWEdnyt6Qrk9Kyo0pvwbELhPbKEiDNMuYL5+2VQP2dzK8ZT7M91YfAU
HXQEd2F7GB6TVfCWA3CQrxdM9YI4xTw7EaPTsi6trC5fLzG1RqF1pD4Kmu2OrLPS
Jvy83mXsWObFgIH7T01aMYL5AoGAfnEZjetRWGTccrJQSHCjq38ORg5B7DANtR3A
Z5KJE2Ej1FtA7V/begtPSWba70ow3B91MGswleq0P5RC20FtoNxmS4bPFOCwrB9k
YgskC1FvrAnaarkY8fOmcO+Y7TnoS9ppqllM49t1H3vDdWEJvY/fYvOBFPLvQ4cG
A1YQgqECgYAFFan2RbUn7pUXyrRchPG9dknayEpaqK5a9CxmsAohJc1Gn5IwtWhQ
IV425nQHa4A2OlOo4D9mNVSVpufKqLfw/ld1kSd9fhD+wBgz+qvDy0JHXPU2W8Sy
fMaZ6pemtyzA9kmJ+2E1mBnjdY9rTTf4A3pEfNdrNziiI78fop6gAw==
-----END RSA PRIVATE KEY-----`

var PublicKey2 string = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1oUFD7aYYWh7p9HdNVYU
M+R7zdkzsUMCM5EXkwxIV1FDR/3BGQ5T1uAzSphdfqjgBxbfv9x5gmieTvKXW84f
126VZBNXtIGXjfUDakOuGiqDdKl2He9KBf9lfitwkLih9zYsWTuiBh8o7V3FFtM1
i3DuzP4HCo57aUxtUy5b7ij8LQxZU0XwOXZHaszVl+nAOCL+0/+VmLJkyio3kiU9
C0hJDd+j4tvJC3jyLipOYX50ID4GPfZxIi7FCvCHKHBadfehlNcW9EtDpxPhGcyV
pP7AsgvStQvyU0TZz7vhLHDfUtY/nl7dlZ+5PWUmvuw4myHGzDmvQ6vcUVW7DRMG
AQIDAQAB
-----END PUBLIC KEY-----`

func main() {
	// Public key encryption private key decryption
	if err := applyPubEPriD(); err != nil {
		fmt.Println("applyPubEPriD", err)
	}
	// Public key decryption private key encryption
	if err := applyPriEPubD(); err != nil {
		fmt.Println("applyPriEPubD", err)
	}
}

// Public key encryption private key decryption
func applyPubEPriD() error {
	data := `v=1&usr=zdns-rpki-f0c7d7c8&uniq=a8c3a2bd070533bbd22ad6d61a9bfbf17361e5377ffe3a3431048631bcabf76c&reg=2022-05-13T22:21:15+08:00&st=2022-05-13T22:21:15+08:00&ed=2122-04-19T22:21:15+08:00`
	grsa := gorsa.RSASecurity{}
	fmt.Println(PublicKey2)
	grsa.SetPublicKey(PublicKey2)

	rsadata, err := grsa.PubKeyENCTYPT([]byte(data))
	if err != nil {
		belogs.Error("RsaEncryptByPublicKey(): PubKeyENCTYPT fail:", err)
		return err
	}
	fmt.Println("applyPubEPriD(): PubKeyENCTYPT rsadata:", convert.PrintBytes(rsadata, 8))
	fmt.Println("rsadata:", hex.EncodeToString(rsadata))
	/*
		grsa := gorsa.RSASecurity{}
		grsa.SetPublicKey(PublicKey2)

		rsadata, err := grsa.PubKeyENCTYPT([]byte(data))
		if err != nil {
			fmt.Println("applyPubEPriD(): PubKeyENCTYPT:", err)
			return err
		}
		fmt.Println("applyPubEPriD(): PubKeyENCTYPT rsadata:", convert.PrintBytes(rsadata, 8))
	*/

	encryptData := `562bb3900a44cbad5d90cc6baa3fc1fe54fcf5b8bdf068accfcfd60a5d2802838a583ddd572407892620f2d28ac8d3fc5d3f98656599a504e1fc92eb7cd4f6602f74ad2460008d6eb2e39dbe78bafb8d162ebba35c7bf51a855e830dead6b0f3aac6220137698a066f8f78518c4e01e998fbe3418caa010e4e60dbfbb99b0f6d80279eb79c21aed7327558a02d679bef9f0a88461593615097074e8b8daef43e7e037af6219cb689f6b4ceae6936c36119a22464f83ba6fa462a968bdd71573b0324080975eb8c651a530aaa43f28fa28734f3a0f124a94c24746cd744951f1eeb17597e1bfaf97802ce49114d6570f7281e22c89656e5edae9d6f44cb3b2746`
	fmt.Println("encryptData:", encryptData)
	rsadata, err = hex.DecodeString(encryptData)
	if err != nil {
		fmt.Println("applyPubEPriD(): DecodeString:", err)
		return err
	}

	grsa2 := gorsa.RSASecurity{}
	err = grsa2.SetPrivateKey(PrivateKey2)
	if err != nil {
		fmt.Println("applyPubEPriD(): SetPrivateKey:", err)
		return err
	}

	data2, err := grsa2.PriKeyDECRYPT(rsadata)
	if err != nil {
		fmt.Println("applyPubEPriD(): PriKeyDECRYPT:", err)
		return err
	}
	fmt.Println("applyPubEPriD():PriKeyDECRYPT data2:", convert.PrintBytes(data2, 8), string(data2))
	/*
			pubenctypt, err := gorsa.PublicEncrypt(`hello world`, Pubkey)
			if err != nil {
				fmt.Println("applyPubEPriD() PublicEncrypt:", err)
				return err
			}
			fmt.Println("applyPubEPriD(): pubenctypt:", pubenctypt)
			pridecrypt, err := gorsa.PriKeyDecrypt(pubenctypt, Pirvatekey)
			if err != nil {
				fmt.Println("applyPubEPriD() PriKeyDecrypt:", err)
				return err
			}

		fmt.Println("applyPubEPriD(): pridecrypt:", pridecrypt)
		if string(pridecrypt) != `hello world` {
			return errors.New(`Decryption failed`)
		}
	*/
	return nil
}

// Public key decryption private key encryption
func applyPriEPubD() error {
	data := `hello world`
	grsa := gorsa.RSASecurity{}
	grsa.SetPrivateKey(PrivateKey2)

	rsadata, err := grsa.PriKeyENCTYPT([]byte(data))
	if err != nil {
		return err
	}
	fmt.Println("applyPriEPubD(): PriKeyEncrypt rsadata:", convert.PrintBytes(rsadata, 8))
	grsa2 := gorsa.RSASecurity{}
	err = grsa2.SetPublicKey(PublicKey2)
	if err != nil {
		return err
	}

	data2, err := grsa2.PubKeyDECRYPT(rsadata)
	if err != nil {
		return err
	}
	fmt.Println("applyPriEPubD():PubKeyDECRYPT data2:", convert.PrintBytes(data2, 8), string(data2))
	return nil
	/*
		prienctypt, err := gorsa.PriKeyEncrypt(`hello world`, Pirvatekey)
		if err != nil {
			fmt.Println("applyPriEPubD() PriKeyEncrypt:", err)
			return err
		}
		fmt.Println("applyPriEPubD(): PriKeyEncrypt rsadata:", convert.PrintBytes(prienctypt, 8))

		pubdecrypt, err := gorsa.PublicDecrypt(prienctypt, Pubkey)
		if err != nil {
			fmt.Println("applyPriEPubD() PublicDecrypt:", err)
			return err
		}
		fmt.Println("applyPubEPriD(): pubdecrypt:", pubdecrypt, string(pubdecrypt))
		if string(pubdecrypt) != `hello world` {
			return errors.New(`Decryption failed`)
		}
		return nil
	*/
}
