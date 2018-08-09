package main

import (
	_ "qixijie/routers"
	"github.com/astaxie/beego"
	"qixijie/def"
	"qixijie/util"
)
func init(){
	util.InitCache()
}
func main() {
	beego.SetStaticPath("/MP_verify_oSClQLOUTyzPRg6o.txt","MP_verify_oSClQLOUTyzPRg6o.txt")
	def.Outtradeno=beego.AppConfig.String("outtradeno")
	beego.Run()
}

