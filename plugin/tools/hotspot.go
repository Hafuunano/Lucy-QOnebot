package tools // 参考了 Cha0sIDL 的 zbp hotspot热点

import (
	"strconv"
	"strings"

	"github.com/FloatTech/zbputils/binary"
	"github.com/FloatTech/zbputils/web"
	"github.com/antchfx/htmlquery"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine.OnFullMatch("微博热搜", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		rsp := "微博实时热榜:\n"
		url := "http://api.weibo.cn/2/guest/search/hot/word"
		data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA())
		if err != nil {
			msg := message.Text("ERROR:", err)
			ctx.SendChain(msg)
			return
		}
		json := gjson.Get(binary.BytesToString(data), "data").Array()
		for idx, hot := range json {
			if hot.Get("word").String() == "" {
				continue
			}
			rsp = rsp + strconv.Itoa(idx+1) + ":" + hot.Get("word").String() + "\n"
		}
		ctx.SendChain(message.Text(rsp))
	})

	engine.OnFullMatch("github热搜", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		msg := "GitHub实时热榜:\n"
		doc, err := htmlquery.LoadURL("https://github.com/trending")
		if err != nil {
			panic("htmlQuery error")
		}
		article := htmlquery.Find(doc, "//*[@id=\"js-pjax-container\"]/div[3]/div/div[2]/article[@*]")
		for idx, a := range article {
			titlePath := htmlquery.FindOne(a, "/h1/a")
			title := htmlquery.SelectAttr(titlePath, "href")
			msg += strconv.Itoa(idx+1) + "：" + strings.TrimPrefix(title, "/") + "\n" + "地址：https://github.com" + title + "\n"
		}
		ctx.SendChain(message.Text(msg))
	})
}
