package ecc_tool

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
)

type EccSign struct {
	Rtext string `json:"rtext"`
	Stext string `json:"stext"`
}

func (this *EccSign) ToJson() (string, error) {
	b, err := json.Marshal(this)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (this *EccSign) FromJson(data_str string) error {
	data := []byte(data_str)
	err := json.Unmarshal(data, this)
	return err
}

//生成ECC椭圆曲线密钥对
func GenerateECCKeyString() (string, string, error) {

	//生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
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

//对消息的散列值生成数字签名
func Sign(msg string, private_str string) (string, error) {
	//取得私钥
	privateKey, err := PrivateKeyFromString(private_str)
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

	tmp_str, err := sign.ToJson()
	if err != nil {
		return "", err
	}

	return tmp_str, nil
}

//验证数字签名
func VerifySign(msg string, sign string, publish_key_str string) bool {
	//读取公钥
	publicKey, err := PublicKeyFromString(publish_key_str)
	if err != nil {
		return false
	}

	ecc_sign := EccSign{}
	err = ecc_sign.FromJson(sign)
	if err != nil {
		return false
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
	return verify
}
