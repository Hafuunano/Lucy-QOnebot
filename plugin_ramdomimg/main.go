// Package ramimg 随机本地图片
package ramdomimg

import (

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/rule"

	"github.com/FloatTech/ZeroBot-Plugin/order"
)

const (
	datapath = "data/ramimg"
	dbfile   = datapath + "/data.db"
	cfgfile  = datapath + "/ramimgpath.txt"
)

var (
	ramimgpath = "/tmp" // 绝对路径，图片根目录
)

func init() {
	engine := control.Register("ramdomimg", order.Prioramimg, &control.Options{
		DisableOnDefault: false,
		Help: "\n" +
			"- 来份[xxx]\n",
	})
	engine.OnRegex(`^来份(.*)$`, rule.FirstValueInList(ns)).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			imgtype := ctx.State["regex_matched"].([]string)[1]
			sc := new(ramimgclass)
			ns.mu.RLock()
			err := ns.db.Pick(imgtype, sc)
			ns.mu.RUnlock()
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
			} else {
				p := "file:///" + ramimgpath + "/" + sc.Path
				ctx.SendChain(message.Image(p))
			}
		})
		engine.OnFullMatch("刷新随机图片", zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			err := ns.scanall(ramimgpath)
			if err == nil {
				ctx.SendChain(message.Text("成功！"))
			} else {
				ctx.SendChain(message.Text("ERROR: ", err))
			}
		})
}
