package licenser_normal

import (
	"encoding/json"
	"fmt"
)

type LicenserData struct {
	AppName     string `bson:"appname" json:"appname"`         //程序名
	AppCode     string `bson:"appcode" json:"appcode"`         //程序名
	CompanyName string `bson:"companyname" json:"companyname"` //授权公司
	HardSn      string `bson:"hardsn" json:"hardsn"`           //授权硬件id
	MaxNum      int64  `bson:"maxnum" json:"maxnum"`           //最大数量
	ExpireAt    int64  `bson:"expire_at" json:"expire_at"`     //过期时间签
	Copyright   string `bson:"copyright" json:"copyright"`     //版权所有人
	Desc        string `bson:"desc" json:"desc"`               //描述
	Sign        string `bson:"sign" json:"sign"`               //签名字符串
	PubKey      string `bson:"pub_key" json:"pub_key"`         //公钥
}

func MakeLicenserData(appcode, hardsn string) *LicenserData {

	licenserData := &LicenserData{
		AppCode: appcode,
		HardSn:  hardsn,
	}
	return licenserData
}

//生成验签的字符串
func (this *LicenserData) ToString() string {

	text := fmt.Sprintf("%s%s%s%s-%d%d-%s%s%s",
		this.AppName, this.AppCode, this.CompanyName, this.HardSn,
		this.MaxNum, this.ExpireAt,
		this.Copyright, this.Desc, this.PubKey)

	return text
}

func (this *LicenserData) ToJson() (string, error) {
	b, err := json.Marshal(this)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (this *LicenserData) FromJson(data_str string) error {
	data := []byte(data_str)
	err := json.Unmarshal(data, this)
	return err
}

func (this *LicenserData) OutString() string {
	out_str := ""
	out_str += "CompanyName:" + this.CompanyName + "\r\n"
	out_str += "AppName    :" + this.AppName + "\r\n"
	out_str += "AppCode    :" + this.AppCode + "\r\n"
	out_str += "MaxNum     :" + fmt.Sprintf("%d", this.MaxNum) + "\r\n"
	out_str += "ExpireAt   :" + UnixToTimeStr(this.ExpireAt) + "\r\n"
	return out_str
}
