package beego_controllers

import (
	"encoding/json"
	"livertc/core/beego_assert"
	"livertc/core/beego_filter"
	"livertc/core/beego_models"
	"livertc/core/intfs"
	"livertc/core/types"
	"regexp"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type UserController struct {
	beego_assert.BaseController
	NoCheckToken bool //	不检查token
	Prefix       string
}

func (this *UserController) Router(app intfs.IApplication, prefix string) {
	this.App = app
	beego.Router(prefix+"users", this, "get:GetAll")  //发布者
	beego.Router(prefix+"user", this, "post:Add")     //发布者
	beego.Router(prefix+"user", this, "get:UserInfo") //发布者
	beego.Router(prefix+"user/:id", this)             //发布者

	beego_filter.Add_CheckAdmin(prefix+"user", "post")   //增加用户-需要管理员权限
	beego_filter.Add_CheckAdmin(prefix+"user", "delete") //删除用户-需要管理员权限

}

func (this *UserController) Add() {

	var user beego_models.User
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)

	match, _ := regexp.MatchString("^[A-Za-z0-9]{1,20}$", user.Username)
	if match == false {
		this.ReturnFail("用户名错误: 必需长度为3-20的字母与数字:" + user.Username)
	}

	if len(user.Password) < 3 {
		this.FailReturn(50010, "密码太短了!")
	}

	_, err := beego_models.ModUser.GetOneByUsername(user.Username)
	if err == nil {
		this.ReturnFail("用户已经存在")
	}
	if err != orm.ErrNoRows {
		this.ReturnFail("服务器内部错误:" + err.Error())
	}

	user_id, err := beego_models.ModUser.AddOne(user)
	if err != nil {
		this.ReturnFail("服务器内部错误:" + err.Error())
	}

	userinfo, err := beego_models.ModUser.GetOne(user_id)
	if err != nil {
		this.ReturnFail("服务器内部错误:" + err.Error())
	} else {
		logs.Info("新加用户:", userinfo)
		this.SuccReturn(userinfo)
	}
}

func (this *UserController) Put() {

	token_info := this.GetTokenData()
	if token_info == nil {
		this.ReturnFail("权限不足: 没有登录")
	}

	user_id, err := this.GetInt64(":id", -1)
	if err != nil {
		this.ReturnFail("id出错")
	}

	var user beego_models.User
	json.Unmarshal(this.Ctx.Input.RequestBody, &user)

	old_user_info, err := beego_models.ModUser.GetOne(user_id)
	if err != nil {
		this.ReturnFail("出错:" + err.Error())
	}

	if token_info.IsAdmin() == false {

		//非管理员...
		if token_info.Username != user.Username {
			this.ReturnFail("权限不足: 不能修改其它用户的信息")
		}

		if old_user_info.Avatar != user.Avatar {
			this.ReturnFail("权限不足: 不能修改角色信息")
		}

		if old_user_info.SectionId != user.SectionId {
			this.ReturnFail("权限不足: 不能修改所属部门")
		}
		if old_user_info.IsEnable != user.IsEnable {
			this.ReturnFail("权限不足: 不能修改所处状态")
		}
	}

	err = beego_models.ModUser.UpdateOne(old_user_info.Id, &user)
	if err != nil {
		this.ReturnFail("更新失败:" + err.Error())
	}
	this.SuccReturn("ok")

}

func (this *UserController) Delete() {

	user_id, err := this.GetInt64(":id", -1)
	if err != nil {
		this.ReturnFail("id出错")
	}

	err = beego_models.ModUser.DeleteOne(user_id)
	if err != nil {
		this.ReturnFail("服务器内部错误:" + err.Error())
	}

	logs.Notice("删除用户: id=", user_id)
	this.SuccReturn("ok")

}

func (this *UserController) Get() {

	token_info := this.GetTokenData()
	if token_info == nil {
		this.ReturnFail("未登录，无法获得用户表")
	}

	user_name := this.GetString(":id")

	if token_info.Username != user_name {
		//如果是查别人的..则需要管理员权限
		this.CheckAdminPermissions()
	}

	tmp, err := beego_models.ModUser.GetOneByUsername(user_name)
	if err != nil {
		this.ReturnFail("取得用户信息出错:" + err.Error())
	}
	tmp.Password = ""
	this.ReturnSuccess(tmp)
}

//取得当前用户信息
func (this *UserController) UserInfo() {

	token_info := this.GetTokenData()
	if token_info == nil {
		this.ReturnFail("未登录，无法获得用户信息")
	}

	tmp, err := beego_models.ModUser.GetOneByUsername(token_info.Username)
	if err != nil {
		this.ReturnFail("取得用户信息出错:" + err.Error())
	}
	tmp.Password = ""
	this.ReturnSuccess(tmp)
}

func (this *UserController) GetAll() {

	token_info := this.GetTokenData()
	if token_info == nil {
		this.ReturnFail("未登录，无法获得用户表")
	}

	if token_info.IsAdmin() == true {
		list, num, err := beego_models.ModUser.GetAll()
		if err != nil {
			this.ReturnFail("错误:" + err.Error())
		}
		this.SuccReturnList(list, int64(num))
	}

	//非管理员，只能看到自身
	var users = []*beego_models.User{}
	userinfo, err := beego_models.ModUser.GetOneByUsername(token_info.Username)
	if err != nil {
		this.ReturnFail("错误:" + err.Error())
		return
	}
	users = append(users, userinfo)
	this.SuccReturnList(users, int64(len(users)))
}

//检查必需为管理员....
func (this *UserController) CheckAdminPermissions() *beego_models.Token {

	tokenData := this.GetTokenData()
	if tokenData != nil && tokenData.Avatar != "admin" {
		this.ReturnFail("权限不足:管理员操作(" + tokenData.Avatar + ")")
	}

	return tokenData
}

//取得token数据 操作..
func (this *UserController) GetTokenData() *beego_models.Token {

	tmp := this.Ctx.Input.GetData(types.TokenInfo_KEY)
	if tmp == nil {
		this.ReturnFail("权限不足: 取得token出错")
		return nil
	}

	tokenData := tmp.(*beego_models.Token)

	return tokenData
}
