package licenser_normal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/lei006/go-assist/tools/ecc_tool"
)

// licenser 更新回调
type LicenserUpdateCallback func()

type Licenser struct {
	AppName    string
	AppCode    string
	HardSn     string
	PublishKey string //公钥
	is_demo    bool   //演示版

	Data      *LicenserData
	callbacks []LicenserUpdateCallback //更新回调列表
}

func MakeLicenser(appname, appcode, hardsn string, publish_key string) *Licenser {

	licenserData := &LicenserData{
		AppName:     appname,
		AppCode:     appcode,
		CompanyName: "这是测试公司",
		HardSn:      hardsn,
		MaxNum:      1,
		Desc:        "这是测试描述",
	}

	licenser := &Licenser{
		AppName:    appname,
		AppCode:    appcode,
		HardSn:     hardsn,
		PublishKey: publish_key,
		Data:       licenserData,
		is_demo:    true,
	}

	return licenser
}

func (this *Licenser) LoadFromFile(filename string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = this.LoadData(string(data))
	return err
}

func (this *Licenser) SetCallback(cb LicenserUpdateCallback) {
	this.callbacks = append(this.callbacks, cb)
}

func (this *Licenser) LoadData(lic_data string) error {

	license_data := &LicenserData{}

	// 1. 解码licsense数据
	err := license_data.FromJson(lic_data)
	if err != nil {
		return errors.New("格式出错")
	}

	// 2. 检查数据...
	err = this.checkData(license_data)
	if err != nil {
		return err
	}

	// 3. 设置数据
	this.is_demo = false
	this.Data = license_data

	// 4. 通知生效
	for _, val := range this.callbacks {
		val()
	}

	return nil
}

//检查数据...不回调..
func (this *Licenser) checkData(license_data *LicenserData) error {

	// 2. 验证 AppName
	if this.AppName != license_data.AppName {
		return errors.New("AppName不同")
	}

	// 3. 验证 AppCode
	if this.AppCode != license_data.AppCode {
		return errors.New("AppCode不同")
	}

	// 4. 验证 hardSn
	if this.HardSn != license_data.HardSn {
		return errors.New("硬件码不同")
	}

	// 5. 验证签名
	text := license_data.ToString()
	ret := this.EccVerifySign(text, license_data.Sign, this.PublishKey)
	if ret == false {
		return errors.New("验签名失败")
	}

	return nil
}

//制做签名
func (this *Licenser) EccSign(data string, pri_key string) (string, error) {
	return ecc_tool.Sign(data, pri_key)
}

//效验签名
func (this *Licenser) EccVerifySign(data string, sign string, pub_key string) bool {
	return ecc_tool.VerifySign(data, sign, pub_key)
}

//是否演示版
func (this *Licenser) IsDemo() bool {

	//过期也认为是演示
	if this.IsExpired() {
		return true
	}

	return this.is_demo
}

func (this *Licenser) GetMaxNum() int64 {

	if this.IsDemo() {
		//如果是演示
		return 2000
	}

	return this.Data.MaxNum
}

//是否过期
func (this *Licenser) IsExpired() bool {
	return time.Now().Unix() > this.Data.ExpireAt
}

type DataItem struct {
	Key string      `bson:"key" json:"key"`
	Val interface{} `bson:"val" json:"val"`
}

func (this *Licenser) GetInfo() interface{} {

	list := []DataItem{}
	list = append(list, DataItem{Key: "公司名", Val: this.Data.CompanyName})
	list = append(list, DataItem{Key: "授权程序", Val: this.Data.AppCode})
	if this.IsDemo() {
		list = append(list, DataItem{Key: "授权描述", Val: "演示版"})
	} else {
		list = append(list, DataItem{Key: "授权描述", Val: this.Data.Desc})
	}

	list = append(list, DataItem{Key: "授权硬件", Val: this.Data.HardSn})
	list = append(list, DataItem{Key: "最大数量", Val: fmt.Sprintf("%d", this.Data.MaxNum)})

	expire_val := UnixToTimeStr(this.Data.ExpireAt)
	list = append(list, DataItem{Key: "过期时间", Val: expire_val})

	return list
}

func UnixToTimeStr(t_data int64) string {
	timeLayout := "2006-01-02 15:04:05" //转化所需模板
	//时间戳转日期
	dataTimeStr := time.Unix(t_data, 0).Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	return dataTimeStr

}
