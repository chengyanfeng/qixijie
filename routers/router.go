package routers

import (
	"qixijie/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//正式路由器

	//登陆
	beego.Router("/seven_night/redirecturl", &controllers.MainController{}, "*:Redirecturl")
	beego.Router("/seven_night/index", &controllers.MainController{}, "get:Index")
	//分享接口
	beego.Router("/seven_night/upimageAndmessage", &controllers.MainController{}, "*:UpImageAndMessage")
	beego.Router("/seven_night/getpayid", &controllers.MainController{}, "*:GetWxPayId")
	beego.Router("/seven_night/gethistorymessage", &controllers.MainController{}, "post:GetHistoryMessage")
	//分享成功后通知后端添加一次机会
	beego.Router("/seven_night/sharesucess", &controllers.MainController{}, "get:Sharesuccess")
	//更新奖项记录
	beego.Router("/seven_night/prize", &controllers.MainController{}, "post:Prize")

	//查看支付状态
	beego.Router("/seven_night/checkpay", &controllers.MainController{}, "*:CheckPay")
	beego.Router("/seven_night/share/get_ticker", &controllers.MainController{}, "get:GetTicker")
	beego.Router("/seven_night/share/get_user_token", &controllers.MainController{}, "post:GetToken")




}
