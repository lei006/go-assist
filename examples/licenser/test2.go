package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lei006/go-assist/tools/ecc_tool"
	"github.com/lei006/go-assist/tools/licenser_normal"
	"github.com/lei006/go-assist/utils"
)

func test2() {

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

		lic_data.AppName = appName
		lic_data.CompanyName = CompanyName
		lic_data.MaxNum = maxNum
		lic_data.Copyright = Copyright

		lic_data_str := lic_data.ToString()
		//fmt.Println("lic_data_str --> ", lic_data_str)

		sign, err := licenser.EccSign(lic_data_str, pri_key)
		if err != nil {
			fmt.Println("error ", i, err.Error())
		} else {

			lic_data.Sign = sign
			//lic_data.PubKey = pub_key

			lic_data_str111, err_json := lic_data.ToJson()
			if err_json != nil {
				fmt.Println("err_json ", err_json.Error())
				continue
			}

			bret := licenser.SetData(lic_data_str111, func(data *licenser_normal.LicenserData) {

			})
			if bret != nil {
				fmt.Println("-----------text:", lic_data_str)
				fmt.Println("-----------pub_key:", pub_key)
				fmt.Println("-----------sign:", sign)
			}

			fmt.Println("bret =", i, bret)
			fmt.Println("------------------------------------------------")
			fmt.Println("------------------------------------------------")
			fmt.Println("------------------------------------------------")
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
