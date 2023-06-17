// Package slash https://github.com/Rongronggg9/SlashBot
package slash

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

var (
	// so noisy so use this.
	engine = control.Register("slash", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "Hi NekoPachi!\n",
	})
)

func init() {
	engine.OnRegex(`^/(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPatternInfo := ctx.State["regex_matched"].([]string)[1]
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ctx.CardOrNickName(ctx.Event.UserID)+getPatternInfo+"了他自己~"))
	})
	engine.OnRegex(`^(\[CQ:reply,id=(\d+)\])\s/(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPatternUserMessageID := ctx.State["regex_matched"].([]string)[2]
		getPatternInfo := ctx.State["regex_matched"].([]string)[3]
		getSplit := strings.Split(getPatternInfo, " ")
		rsp := ctx.CallAction("get_msg", zero.Params{
			"message_id": getPatternUserMessageID,
		}).Data.String()
		sender := gjson.Get(rsp, "sender.user_id").Int()
		if len(getSplit) == 2 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ctx.CardOrNickName(ctx.Event.UserID)+" "+getSplit[0]+" 了 "+ctx.CardOrNickName(sender)+" "+getSplit[1]))
		} else {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(ctx.CardOrNickName(ctx.Event.UserID)+" "+getPatternInfo+" 了 "+ctx.CardOrNickName(sender)))
		}
	})
}
