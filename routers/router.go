package routers

import (
	"qixijie/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//正式路由器
	//登陆
	beego.Router("/redirecturl", &controllers.MainController{}, "*:Redirecturl")
	beego.Router("/index", &controllers.MainController{}, "get:Index")
	//分享接口
	beego.Router("/upimageAndmessage", &controllers.MainController{}, "*:UpImageAndMessage")

	beego.Router("/share/get_ticker", &controllers.MainController{}, "get:GetTicker")
	beego.Router("/share/get_user_token", &controllers.MainController{}, "post:GetToken")


}
