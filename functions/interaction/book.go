package interaction

import (
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"math/rand"
	"os"
	"strconv"
)

func init() {
	data, err := os.ReadFile(engine.DataFolder() + "answers.json")
	if err != nil {
		panic(err)
	}
	engine.OnFullMatch("答案之书").SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦, 可以和咱说下是什么问题呢"))
		nextstep := ctx.FutureEvent("message", ctx.CheckSession())
		nextstep.Repeat() // get repeat but here no reply yet.
		answerListInt := rand.Intn(268)
		answerListStr := strconv.Itoa(answerListInt)
		answer := gjson.Get(helper.BytesToString(data), answerListStr+".answer")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(answer))
	})
}
