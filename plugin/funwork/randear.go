package funwork

import (
	"math/rand"

	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine.OnFullMatchGroup([]string{"今天用什么耳机", "抽耳机"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("建议"),
				randText("AKG K52", "AKG K72", "AKG K92", "AKG K121", "AKG K141", "AKG K167", "AKG K181DJ", "AKG K182", "AKG K240", "AKG K245", "AKG K271", "AKG K275", "AKG K361", "AKG K371", "AKG K400", "AKG K401", "AKG K403", "AKG K420", "AKG K430", "AKG K440NC", "AKG K450", "AKG K495NC", "AKG K500", "AKG K501", "AKG K520", "AKG K530", "AKG K540", "AKG K541", "AKG K545", "AKG K550", "AKG K551", "AKG K553", "AKG K601", "AKG K602", "AKG K612PRO", "AKG K618DJ", "AKG K619DJ", "AKG K701", "AKG K702", "AKG K712PRO", "AKG K812", "AKG K872", "AKG K1000", "AKG Q200", "AKG Q460", "AKG Q701", "AKG Y30", "AKG Y40", "AKG Y45BT", "AKG Y50", "AKG Y60NC", "AKG Y400", "AKG Y500", "AKG Y600NC", "AKG N60NC", "AKG N700NC"))
		})
	engine.OnFullMatch("猪").Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		switch {
		case ctx.Event.UserID == 363128:
			message.At(2896285821)
		default:
		}
	})
	engine.OnFullMatch("呆呆瓜").Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		switch {
		case ctx.Event.UserID == 2896285821:
			message.At(363128)
		default:
		}
	})

}

func randText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}
