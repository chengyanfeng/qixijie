package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"qixijie/util"
	"qixijie/def"
	"qixijie/db"
	"qixijie/model"
	"time"
)

var mongp = util.P{}

func init() {
	mongp["host"] = beego.AppConfig.String("mongodburl")
	mongp["name"] = beego.AppConfig.String("mongodbdb")
	mongp["username"] = beego.AppConfig.String("mongodbuser")
	mongp["password"] = beego.AppConfig.String("mongodbpass")
	mongp["userinfo"] = beego.AppConfig.String("mongodbdocm")
	mongp["userdata"] = beego.AppConfig.String("mongodbdata")
}

type MainController struct {
	beego.Controller
}

//把表白信息和图片保存到数据库
func (c *MainController) UpImageAndMessage() {
	data := model.Node{}
	userdata := mongp["userdata"]
	docm2string := util.ToString(userdata)
	mongdb := db.D(docm2string, mongp)
	medId := c.GetString("medId")
	data.From = c.GetString("userform")
	data.To = c.GetString("touser")
	data.Word = c.GetString("word")
	fmt.Print(medId)
	data.Timestamp = util.ToString(time.Now().Unix())
	/*imagePath := util.GetImageFromCould(medId, "./image/")*/
	imagePath := "aaaa"
	data.ImageUrl = imagePath
	ethaddr, payId:= util.GetEthAddress()
	data.Addr = ethaddr
	openid := c.GetString("openid")
	mongdb.Query = &util.P{"userOpenId": openid}
	count := mongdb.Count()
	if count <10{
		mongdb.Add(util.P{"data": data, "userOpenId": openid, "IfPay": "1", "Addr": data.Addr,})
		c.Data["json"]=map[string]interface{}{"userOpenId":openid,"addr":data.Addr,"payId":payId,"isPay":1,"code":0}
		c.ServeJSON()
	} else {
		ifshare, time := getShareInfo(openid)
		if ifshare == "Yes" && time == "0" {
			mongdb.Add(util.P{"data": data, "userOpenId": openid, "IfPay": "1", "addr": data.Addr})
			c.Data["userOpenId"]=openid
			c.Data["Addr"]=data.Addr
			c.Data["payId"]=payId
			c.Data["isPay"]=1
			c.TplName="share.html"
		} else {
			mongdb.Add(util.P{"data": data, "userOpenId": openid, "IfPay": "0", "addr": data.Addr})
			c.Data["userOpenId"]=openid
			c.Data["Addr"]=data.Addr
			c.Data["payId"]=payId
			c.Data["isPay"]=0
			c.TplName="share.html"
			}
	}
}

//获取所有表白信息
func (c *MainController) GetUserMessage() {
	getMessage("", "")
	c.Data["json"] = getMessage("", "")
	c.ServeJSON()
}

//跳转微信的url
func (c *MainController) Redirecturl() {
	Addr := c.GetString("addr")
	ShareOpenid := c.GetString("shareopenid")
	util.S("addr", Addr)
	util.S("shareopenid", ShareOpenid)
	redirecturl := "http%3a%2f%2fchengyanfeng.natapp1.cc%2findex"
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + def.WEIXINAPPID + "&redirect_uri=" + redirecturl +
		"&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect"
	c.Redirect(url, 302)


}

//由微信服务器跳转回来的rul
func (c *MainController) Index() {
	Addr := c.GetString("addr")
	ShareOpenid := c.GetString("shareopenid")
	if util.IsEmpty(ShareOpenid) {
		user := &util.P{}
		docm := mongp["userinfo"]
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
		if count == 0 {
			(*user)["ifshare"] = "No"
			(*user)["usersharetime"] = "0"
			mongdb.Add(user)
		}
		c.Data["openid"] = (*userinfo)["openid"].(string)
		c.Data["oder"]="first"
		c.TplName = "home.html"
	} else {

		p := getMessage(ShareOpenid, Addr)
		data:=(*p)[0]["data"]
		datap:=data.(util.P)
		from:=datap["from"]
		to:=datap["to"]
		word:=datap["word"]
		c.Data["from"]=from
		c.Data["to"]=to
		c.Data["word"]=word
		c.Data["oder"]="share"
		c.TplName = "home.html"
	}
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

//获取表白信息
func getMessage(openid, addr string) (p *[]util.P) {
	userdata := mongp["userdata"]
	docm2string := util.ToString(userdata)
	mongdb := db.D(docm2string, mongp)
	if  len(addr) == 0 {
		p = mongdb.Find(util.P{"userOpenId": openid}).All()
		return p
	} else {
		p = mongdb.Find(util.P{"userOpenId": openid, "Addr": addr}).All()
		return p
	}

}

//获取分享信息和次数
func getShareInfo(useropenid string) (ifshare, usersharetime string) {
	docm := mongp["userinfo"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)
	p := mongdb.Find(util.P{"userOpenId": useropenid}).One()
	ifshare = (*p)["ifshare"].(string)
	usersharetime = (*p)["usersharetime"].(string)
	return
}
