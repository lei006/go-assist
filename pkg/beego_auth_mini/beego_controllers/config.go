package beego_controllers

import (
	"livertc/core/beego_assert"
	"livertc/core/beego_filter"
	"livertc/core/intfs"
	"livertc/core/tools"

	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

type ConfigController struct {
	beego_assert.BaseController
}

func (this *ConfigController) Router(app intfs.IApplication, prefix string) {
	this.App = app

	beego.Router(prefix+"configs", this, "get:GetAll")
	beego.Router(prefix+"config/:id", this, "get:GetConfig")
	beego.Router(prefix+"config/:id", this, "put:SetConfig")

	//需要管理员权限
	beego_filter.Add_CheckAdmin(prefix+"config", "put") //配置系统-需要管理员权限

}

func (this *ConfigController) GetAll() {

	this.SuccReturn("ok")
}

func (this *ConfigController) GetConfig() {
	key := this.GetString(":id")

	val, err := tools.GetConfigItem(key)
	if err != nil {
		//不需要重数据库中取，因为程序里，都没有，说明没有用...
		logs.Warn("ControllerGetConfig ", key, " error:", err.Error())
		this.ReturnFail(err.Error())
	}

	this.SuccReturn(val)
}

func (this *ConfigController) SetConfig() {

	key := this.GetString(":id")
	body := string(this.Ctx.Input.RequestBody)
	data := gjson.Get(body, "data").String()

	// 1. 必需先设置，
	err := tools.SetConfigItem(key, data)
	if err != nil {
		logs.Warn("ControllerSetConfig ", key, "=", data, " error:", err.Error())
		this.FailReturn(50010, err.Error())
	}

	//2. 再保存到系统
	_, err = tools.ReplaceOneConfig(key, data)
	if err != nil {
		logs.Warn("ReplaceOneConfig ", key, " error:", err.Error())
		this.FailReturn(50010, err.Error())
	}

	this.SuccReturn("ok")
}
