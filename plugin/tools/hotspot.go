package tools // 参考了 Cha0sIDL 的 zbp hotspot热点

import (
	"strconv"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
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

	engine.OnFullMatch("今日早报").Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.Send(message.Image("https://api.03c3.cn/zb/"))
	})
}
