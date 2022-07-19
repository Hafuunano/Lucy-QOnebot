package funwork // 每日龙王

import (
	"time"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine.OnNotice().SetBlock(false).Handle(func(ctx *zero.Ctx) {
		if ctx.Event.NoticeType == "current_talkative" {
			list := ctx.GetGroupHonorInfo(ctx.Event.GroupID, "talkative")
			temp := list.String()
			id := gjson.Get(temp, "current_talkative.user_id")
			name := ctx.CardOrNickName(id.Int())
			ctx.SendChain(message.Text("今日的龙王是~", name, "哦"))
			if id.Int() == ctx.Event.SelfID {
				time.Sleep(time.Second * 20)
				ctx.Send(message.Text("好欸~今天咱是龙王哦www"))
			} else {
				time.Sleep(time.Second * 20)
				ctx.Send(message.Text("今天咱没有拿到龙王qaq"))
			}
		}
	})
	engine.OnFullMatch("今日龙王", zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		list := ctx.GetGroupHonorInfo(ctx.Event.GroupID, "talkative")
		temp := list.String()
		id := gjson.Get(temp, "current_talkative.user_id")
		name := ctx.CardOrNickName(id.Int())
		ctx.SendChain(message.Text("今日的龙王是~", name, "哦"))
	})
}
