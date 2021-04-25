package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lei006/go-assist/tools/ecc_tool"
	"github.com/lei006/go-assist/tools/licenser_normal"
	"github.com/lei006/go-assist/utils"
)

func test1() {

	//licensefile := "./license.lic"

	//生成ECC密钥对文件
	pub_key, pri_key, err := ecc_tool.GenerateECCKeyString()
	if err != nil {
		fmt.Println("生成密钥对出错")
		return
	}

	for i := 0; i < 10000; i++ {

		appName := utils.RandomString1(10)
		appCode := utils.RandomString1(10)
		hardsn := utils.RandomString1(32)
		CompanyName := utils.RandomString1(32)
		Copyright := utils.RandomString1(32)
		maxNum := rand.Int63()

		licenser := licenser_normal.MakeLicenser(appName, appCode, hardsn, pub_key)
		lic_data := licenser_normal.MakeLicenserData(appCode, hardsn)

		lic_data.CompanyName = CompanyName
		lic_data.MaxNum = maxNum
		lic_data.Copyright = Copyright

		lic_data_str := lic_data.ToString()
		//fmt.Println("lic_data_str --> ", lic_data_str)

		sign, err := licenser.EccSign(lic_data_str, pri_key)
		if err != nil {
			fmt.Println("error ", i, err.Error())
		} else {
			bret := licenser.EccVerifySign(lic_data_str, sign, pub_key)
			fmt.Println("bret =", i, bret)
		}
	}

	/*
		err := licenser.LoadFromFile(licensefile)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("ok")
		}
	*/

	time.Sleep(time.Duration(10) * time.Second)
}
