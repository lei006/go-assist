package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/lei006/go-assist/license"
)

func test_sign() {
	key, err := license.MakEccP521Key()
	if err != nil {
		fmt.Println("生成KEY错误:", err)
		return
	}
	acceptmsg := license.RandString(100000)

	verify, err := key.TestSign(acceptmsg)
	fmt.Println("测试签名验证：", verify, err)
}

func test_key() {

	verify___to_string := false
	verify___pub_key := false

	key, err := license.MakEccP521Key()
	if err != nil {
		fmt.Println("生成KEY错误:", err)
		return
	}

	{

		key_str, err := key.ToBase64String()
		if err != nil {
			fmt.Println("key变为字符串出错:", err)
			return
		}

		new_key_0 := &license.LicenseKey{}
		err = new_key_0.FromBase64String(key_str)
		if err != nil {
			fmt.Println("字符串变为key出错:", err)
			return
		}

		if new_key_0.ToCompare(key) == false {

			json_0, _ := key.ToJson()
			json_1, _ := new_key_0.ToJson()
			fmt.Println("字符串变为key, 变前变后不一样" + json_0 + json_1)
			return
		}

		verify___to_string, err = new_key_0.TestSign(license.RandString(1000))
		if err != nil {
			fmt.Println("验证失败:", err)
			return
		}

	}

	{

		rand_data := license.RandString(1000)

		sign, err := key.Sign(rand_data)
		if err != nil {
			fmt.Println("对数据签名出错")
			return
		}

		pub_key__bate64 := key.GetBase64PubKey()

		new_key_1 := &license.LicenseKey{}
		err = new_key_1.SetBase64PubKey(pub_key__bate64)
		if err != nil {
			fmt.Println("设置公钥失败:", err)
			return
		}

		verify___pub_key, err = new_key_1.VerifySign(rand_data, sign)
		if err != nil {
			fmt.Println("验证失败:", err)
			return
		}

	}
	fmt.Println("验证结果:", verify___to_string, verify___pub_key)

}

func test_license00() {

	suc, err := license.TestLicense()
	fmt.Println("test license ", suc, err)

}

func test_license01() (bool, error) {

	key, err := license.MakEccP521Key()
	if err != nil {
		return false, errors.New("生成KEY错误:" + err.Error())
	}

	expireTime := time.Now().Add(3 * time.Second) //过期时间取值

	enc_data := &license.LicenseClaims{}
	enc_data.StandardClaims.Id = license.RandString(10)
	enc_data.StandardClaims.Subject = license.RandString(100)
	enc_data.StandardClaims.Issuer = "乐园天"
	enc_data.StandardClaims.ExpiresAt = expireTime.Unix()
	enc_data.Number = time.Now().Unix()

	lic_data, err := license.KeyEncryption(enc_data, key)
	if err != nil {
		return false, errors.New("license Encryption error:" + err.Error())
	}
	dec_data, err := license.Decryption(lic_data)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return false, errors.New("license Decryption error:" + err.Error())
	}

	//比较解码前后是否一样...
	if dec_data.ToCompare(enc_data) == false {
		fmt.Println("解码前后数据不一致")
		return false, errors.New("解码前后不一致")
	}

	fmt.Println("lic_data==>", "解码前后数据一致")

	return true, nil

}

func main() {
	for i := 0; i < 20; i++ {
		test_sign()
	}
	for i := 0; i < 20; i++ {
		test_key()
	}
	for i := 0; i < 20; i++ {
		test_license00()
	}
	for i := 0; i < 20; i++ {
		test_license01()
	}

	//test_license(100)

}
