package GroupForMe

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"

	"github.com/FloatTech/ZeroBot-Plugin/order"
	control "github.com/FloatTech/zbputils/control"
)
const (
    
    img = "file:///root/qqbot/img/"
    
    )


func init() {
	engine := control.Register("GroupForMe", order.Priomodernok, &control.Options{
		DisableOnDefault: false,
		Help:             "给自己定制的 \n- 我也不知道写什么.jpg",
	})

	//针对于自己的添加w 我很懒不要打我a.a
	engine.OnKeywordGroup([]string{"主人"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts(
				"主人是MoeMagicMango",
				"夹子酱",
				"ww~你猜嘛www",
				"大笨蛋~是夹子惹ww是(≧∇≦)ﾉ",
				"架子~Σ( ° △ °|||)︴说错辣!",
				"看一下~柠檬味的布丁盒子叭w",
			))
		})
		
	engine.OnKeywordGroup([]string{"柠檬味的布丁盒子"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts(
				"~呼呼，你猜一下w",
				"这个是主人的博客~",
				"据说这个博客使用了世界最先进的非对称加密安全标准~(指HTTPS)",
				"麻烦看到架子没写博客的话 把他拍死>_<",
			))
		})
	engine.OnKeywordGroup([]string{"名字"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts(
				"~呼呼，你猜一下w",
				"来源于PicoPico~Magic~",
				"或许你可以去问一下夹子酱嗷w",
			))
		})
	engine.OnKeywordGroup([]string{"会什么"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts(
				"麻烦看一下说明书惹 输入Hana.help即可w",
			))
		})
	engine.OnKeywordGroup([]string{"会csgo"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts(
				"我不会！主人会ww",
				"你可以试着让夹子玩玩w，他玩的可比我菜hdndkjdjejsjjdbdh",
				"什么是csgo? 你可以教我嘛ww",
				"或许哪天就可以一起玩了嘛~",
			))
		})
	engine.OnFullMatchGroup([]string{"使用方法"},zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts(
				"https://github.com/MoYoez/Hana-Introduction/tree/main/intro/readme.md",
				"我会的可多了啦www~~~~自己猜猜哦♪(^∇^*)",
				"我会什么都是由夹子的开发能力绝对惹.jpg",
				"提醒一下~架子是个大笨蛋 只会瞎改惹www",
			))
		})


     engine.OnKeyword("老婆",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
		    switch rand.Intn(5){
		        case 0:
				ctx.SendChain(randtexts("?", "你在说什么?", "嗯？", "(。´・ω・)ん?", "是笨蛋?"))
			    case 1, 2:
			    ctx.SendChain(randImage("2984862825214.jpg","2977250804540.jpg","2972603024125.jpg","2932544972786.jpg","2948706690280.jpg"))
		    }
		})


engine.OnFullMatch("捏脸",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts("大笨蛋！不许捏٩(๑`^´๑)۶","疼....","(捏你的脸~)",))
		})
engine.OnFullMatch("摸头",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("2953919819601.jpg",))
		})
engine.OnFullMatch("Hana酱",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts("Hana酱在这边~","略略略~这边是Hana(*/ω＼*)","Hana在忙哦w 有什么事情嘛",))
		})
engine.OnFullMatch("摸摸",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts("啾啾~","呼呼~","摸摸你~",))
		})
engine.OnKeyword("教我",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text("你可以去问一下夹子哦w","如果夹子在的话~你可以去问一下他哦w",))
		})


engine.OnKeywordGroup([]string{"MoYoez","MoeMagicMango"},zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts("这些都是主人的名字~","这些是架子(划掉)使用的名字w",))
		})		
		
engine.OnKeyword("人设",zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randtexts("人设的话...目前夹子没打算做哦","要不你帮我画一个叭(雾)",))
		})		
		
}
func randtexts(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(img + file[rand.Intn(len(file))])
}
