package main

import (
	"fmt"

	"github.com/lei006/go-assist/license"
)

func test_sign(num int) {
	key, err := license.MakEccP521Key()
	if err != nil {
		fmt.Println("生成KEY错误:", err)
		return
	}
	acceptmsg := "hello world"

	for i := 0; i < num; i++ {
		verify, err := key.TestSign(acceptmsg)
		fmt.Println("验证结果：", i, verify, err)
	}
}

func test_license(num int) {

	key, err := license.MakEccP521Key()
	if err != nil {
		fmt.Println("生成KEY错误:", err)
		return
	}

	enc_data := license.LicenseClaims{}
	enc_data.StandardClaims.Id = "aaa"

	lic_data, err := license.Encryption(enc_data, key)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	dec_data, err := license.Decryption(lic_data, key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	if enc_data.StandardClaims.Id == dec_data.StandardClaims.Id {
		fmt.Println("ok")
	} else {
		fmt.Println("err")
	}

}

func main() {
	//test_sign(100)

	test_license(100)
	suc, err := license.TestLicense()
	fmt.Println("TestLicense ret:", suc, err)

	/*
		claims_str, err := dec_data.ToJson()
		if err != nil {
			fmt.Println("claims to json error:", err)
			return
		}
	*/
}
