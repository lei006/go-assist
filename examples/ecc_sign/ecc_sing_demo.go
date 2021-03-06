package main

import (
	"fmt"

	"github.com/lei006/go-assist/crypto/ecc_tool"
)

func test_sign(index int) {

	//生成ECC密钥对文件
	pub_key, pri_key, err := ecc_tool.GenerateECCKeyString()
	if err != nil {
		fmt.Println("生成密钥对出错")
		return
	}

	//fmt.Println("pub_key", pub_key)
	//fmt.Println("pri_key", pri_key)

	//模拟发送者
	//要发送的消息
	msg := "hello world"
	//生成数字签名
	sign, err := ecc_tool.Sign(msg, pri_key)
	if err != nil {
		fmt.Println("对数据签名出错")
		return
	}
	//模拟接受者
	//接受到的消息
	acceptmsg := "hello world"

	//验证签名
	verifySignECC := ecc_tool.VerifySign(acceptmsg, sign, pub_key)

	fmt.Println("验证结果：", index, verifySignECC, sign)

}

func main() {

	for i := 0; i < 1000; i++ {
		test_sign(i)
	}

}
