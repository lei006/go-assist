package license

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
)

type LicenseKey struct {
	PubKey string `json:"pub_key"` //公钥
	PriKey string `json:"pri_key"` //私钥
}

func (key *LicenseKey) ToString() string {
	return "pub:" + key.PubKey + ",pri:" + key.PriKey
}

func (key *LicenseKey) ToJson() (string, error) {

	b, err := json.Marshal(key)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (key *LicenseKey) TestSign(msg string) (bool, error) {

	//生成数字签名
	//sign, err := ecc_tool.Sign(msg, pri_key)

	sign, err := key.Sign(msg)
	if err != nil {
		fmt.Println("对数据签名出错")
		return false, errors.New("生成签名出错:" + err.Error())
	}

	verify, err := key.Verify(msg, sign)

	return verify, err
}

func (key *LicenseKey) Sign(msg string) (string, error) {

	//取得私钥
	privateKey, err := PrivateKeyFromString(key.PriKey)
	if err != nil {
		return "", err
	}

	//计算哈希值
	hash := sha256.New()
	//填入数据
	hash.Write([]byte(msg))
	bytes := hash.Sum(nil)
	//对哈希值生成数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, bytes)
	if err != nil {
		return "", err
	}
	rtext, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	stext, err := s.MarshalText()
	if err != nil {
		return "", err
	}

	sign := &EccSign{
		Rtext:  string(rtext),
		Stext:  string(stext),
		PubKey: key.PubKey,
	}
	//sign.PubKey = "test"
	tmp_str, err := sign.ToJson()
	if err != nil {
		return "", err
	}

	return tmp_str, nil
}

//验证数字签名
func (key *LicenseKey) Verify(msg string, sign string) (bool, error) {

	//读取公钥
	publicKey, err := PublicKeyFromString(key.PubKey)
	if err != nil {
		return false, err
	}

	ecc_sign := EccSign{}
	err = ecc_sign.FromJson(sign)
	if err != nil {
		return false, err
	}
	//计算哈希值
	hash := sha256.New()
	hash.Write([]byte(msg))
	bytes := hash.Sum(nil)

	//验证数字签名
	var r, s big.Int
	r.UnmarshalText([]byte(ecc_sign.Rtext))
	s.UnmarshalText([]byte(ecc_sign.Stext))

	verify := ecdsa.Verify(publicKey, bytes, &r, &s)

	return verify, nil
}

//取得ECC私钥
func PrivateKeyFromString(private_str string) (*ecdsa.PrivateKey, error) {
	//读取私钥

	//pem解码
	block, _ := pem.Decode(([]byte)(private_str))
	//x509解码
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

//取得ECC公钥
func PublicKeyFromString(publish_str string) (*ecdsa.PublicKey, error) {
	//读取公钥
	//pem解密
	block, _ := pem.Decode(([]byte)(publish_str))

	//x509解密
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)
	return publicKey, nil
}
