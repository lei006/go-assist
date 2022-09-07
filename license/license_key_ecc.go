package license

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func MakEccP521Key() (*LicenseKey, error) {

	ecc_key := &LicenseKey{}

	pub, pri, err := generateECCKeyString(elliptic.P521())
	if err != nil {
		return nil, err
	}
	ecc_key.PriKey = pri
	ecc_key.PubKey = pub

	return ecc_key, nil
}

func generateECCKeyString(curve elliptic.Curve) (string, string, error) {

	//生成密钥对
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", "", err
	}

	//x509编码
	eccPrivateKey, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}
	//pem编码
	privateBlock := pem.Block{
		Type:  "ecc private key",
		Bytes: eccPrivateKey,
	}
	tmp := pem.EncodeToMemory(&privateBlock)
	if tmp == nil {
		return "", "", errors.New("error")
	}

	private_str := string(tmp)

	//保存公钥
	publicKey := privateKey.PublicKey

	//x509编码
	eccPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return "", "", err
	}

	//pem编码
	block := pem.Block{Type: "ecc public key", Bytes: eccPublicKey}
	tmp = pem.EncodeToMemory(&block)
	if tmp == nil {
		return "", "", errors.New("error")
	}
	//公钥字符串
	publish_str := string(tmp)

	return publish_str, private_str, nil
}
