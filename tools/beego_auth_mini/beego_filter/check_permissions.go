package beego_filter

import (
	"livertc/core/beego_assert"
	"livertc/core/beego_models"
	"livertc/core/types"
	"net/url"
	"strings"

	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/server/web/context"
)

//检查管理员权限
var g_checkAdmin []FilterItem //  whitelist

//增加管理员检查表
func Add_CheckAdmin(api_path, method string) {
	method = strings.ToUpper(method)
	g_checkAdmin = append(g_checkAdmin, FilterItem{Method: method, Url: api_path})
}

//过滤检查token
func FilterPremissions(ctx *context.Context) {

	//取得当前的 url
	cur_url := ""
	tmpUrl, err := url.Parse(ctx.Request.RequestURI)
	if err == nil {
		cur_url = strings.ToLower(tmpUrl.Path)
	}

	//取得当前的 方法
	cur_method := strings.ToUpper(ctx.Request.Method)

	for _, item := range g_checkAdmin {
		pos := strings.Index(cur_url, item.Url)

		if (pos == 0) && (cur_method == item.Method) {

			// url 方法都符合，则检查管理员权限
			tmp := ctx.Input.GetData(types.TokenInfo_KEY)
			if tmp == nil {
				//没找到数据，则放过....不检查权限---应该不会走到这一步...
				logs.Warn("检查管理员权限时没有找到 token 信息,这是不应该的...")
				beego_assert.ReturnFail(ctx, "权限不足")
				return
			}

			//需要检查管理员
			token := tmp.(*beego_models.Token)
			if token.IsAdmin() == false {
				beego_assert.ReturnFail(ctx, "权限不足")
				return
			}
		}
	}
}
