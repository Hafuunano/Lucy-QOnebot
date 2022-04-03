package action

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"

	"github.com/FloatTech/ZeroBot-Plugin/order"
	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/process"
)
const (
    
    img = "file:///root/qqbot/img/"
    
    )

	func init() {
		engine := control.Register("action", order.Prioaction, &control.Options{
			DisableOnDefault: true,
			Help:             "Hana容易被动触发语言 \n- 默认禁用 可以自行打开哦",
		})
	//Okk

	engine.OnFullMatchGroup([]string{"？", "?", "¿"}).SetBlock(true).
	Handle(func(ctx *zero.Ctx) {
		process.SleepAbout1sTo2s()
		switch rand.Intn(5) {
			case 0:
			ctx.SendChain(randtexts("?", "？", "嗯？", "(。´・ω・)ん?", "ん？"))
			case 1, 2:
			ctx.SendChain(randImage("2989062593389.jpg","2929073585339.gif"))
		}
	})

	engine.OnKeyword("离谱").SetBlock(true).
	Handle(func(ctx *zero.Ctx) {
		switch rand.Intn(5){
			case 0:
			ctx.SendChain(randtexts("?", "？", "嗯？", "(。´・ω・)ん?", "ん？"))
			case 1, 2:
			ctx.SendChain(randImage("3003142120311.jpg"))
		}
	})
	engine.OnFullMatchGroup([]string{"呜呜","呜呜呜"}).SetBlock(true).
	Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(randImage("2925511468257.png") , message.Text(
			"摸摸~"))
	})
	engine.OnFullMatch("你好").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("-1cc4a39a895d93ae.jpg"))
		})
	engine.OnFullMatch("都怪Zheic").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("-122f9251d280c34c.jpg"))
		})
		engine.OnFullMatchGroup([]string{"喵","喵喵","喵喵喵"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts("喵喵~","喵w~",))
		})


		//end

	}
	func randtexts(text ...string) message.MessageSegment {
		return message.Text(text[rand.Intn(len(text))])
	}
	
	func randImage(file ...string) message.MessageSegment {
		return message.Image(img + file[rand.Intn(len(file))])
	}
	