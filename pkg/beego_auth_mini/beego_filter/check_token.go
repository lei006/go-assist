package beego_filter

import (
	"fmt"
	"livertc/core/beego_assert"
	"livertc/core/beego_models"
	"livertc/core/intfs"
	"livertc/core/types"
	"net/url"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web/context"

	"github.com/beego/beego/v2/adapter/logs"
)

type FilterItem struct {
	Method string
	Url    string
}

var g_nocheckList []FilterItem //  whitelist

var (
	NoCheckToken    = false                           //不检查 token 默认: 检查
	TokenExpireTime = time.Duration(10) * time.Minute //有效时间 10分钟
	TokenDeleteTime = time.Duration(60) * time.Minute //删除时间 60分钟
)

var App intfs.IApplication

//过滤检查token
func FilterCheckToken(ctx *context.Context) {

	//////////////////////////////////
	// 检查白名单
	//取得当前的 url
	cur_url := ""
	tmpUrl, err := url.Parse(ctx.Request.RequestURI)
	if err == nil {
		cur_url = strings.ToLower(tmpUrl.Path)
	}

	//取得当前的 方法
	cur_method := strings.ToUpper(ctx.Request.Method)

	fmt.Println("FilterCheckToken -- > ", ctx.Input.URL())

	for _, item := range g_nocheckList {
		pos := strings.Index(cur_url, item.Url)
		if (pos == 0) && (cur_method == item.Method) {
			//不检查白名单的内容
			return
		}
	}

	//////////////////////////////////
	// 记录操作的 token
	token := ctx.Request.Header.Get("x-token")
	tokenInfo, err := beego_models.ModToken.GetOne(token)
	if err != nil {
		logs.Error("get token:" + err.Error())
		beego_assert.ReturnFail(ctx, "token error:"+err.Error())
		return
	}

	ctx.Input.SetData(types.TokenInfo_KEY, tokenInfo)

	////////////////////////////////////
	// 如果是演示版...

	if App.GetLicenser().IsDemo() {
		//演示时间为 10分钟
		if (time.Now().Unix() - tokenInfo.CreatedAt.Unix()) > 10*60 {
			beego_assert.ReturnFailCode(ctx, 50014, "演示 token 已过期")
		}
	}

	//////////////////////////////////
	// 检查 Token
	bret, err := beego_models.ModToken.IsExpired(token)
	if err != nil {
		logs.Warn("get token err:" + err.Error())
		beego_assert.ReturnFail(ctx, "token err:"+err.Error())
		return
	}

	if bret == true {
		beego_assert.ReturnFailCode(ctx, 50014, "token 已过期")
		return
	}
	//更新时间...
	err = beego_models.ModToken.UpdateExpiredAt(token, TokenExpireTime)
	if err != nil {
		logs.Warn("UpdateExpiredAt err:" + err.Error())
	}

	//删除过期token
	beego_models.ModToken.DeleteOverTime(TokenDeleteTime)

}

func Add_NoCheckToken(api_path, method string) {
	method = strings.ToUpper(method)
	g_nocheckList = append(g_nocheckList, FilterItem{Method: method, Url: api_path})
}
