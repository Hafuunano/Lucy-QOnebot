package action // action for Lucy

import (
	"math/rand"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	limit   = rate.NewManager[int64](time.Minute*10, 15)
	Lucyimg = "file:///root/Lucy_Project/memes/" // lucy的meme表情包地址
)

func init() {
	engine := control.Register("action", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "Lucy容易被动触发语言 \n- 默认禁用 可以自行打开哦",
	})
	// Okk

	engine.OnFullMatchGroup([]string{"？", "?", "¿"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			process.SleepAbout1sTo2s()
			switch rand.Intn(5) {
			case 0:
				ctx.SendChain(randText("?", "？", "嗯？", "(。´・ω・)ん?", "ん？"))
			case 1, 2:
				ctx.SendChain(randImage("cats.jpg", "322E8EBA2B08815460119BE93342E33B.png", "111.jpg"))
			}
		})

	engine.OnFullMatch("草").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			switch rand.Intn(2) {
			case 0:

				ctx.SendChain(randText("（一种植物）", "ん？"))
			case 1, 2:
				ctx.SendChain(randImage("5cee2a0f5dc32a114b1a9d3f60314e5e.jpg", "R-C.jpeg", "hikari.jpg"))
			}
		})
	engine.OnKeyword("离谱").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			switch rand.Intn(3) {
			case 0:
				ctx.SendChain(randText("?", "？", "ん？"))
			case 1, 2:
				ctx.SendChain(randImage("-33ee3a0711f11810.jpg", "111.jpg"))
			}
		})
	engine.OnFullMatchGroup([]string{"呜呜", "呜呜呜"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			switch rand.Intn(2) {
			case 0:
				ctx.SendChain(randImage("2925511468257.png", "FBFBBBBA433464163949F55085266356.png"), message.Text(
					"摸摸~"))
			case 1:
				ctx.SendChain(message.Text(
					"抱抱~"))
			}
		})

	engine.OnFullMatchGroup([]string{"喵", "喵喵", "喵喵喵"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			switch rand.Intn(2) {
			case 0:
				ctx.SendChain(randText("喵喵~", "喵w~"))
			case 1:
				ctx.SendChain(randImage("6152277811454.jpg", "meow.jpg", "ww.jpg"))
			}
		})
	engine.OnFullMatch("咕咕").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			ctx.SendChain(randText("抓到一只鸽子OwO", "是鸽子 炖了~", "咕咕咕", "不许咕咕咕"))
		})
		// 114514
	engine.OnRegex(`^我要(.*)份涩图`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			ctx.SendChain(message.Image("https://gchat.qpic.cn/gchatpic_new/1770747317/1049468946-3068097579-76A49478EFA68B4750B10B96917F7B58/0?term=3"))
		})
	// end
	engine.OnFullMatch("都怪Zheic", zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			ctx.SendChain(message.Image("https://gchat.qpic.cn/gchatpic_new/1292581422/812841489-2901631038-0DE3B3F02343B09C591FDFFEE586A353/0?term=3"))
		})
}

func randText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(Lucyimg + file[rand.Intn(len(file))])
}
