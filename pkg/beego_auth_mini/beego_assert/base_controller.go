package beego_assert

import (
	"encoding/json"
	"livertc/core/intfs"
	"time"

	"github.com/beego/beego/v2/adapter/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

type BaseController struct {
	beego.Controller
	App intfs.IApplication
}

type JsonReturn struct {
	Msg  string      `json:"message"`
	Code int         `json:"code"`
	Data interface{} `json:"data"` //Data字段需要设置为interface类型以便接收任意数据
	Now  int64       `json:"now"`
}

func (this *JsonReturn) ToJson() (string, error) {

	b, err := json.Marshal(this)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *BaseController) ApiJsonReturn(code int, msg string, data interface{}) {

	origin := c.Ctx.Request.Header.Get("Origin")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)

	var _jsonReturn JsonReturn
	_jsonReturn.Msg = msg
	_jsonReturn.Code = code
	_jsonReturn.Data = data
	_jsonReturn.Now = time.Now().Unix()
	c.Data["json"] = _jsonReturn //将结构体数组根据tag解析为json
	c.ServeJSON()                //对json进行序列化输出
	c.StopRun()                  //终止执行逻辑

}

func CtxReturn(ctx *context.Context, code int, msg string, data interface{}) {

	origin := ctx.Request.Header.Get("Origin")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)

	var _jsonReturn JsonReturn
	_jsonReturn.Msg = msg
	_jsonReturn.Code = code
	_jsonReturn.Data = data
	_jsonReturn.Now = time.Now().Unix()

	msg, err := _jsonReturn.ToJson()
	if err != nil {
		logs.Error(err.Error())
	}
	ctx.WriteString(msg)
}

type JsonDataList struct {
	Items interface{} `json:"items"` //Data字段需要设置为interface类型以便接收任意数据
	Total int64       `json:"total"`
}

type JsonReturnList struct {
	Msg  string       `json:"message"`
	Code int          `json:"code"`
	Data JsonDataList `json:"data"` //Data字段需要设置为interface类型以便接收任意数据
	Now  int64        `json:"now"`
}

func (c *BaseController) SuccReturnList(data interface{}, total int64) {
	origin := c.Ctx.Request.Header.Get("Origin")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)

	var _ret JsonReturnList
	_ret.Msg = "success"
	_ret.Code = 20000
	_ret.Data.Items = data
	_ret.Data.Total = total
	_ret.Now = time.Now().Unix()

	c.Data["json"] = _ret //将结构体数组根据tag解析为json
	c.ServeJSON()         //对json进行序列化输出
	c.StopRun()           //终止执行逻辑
}

func (c *BaseController) SuccReturn(data interface{}) {
	c.ApiJsonReturn(20000, "success", data)
}

func (c *BaseController) FailReturn(code int, msg string) {
	c.ApiJsonReturn(code, msg, nil)
}

func (c *BaseController) ReturnSuccess(data interface{}) {
	c.ApiJsonReturn(20000, "success", data)
}

func (c *BaseController) ReturnList(data interface{}, total int64) {
	c.SuccReturnList(data, total)
}

func (c *BaseController) ReturnFail(msg string) {
	logs.Warn(msg)
	c.ApiJsonReturn(50010, msg, nil)
}

func (c *BaseController) CtxReturnFail(ctx *context.Context, msg string) {
	logs.Warn(msg)
	CtxReturn(ctx, 50010, msg, nil)
}

func (c *BaseController) ReturnFailCode(code int, msg string) {
	c.ApiJsonReturn(code, msg, nil)
}

func ReturnFail(ctx *context.Context, msg string) {
	CtxReturn(ctx, 50010, msg, nil)
}

func ReturnFailCode(ctx *context.Context, code int, msg string) {
	CtxReturn(ctx, code, msg, nil)
}
