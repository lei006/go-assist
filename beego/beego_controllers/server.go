package beego_controllers

import (
	"livertc/core/beego_assert"
	"livertc/core/intfs"
	"livertc/core/servers/server_monitor"

	beego "github.com/beego/beego/v2/server/web"
)

type ServerController struct {
	beego_assert.BaseController
	Monitor server_monitor.MonitorServer
}

func (this *ServerController) Router(app intfs.IApplication, prefix string) {
	this.App = app

	beego.Router(prefix+"servers", this, "get:GetAll")

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

func (this *ServerController) GetAll() {

	obs := this.App.ServerItems("")

	this.SuccReturn(obs)
}
