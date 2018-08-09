package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"qixijie/util"
	"qixijie/def"
	"qixijie/db"
)

var mongp = util.P{}

func init() {
	mongp["host"] = beego.AppConfig.String("mongodburl")
	mongp["name"] = beego.AppConfig.String("mongodbdb")
	mongp["username"] = beego.AppConfig.String("mongodbuser")
	mongp["password"] = beego.AppConfig.String("mongodbpass")
	mongp["docm"] = beego.AppConfig.String("mongodbdocm")
}

type MainController struct {
	beego.Controller
}

//跳转微信的url
func (c *MainController) UpImageAndMessage() {
	medId := c.GetString("medId")
	imagePath := util.GetImageFromCould(medId, "./image/")
	/*err:=db.D(docm2string, mongp).Add(user)*/
	fmt.Print(imagePath)
}

//跳转微信的url
func (c *MainController) Redirecturl() {
	redirecturl := "http%3a%2f%2fchengyanfeng.natapp1.cc%2findex"
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + def.WEIXINAPPID + "&redirect_uri=" + redirecturl +
		"&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect"
	c.Redirect(url, 301)
}

//由微信服务器跳转回来的rul
func (c *MainController) Index() {

	user := &util.P{}
	docm := mongp["docm"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)

	code := c.GetString("code")
	if code == "" {
		c.Ctx.WriteString("请用微信登陆")
	}
	userinfo := util.GetUserInfo(code)
	(*user)["country"] = (*userinfo)["country"].(string)
	(*user)["city"] = (*userinfo)["city"].(string)
	(*user)["userOpenId"] = (*userinfo)["openid"].(string)
	(*user)["name"] = (*userinfo)["nickname"].(string)
	(*user)["province"] = (*userinfo)["province"].(string)
	mongdb.Query = &util.P{"userOpenId": (*userinfo)["openid"].(string)}
	count := mongdb.Count()
	if count < 0 {
		mongdb.Add(user)
	}
	c.Data["openid"] = (*userinfo)["openid"].(string)
	c.TplName = "index.html"
}

//微信获取的转发token
func (c *MainController) GetToken() {
	Forwardtoken := util.GetForwardToken()
	c.Ctx.WriteString(util.ToString(Forwardtoken))
}

//微信获取分享的ticker
func (c *MainController) GetTicker() {
	token := util.S("forword_token")
	tick := util.GetTicket(util.ToString(token))
	fmt.Print(tick, "---------------tick----------")
	c.Ctx.WriteString(string(tick))
}
