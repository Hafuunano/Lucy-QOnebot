package action

import (
	"math/rand"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	img = "file:///root/Lucy_Project/memes/"
)

var (
	limit = rate.NewManager[int64](time.Minute*10, 15)
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
				ctx.SendChain(randtexts("?", "？", "嗯？", "(。´・ω・)ん?", "ん？"))
			case 1, 2:
				ctx.SendChain(randImage("6148451828070.jpg", "2929073585339.gif"))
			}
		})

	engine.OnFullMatch("草").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			switch rand.Intn(2) {
			case 0:
				ctx.SendChain(randtexts("（一种植物）", "ん？"))
			case 1, 2:
				ctx.SendChain(randImage("5cee2a0f5dc32a114b1a9d3f60314e5e.jpg", "R-C.jpeg", "26329490470319.jpg"))
			}
		})

	engine.OnFullMatch("哈哈哈").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			process.SleepAbout1sTo2s()
			ctx.SendChain(message.Text("草"))
		})

	engine.OnKeyword("离谱").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			switch rand.Intn(3) {
			case 0:
				ctx.SendChain(randtexts("?", "？", "ん？"))
			case 1, 2:
				ctx.SendChain(randImage("-33ee3a0711f11810.jpg", "2929073585339.gif"))
			}
		})
	engine.OnFullMatchGroup([]string{"呜呜", "呜呜呜"}).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			switch rand.Intn(2) {
			case 0:
				ctx.SendChain(randImage("2925511468257.png", "-6c72f212a4f62980.jpg"), message.Text(
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
				ctx.SendChain(randtexts("喵喵~", "喵w~"))
			case 1:
				ctx.SendChain(randImage("6152277811454.jpg", "26329403092194.gif"))
			}
		})
	engine.OnFullMatch("咕咕").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			ctx.SendChain(randtexts("抓到一只鸽子OwO", "是鸽子 炖了~", "咕咕咕", "不许咕咕咕"))
		})

	engine.OnFullMatch("小鱼干").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			ctx.SendChain(randImage("6121912482976.gif", "26329316800111.jpg", "26329403092194.gif"), message.Text("我要恰！"))
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
}
func randtexts(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(img + file[rand.Intn(len(file))])
}
