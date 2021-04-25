package beego_filter

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

//过滤域名
func FilterCheckDomain(ctx *context.Context) {

	origin := ctx.Request.Header.Get("Origin")
	ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)

	filterFunc := cors.Allow(&cors.Options{
		//AllowAllOrigins: this.allowAllDomain,
		//AllowAllOrigins: true,
		AllowOrigins: []string{origin},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Accept", "Content-Type", "x-token"},
		//AllowHeaders:     []string{"*", "x-token", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	})

	filterFunc(ctx)
}
