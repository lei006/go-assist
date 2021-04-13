package beego_controllers

import (
	"fmt"
	"livertc/core/beego_assert"
	"livertc/core/beego_filter"
	"livertc/core/beego_models"
	"livertc/core/intfs"
	"livertc/core/types"
	"livertc/core/utils/errno"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"
)

var (

	// 系统错误, 前缀为 100
	ERR_AUTH_ERROR = &errno.Errno{Code: 10101, Message: "认证出错!"}
)

// Operations about Users
type AuthController struct {
	beego_assert.BaseController
}

func (this *AuthController) Router(app intfs.IApplication, prefix string) {
	this.App = app
	beego.Router(prefix+"auth/login", this, "post:Login")   //发布者
	beego.Router(prefix+"auth/logout", this, "post:Logout") //发布者
	beego.Router(prefix+"auth/info", this, "get:TokenInfo") //发布者

	//不检查token-白名单
	beego_filter.Add_NoCheckToken(prefix+"auth/login", "post") //登录接口

	///////////////////////////////////////////////////////////
	// 过滤中间件....
	beego_filter.App = app
	beego.InsertFilter(prefix+"*", beego.BeforeRouter, beego_filter.FilterCheckDomain) //检查域名
	beego.InsertFilter(prefix+"*", beego.BeforeRouter, beego_filter.FilterCheckToken)  //检查token 是否存在...
	beego.InsertFilter(prefix+"*", beego.BeforeRouter, beego_filter.FilterPremissions) //检查权限
}

/**
 * @api {get} /auth/login?password=:password 1.登录
 * @apiGroup 2.登录
 * @apiName Login
 * @apiParam {String} password 密码
 * @apiSuccessExample {json} 成功响应:
	{
		code: 20000
		data: {id: 4, token: "alw73dd0j8d5ae6en381hf0rhm9r", expires: 600, created: 1599556386, updated: 1599556386}
		message: "success"
	}
*/

func (ctl *AuthController) Login() {

	json_data := string(ctl.Ctx.Input.RequestBody)

	username := gjson.Get(json_data, "username").String()
	password := gjson.Get(json_data, "password").String()

	is_exist, err := beego_models.ModUser.IsExist(username)
	if err != nil {
		logs.Warn("服务器内部错误:" + err.Error())
		ctl.ReturnFail("服务器内部错误:")
	}

	//如果不存在，且是 admin
	if is_exist == false && username == "admin" {

		val := beego.AppConfig.DefaultString(types.ConfigName_AdminPassword, types.Default_AdminPassword)
		if password == val {
			tmp := &beego_models.User{
				Username:  "admin",
				Avatar:    "admin",
				SectionId: 99999999,
			}
			new_token, err := beego_models.ModToken.AddOne(tmp)
			if err != nil {
				logs.Error("Add One token error:" + err.Error())
				ctl.ReturnFail("服务器内部错误")
			}
			tokenInfo, err := beego_models.ModToken.GetOne(new_token)
			if err != nil {
				logs.Error("add new token GetOne (" + new_token + ")  err:" + err.Error())
				ctl.ReturnFail("服务器内部错误")
			}
			fmt.Println("login success", tokenInfo)
			ctl.SuccReturn(tokenInfo)
			return
		}
	}

	user, err := beego_models.ModUser.GetOneByUsername(username)
	if err == orm.ErrNoRows {
		ctl.ReturnFail("用户不存在:" + username)
	}
	if err != nil {
		ctl.ReturnFail("服务器内部错误:" + username)
	}
	if user.IsEnable == false {
		ctl.ReturnFail("此用户已经禁用:" + username)
	}

	if user.Password != password {
		ctl.ReturnFail("认证失败,用户名或密码错误")
	}

	new_token, err := beego_models.ModToken.AddOne(user)
	if err != nil {
		logs.Error("Add One token error:" + err.Error())
		ctl.ReturnFail("服务器内部错误")
	}

	tokenInfo, err := beego_models.ModToken.GetOne(new_token)
	if err != nil {
		logs.Error("add new token GetOne (" + new_token + ")  err:" + err.Error())
		ctl.ReturnFail("服务器内部错误")
	}

	ctl.SuccReturn(tokenInfo)
}

/**
 * @api {get} /auth/logout 2.登出
 * @apiGroup 2.登录
 * @apiName Logout
 * @apiHeader {String} x-token 授权码
 * @apiSuccessExample {json} 成功响应:
	{
	"message": "success",
	"code": 20000,
	"data": {}
	}
*/
func (this *AuthController) Logout() {

	token := this.Ctx.Request.Header.Get("x-token")
	err1 := beego_models.ModToken.DeleteOne(token)
	if err1 != nil {
		this.ReturnFail("delete token error:" + err1.Error())
	}

	this.SuccReturn("ok")
}

/**
 * @api {get} /auth/info 3.信息
 * @apiGroup 2.登录
 * @apiName info
 * @apiHeader {String} x-token 授权码
 * @apiSuccessExample {json} 成功响应:
	{
	"message": "success",
	"code": 20000,
	"data": {
		created: 1599554171
		expires: 600
		id: 3
		token: "ov2wwsz7dfcitcyac2uhd7ssrega"
		updated: 1599554435
	}
	}
*/

func (this *AuthController) TokenInfo() {

	token := this.Ctx.Request.Header.Get("x-token")

	info, err := beego_models.ModToken.GetOne(token)
	if err != nil {
		this.ReturnFail("find token error:" + err.Error())
	}
	info.Avatar = "['admin']"

	this.SuccReturn(info)
}
