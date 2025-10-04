package license

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
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

func MakeLicenseKey(pub_key string, pri_key string) *LicenseKey {
	key := &LicenseKey{
		PubKey: pub_key,
		PriKey: pri_key,
	}
	return key
}

func (key *LicenseKey) ToCompare(tmp_key *LicenseKey) bool {

	if tmp_key.PubKey == key.PubKey && tmp_key.PriKey == key.PriKey {
		return true
	}
	return false
}

//签名一个数具
func (key *LicenseKey) ToJson() (string, error) {

	b, err := json.Marshal(key)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (key *LicenseKey) FromJson(data_str string) error {

	data := []byte(data_str)
	err := json.Unmarshal(data, key)
	return err

}

func (key *LicenseKey) ToBase64String() (string, error) {

	str_json, err := key.ToJson()
	if err != nil {
		return "", err
	}

	base64_str := base64.StdEncoding.EncodeToString([]byte(str_json))

	return base64_str, nil
}

func (key *LicenseKey) FromBase64String(data_str string) error {

	byte_data, err := base64.StdEncoding.DecodeString(data_str)
	if err != nil {
		return err
	}
	err = key.FromJson(string(byte_data))
	if err != nil {
		return err
	}

	return nil
}

//取得公钥字符串
func (key *LicenseKey) GetBase64PubKey() string {
	return base64.StdEncoding.EncodeToString([]byte(key.PubKey))
}

//设置公钥字符串
func (key *LicenseKey) SetBase64PubKey(base64_key_str string) error {
	byte_data, err := base64.StdEncoding.DecodeString(base64_key_str)
	if err != nil {
		return err
	}
	key.PubKey = string(byte_data)
	return nil
}

func (key *LicenseKey) TestKey() (bool, error) {

	return false, nil
}

func (key *LicenseKey) TestSign(msg string) (bool, error) {

	sign, err := key.Sign(msg)
	if err != nil {
		fmt.Println("对数据签名出错")
		return false, errors.New("生成签名出错:" + err.Error())
	}

	verify, err := key.VerifySign(msg, sign)

	return verify, err
}

func (key *LicenseKey) Sign(msg string) (string, error) {

	//取得私钥
	privateKey, err := key.privateKeyFromString(key.PriKey)
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
		Rtext: string(rtext),
		Stext: string(stext),
	}
	//sign.PubKey = "test"
	tmp_str, err := sign.ToJson()
	if err != nil {
		return "", err
	}

	return tmp_str, nil
}

//验证数字签名
func (key *LicenseKey) VerifySign(msg string, sign string) (bool, error) {

	//读取公钥
	publicKey, err := key.publicKeyFromString(key.PubKey)
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
func (key *LicenseKey) privateKeyFromString(private_str string) (*ecdsa.PrivateKey, error) {
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
func (key *LicenseKey) publicKeyFromString(publish_str string) (*ecdsa.PublicKey, error) {
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
