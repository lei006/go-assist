package beego_controllers

import (
	"livertc/core/beego_assert"
	"livertc/core/beego_filter"
	"livertc/core/beego_models"
	"livertc/core/intfs"
	"livertc/core/servers/server_monitor"
	"livertc/core/types"
	"os"

	beego "github.com/beego/beego/v2/server/web"

	"github.com/beego/beego/v2/adapter/logs"
)

// Operations about object
type SystemController struct {
	beego_assert.BaseController
}

func MakeSystemController() *SystemController {
	tmp := &SystemController{}
	return tmp
}
func (this *SystemController) Router(app intfs.IApplication, prefix string) {
	this.App = app
	beego.Router(prefix+"system/about", this, "get:GetAbout")   //软件信息
	beego.Router(prefix+"system/restart", this, "post:Restart") //重启服务--这里只是关闭它，守护进程会自动启动它
	beego.Router(prefix+"system/license", this, "get:GetLicense")
	beego.Router(prefix+"system/license", this, "put:SetLicense")

	beego_filter.Add_CheckAdmin(prefix+"system/restart", "post") //重启系统-需要管理员权限
	beego_filter.Add_NoCheckToken(prefix+"system/about", "get")  //取得软件信息

}

/**
 * @api {get} /About 1.版本数据
 * @apiGroup 5.系统信息
 * @apiName Version
 * @apiHeader {String} x-token 授权码
 * @apiSuccessExample {json} 成功响应(只保留三份数据,实际更多.):
{
  "message": "success",
  "code": 20000,
  "data": [
    {
      "key": "硬件信息",
      "val": "xxx3 公司"
    },
    {
      "key": "运行时间",
      "val": "xxx2 公司"
    },
    {
      "key": "软件信息",
      "val": "xxx1 公司"
    }
  ]
}

*/

func (this *SystemController) GetAbout() {

	this.SuccReturn(this.App.GetAbout())
}

/**
 * @api {get} /system/stat 2.统计数据
 * @apiGroup 5.系统信息
 * @apiName Stat
 * @apiHeader {String} x-token 授权码
 * @apiSuccessExample {json} 成功响应(只保留三份数据,实际更多.):
{
  "message": "success",
  "code": 20000,
  "data": {
    "mem": [
      {
        "time": "2020-06-18 09:23:24",
        "使用": 54
      },
      {
        "time": "2020-06-18 09:23:26",
        "使用": 54
      },
      {
        "time": "2020-06-18 09:23:28",
        "使用": 54
      },
    ],
    "cpu": [
      {
        "time": "2020-06-18 09:23:24",
        "使用": 53
      },
      {
        "time": "2020-06-18 09:23:26",
        "使用": 27
      },
      {
        "time": "2020-06-18 09:23:28",
        "使用": 19
      },
    ]
  }

*/

func (this *SystemController) GetMonitorData() {

	obs := server_monitor.GetMonitorData()

	this.SuccReturn(obs)
}

func (this *SystemController) GetLicense() {

	data := this.App.GetLicenser().GetInfo()

	this.ReturnSuccess(data)
}

func (this *SystemController) SetLicense() {

	lic_data := string(this.Ctx.Input.RequestBody)
	err := this.App.GetLicenser().SetLicense(lic_data)
	if err != nil {
		logs.Warn("设置lic出错:" + err.Error())
		this.ReturnFail("设置lic出错")
	}
	this.ReturnSuccess("ok")
}

func (this *SystemController) Restart() {

	logs.Warn("系统重启")
	os.Exit(-1)
	this.SuccReturn("ok")
}

func (this *SystemController) CheckAdmin() {
	tmp := this.Ctx.Input.GetData(types.TokenInfo_KEY)
	if tmp == nil {
		this.ReturnFail("无操作权限")
	}

	token := tmp.(*beego_models.Token)
	if token.IsAdmin() == false {
		this.ReturnFail("无操作权限:管理员操作")
	}
}
