package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"qixijie/util"
	"qixijie/def"
	"qixijie/db"
	"qixijie/model"
	"time"
	"net/http"
	"strings"
	"crypto/sha1"
	"io/ioutil"
	"math/rand"
)

var mongp = util.P{}

func init() {
	mongp["host"] = beego.AppConfig.String("mongodburl")
	mongp["name"] = beego.AppConfig.String("mongodbdb")
	mongp["username"] = beego.AppConfig.String("mongodbuser")
	mongp["password"] = beego.AppConfig.String("mongodbpass")
	mongp["userinfo"] = beego.AppConfig.String("mongodbdocm")
	mongp["userdata"] = beego.AppConfig.String("mongodbdata")
	mongp["userprize"] = beego.AppConfig.String("mongodbprize")

}

type MainController struct {
	beego.Controller
}

//把表白信息和图片保存到数据库
func (c *MainController) UpImageAndMessage() {
	defer func() {
		if err := recover(); err != nil {
			c.Redirect("/seven_night/redirecturl", 302)
		}
	}()
	data := model.Node{}
	userdata := mongp["userdata"]
	docm2string := util.ToString(userdata)
	mongdb := db.D(docm2string, mongp)
	medId := c.GetString("medId")
	data.From = c.GetString("userform")
	data.To = c.GetString("touser")
	data.Word = c.GetString("word")
	ifOpenShow := c.GetString("ifOpenShow")
	fmt.Print(medId)
	data.Timestamp = util.ToString(time.Now().Unix())
	imagePath := ""
	if medId == "" {
		imagePath = "http://service.genyuanlian.com/seven_night/static/images/my.png"
	} else {
		imagePath = util.GetImageFromCould(medId, "./image/")
		if imagePath == "fail" {
			imagePath = "http://service.genyuanlian.com/seven_night/static/images/my.png"
		}
	}
	imagePath1 := imagePath
	data.ImageUrl = imagePath1
	addr, payId := util.GetEthAddress()
	data.Addr = addr
	openid := c.GetString("openid")
	mongdb.Query = &util.P{"userOpenId": openid}
	count := mongdb.Count()
	if count < 1 {
		//数据上链
		ethaddr := UpMessage(data.From + data.Word + data.To)
		//获取数据高度
		height := getQurHig()
		if height == "" {
			height = "0"
		}
		if ethaddr != "false" {
			mongdb.Add(util.P{"data": data, "userOpenId": openid, "IfPay": "1", "addr": data.Addr, "ethaddr": ethaddr, "height": height, "ifOpenShow": ifOpenShow})
			//返回中奖信息
			Prize:=InserPrize(openid,ethaddr)
			c.Data["json"] = map[string]interface{}{"userOpenId": openid, "addr": data.Addr, "payId": payId, "isPay": 1, "code": 0, "ethaddr": ethaddr, "height": height,"prize":Prize}
			c.ServeJSON()
		}
	} else {
		ifshare, time := getShareInfo(openid)
		if ifshare == "Yes" && time == "0" {
			//上链
			ethaddr := UpMessage(data.From + data.Word + data.To)
			//获取区块高度
			height := getQurHig()
			if height == "" {
				height = "0"
			}
			if ethaddr != "false" {
				mongdb.Add(util.P{"data": data, "userOpenId": openid, "IfPay": "1", "addr": data.Addr, "ethaddr": ethaddr, "height": height, "ifOpenShow": ifOpenShow})
				Prize:=InserPrize(openid,ethaddr)
				//更新用户分享信息
				SetShareInfo(openid)
				c.Data["json"] = map[string]interface{}{"userOpenId": openid, "addr": data.Addr, "payId": payId, "isPay": 1, "code": 0, "ethaddr": ethaddr, "height": height,"prize":Prize}
				c.ServeJSON()
			}
		} else {
			mongdb.Add(util.P{"data": data, "userOpenId": openid, "IfPay": "0", "addr": data.Addr, "ifOpenShow": ifOpenShow})
			c.Data["json"] = map[string]interface{}{"userOpenId": openid, "addr": data.Addr, "payId": payId, "isPay": 0, "code": 0}
			c.ServeJSON()
		}
	}
}

//获取所有表白信息
func (c *MainController) GetUserMessage() {
	c.Data["json"] = getMessage("", "", "")
	c.ServeJSON()
}

//跳转微信的url
func (c *MainController) Redirecturl() {
	Addr := c.GetString("addr")
	ShareOpenid := c.GetString("shareopenid")
	util.S("addr", Addr)
	util.S("shareopenid", ShareOpenid)

	redirecturl := "https%3a%2f%2fservice.genyuanlian.com%2fseven_night%2findex"
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + def.WEIXINAPPID + "&redirect_uri=" + redirecturl +
		"&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect"
	c.Redirect(url, 302)

}

//由微信服务器跳转回来的rul
func (c *MainController) Index() {
	defer func() {
		if err := recover(); err != nil {
			c.Redirect("/seven_night/redirecturl", 302)
		}
	}()
	Addr := c.GetString("addr")
	ShareOpenid := c.GetString("shareopenid")
	//没有openid 就进主页
	if util.IsEmpty(ShareOpenid) {
		//获取滚动信息,获取上墙信息，只有isOpenShow 为true时可以
		AllUp := getMessage("", "", "true")
		mp := getShowOpenMessage(*AllUp)

		c.Data["tanmu"] = mp
		user := &util.P{}
		docm := mongp["userinfo"]
		docm2string := util.ToString(docm)
		mongdb := db.D(docm2string, mongp)

		code := c.GetString("code")
		if code == "" || code == "undefined" {
			c.Redirect("/seven_night/redirecturl", 302)
		}
		userinfo := util.GetUserInfo(code)
		if userinfo == nil {
			c.Redirect("/seven_night/redirecturl", 302)
		}
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
		nodelist := []model.History{}
		openid := (*userinfo)["openid"].(string)
		p := getMessage(openid, "", "")
		for _, v := range *p {
			node := model.History{}
			datap := v["data"]
			hight := util.ToString(v["height"])
			data := datap.(util.P)
			node.From = util.ToString(data["from"])
			node.ImageUrl = util.ToString(data["imageurl"])
			node.To = util.ToString(data["to"])
			node.Word = util.ToString(data["word"])
			node.Addr = util.ToString(data["addr"])
			node.Timestamp = util.ToString(data["timestamp"])
			node.EthAddr = util.ToString(v["ethaddr"])
			fmt.Print("----------------------------------------------------------------------------------------------------------------ethaddr-------------------------------")
			fmt.Print(util.ToString(v["ethaddr"]))
			fmt.Print("----------------------------------------------------------------------------------------------------------------ethaddr-------------------------------")

			node.Height = hight
			nodelist = append(nodelist, node)
		}
		c.Data["openid"] = (*userinfo)["openid"].(string)
		c.Data["nodelist"] = nodelist
		countAll:=getMessageCount
		c.Data["count"]=countAll
		c.TplName = "home.html"
	} else {
		p := getMessage(ShareOpenid, Addr, "")
		data := (*p)[0]["data"]
		url:="https://block.genyuanlian.com/tx/"+util.ToString((*p)[0]["ethaddr"])
		datap := data.(util.P)
		from := datap["from"]
		to := datap["to"]
		word := datap["word"]
		imageurl := datap["imageurl"]
		heigh := (*p)[0]["height"]
		c.Data["imageurl"] = imageurl
		c.Data["from"] = from
		c.Data["to"] = to
		c.Data["hrefurl"]=url
		c.Data["word"] = word
		c.Data["height"] = heigh
		c.Data["oder"] = "share"
		c.Data["isReachTime"] = timeoder()
		c.TplName = "home.html"
	}
}

// 获取历史数据
func (c *MainController) GetHistoryMessage() {
	nodelist := []model.History{}
	openid := c.GetString("openid")
	p := getMessage(openid, "", "")
	for _, v := range *p {
		node := model.History{}
		datap := v["data"]
		height := util.ToString(v["height"])
		ethaddr := v["ethaddr"].(string)
		data := datap.(util.P)
		node.From = util.ToString(data["from"])
		node.To = util.ToString(data["to"])
		node.ImageUrl = util.ToString(data["imageurl"])
		node.Word = util.ToString(data["word"])
		node.Timestamp = util.ToString(data["timestamp"])
		node.Addr = util.ToString(data["addr"])
		node.EthAddr = ethaddr
		node.Height = height
		nodelist = append(nodelist, node)
	}
	c.Data["json"] = nodelist
	c.ServeJSON()
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
func getMessage(openid, addr, isOpenShow string) (p *[]util.P) {

	userdata := mongp["userdata"]
	docm2string := util.ToString(userdata)
	mongdb := db.D(docm2string, mongp)
	if len(addr) == 0 && len(openid) > 0 {
		p = mongdb.Find(util.P{"userOpenId": openid, "IfPay": "1"}).All()
		return p
	} else if len(addr) > 0 && len(openid) > 0 {
		p = mongdb.Find(util.P{"userOpenId": openid, "addr": addr}).All()
		return p
	} else {
		p = mongdb.Find(util.P{"ifOpenShow": isOpenShow}).All()
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

//设置已经用过分享的机会了
func SetShareInfo(useropenid string) (ifshare, usersharetime string) {
	docm := mongp["userinfo"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)
	err := mongdb.Upsert(util.P{"userOpenId": useropenid}, util.P{"usersharetime": "1"})
	if err == nil {
	}

	return
}

//获取微信支付的pay_id
func (c *MainController) GetWxPayId() {
	xml := util.GetWXpay_id("osy0OwaMZDR7dxcMcfUPh750GBqA")
	c.Data["json"] = xml
	c.ServeJSON()
}

//分享成功，分享标签为Yes
func (c *MainController) Sharesuccess() {
	openid := c.GetString("openid")

	docm := mongp["userinfo"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)
	err := mongdb.Upsert(util.P{"userOpenId": openid}, util.P{"ifshare": "Yes"})
	if err == nil {
		c.Ctx.WriteString("success")
	} else {

		c.Ctx.WriteString("false")
	}
}

//上链
func UpMessage(data string) string {

	APP_KEY := "5b6bfaf6a6dd527199fce0c1"
	APP_SECRECT := "f3f0fad23d0803d01619ad06b9cd1469c9ac61a37d794d76cafd5247c272fe38259a07784dfde3cad9e6a6a55340ed17";
	sign := fmt.Sprintf("%x", sha1.Sum([]byte(APP_KEY+data+APP_SECRECT)))
	mp := util.P{}
	mp["data"] = data
	mp["app_key"] = APP_KEY
	mp["sign"] = sign
	response, _ := http.Post("http://chromeapi.genyuanlian.com:9005/api/upload/org", "application/json;charset=utf-8", strings.NewReader(util.JsonEncode(mp)))
	defer response.Body.Close()
	eth_body, _ := ioutil.ReadAll(response.Body)
	p := *util.JsonDecode([]byte(string(eth_body)))
	flag := p["code"].(float64)

	if flag == 0 {
		data := p["data"].([]interface{})
		s := data[0].(string)
		return s
	} else {
		return "false"
	}
}

//查询支付状态，如果支付则上链和更新数据支付信息
func (c *MainController) CheckPay() {
	payid := c.GetString("payid")
	openid := c.GetString("openid")
	addr := c.GetString("addr")
	fla := util.CheckIfPay(payid)
	if fla {

		from, to, word, ethaddr, height := setUpMessage(openid, addr)
		pay := model.Pay{}
		pay.Paid = true
		pay.Upload_error = false
		pay.From = from
		pay.To = to
		pay.Word = word
		pay.Height = height
		pay.EthAddr = ethaddr
		if len(ethaddr) > 0 {
			pay.Success = true
			//支付成功返回抽奖字段
		 pay.Prize=InserPrize(openid,ethaddr)
		}
		c.Data["json"] = map[string]interface{}{"code": 0, "data": pay}
		c.ServeJSON()

	} else {
		c.Data["json"] = map[string]interface{}{"code": 1, "ethaddr": ""}
		c.ServeJSON()
	}
}

//抽奖的接口
func (c *MainController) Prize() {
	openid := c.GetString("openid")
	ethAddr := c.GetString("ethAddr")
	phoneNumber := c.GetString("number")
	useraddr := c.GetString("useraddr")
	time := time.Now().Unix()
	docm := mongp["userprize"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)
	count:=mongdb.Find(util.P{"userOpenId": openid,"ethAddr": ethAddr}).Count()
	if count<1{
		c.Data["json"] = util.P{"code": 1}
		c.ServeJSON()
	}else {
	err:=mongdb.Upsert(util.P{"userOpenId": openid,"ethAddr": ethAddr}, util.P{"phoneNumber": phoneNumber,"userAddr":useraddr,"timestamp":time})
	if err == nil {
		c.Data["json"] = util.P{"code": 0}
		c.ServeJSON()
	} else {
		c.Data["json"] = util.P{"code": 1}
		c.ServeJSON()
	}
	}
}
//抽奖的接口，内部返回接口
func  InserPrize(openid string,ethAddr string)(prizeIn int) {
	docm := mongp["userprize"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)
	//获取prize的中的奖项
	prize,prizeInt:=getPrizeRand()
	err := mongdb.Add(util.P{"userOpenId": openid, "prize": prize,"ethAddr":ethAddr})
	if err == nil {
	return prizeInt
	} else {
		return 11
	}
}

func getPrizeRand() (stringa string,inti int){
	randstring:=[]string{"9.9BSTK","9.9BSTK","9.9BSTK","9.9BSTK","9.9BSTK","99BSTK","蜂蜜","小红薯","算力包","石魁粉"}
	randint:=[]int{1,1,1,1,1,1,2,3,4,5}
	i:=rand.Intn(10)
	stringa=randstring[i]
	inti=randint[i]
	return stringa,inti
}





//上链和更新数据库的支付信息
func setUpMessage(openid, addr string) (from, to, word, ethaddr, height string) {
	docm := mongp["userdata"]
	docm2string := util.ToString(docm)
	mongdb := db.D(docm2string, mongp)
	p := mongdb.Find(util.P{"userOpenId": openid, "addr": addr}).One()
	node := (*p)["data"]
	nodep := node.(util.P)
	from = nodep["from"].(string)
	to = nodep["to"].(string)
	word = nodep["word"].(string)
	data := from + to + word
	//上链
	ethaddr = UpMessage(data)
	//获取区块高度
	height = getQurHig()
	if height == "" {
		height = "0"
	}
	//更新数据库，支付字段，返回的区块链地址值，区块地址
	err := mongdb.Upsert(util.P{"userOpenId": openid, "addr": addr}, util.P{"ethaddr": ethaddr, "IfPay": "1", "height": height})
	if err == nil {

		return from, to, word, ethaddr, height
	} else {
		mongdb.Upsert(util.P{"userOpenId": openid, "addr": addr}, util.P{"ethaddr": ethaddr, "IfPay": "1", "height": height})
		return from, to, data, ethaddr, height
	}

}

//获取区块高度
func getQurHig() string {
	defer func() {
		if err := recover(); err != nil {

			fmt.Print("获取高度失败")

		}

	}()
	retrnbody, _ := http.Get("http://chromeapi.genyuanlian.com:3001/api/status?q=getInfo")
	defer retrnbody.Body.Close()
	eth_body, _ := ioutil.ReadAll(retrnbody.Body)
	p := *util.JsonDecode([]byte(string(eth_body)))
	info := p["info"].(interface{})
	data := info.(map[string]interface{})
	blocks := util.ToString(data["blocks"])
	return util.ToString(blocks)
}

//获取九个数据
func getShowOpenMessage(AllUp []util.P) (mpp []util.P) {

	if len(AllUp) < 9 {
		lenth := 9 - len(AllUp)
		for i := 1; i < lenth; i++ {
			node := util.P{}
			node["word"] = "dsafewfew"
			p := util.P{"data": node}
			AllUp = append(AllUp, p)

		}
		return AllUp
	} else {
		mp := []util.P{}
		for i := 0; i < 9; i++ {
			mp = append(mp, AllUp[rand.Intn(len(AllUp))])

		}
		return mp
	}
}

//获取时间对比
func timeoder() int {
	nowtime := time.Now()
	nowtimenuix := nowtime.Unix()
	the_time, err := time.ParseInLocation("2006-01-02 15:04:05", "2018-08-17 00:00:00", time.Local)
	if err == nil {
		unix_time := the_time.Unix()
		if nowtimenuix < unix_time {
			return 0
		} else {
			return 1
		}
	} else {
		return 1
	}

}
func getMessageCount()(count int){
	userdata := mongp["userdata"]
	docm2string := util.ToString(userdata)
	mongdb := db.D(docm2string, mongp)
	 count=mongdb.Count()
		return
	 }