// Package interaction is mainly for Lucy's Chat base.
package interaction

import (
	"fmt"
	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/MoYoez/Lucy-QOnebot/box/event"
	"github.com/MoYoez/Lucy-QOnebot/box/setname"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	poke = rate.NewManager[int64](time.Minute*10, 8)
	img  = "file:///root/Lucy_Project/memes/"

	engine = control.Register("interaction", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "interaction",
	})
)

func init() {
	engine.OnRegex(`叫我.*?(.*)`, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		onRegexMessage := ctx.State["regex_matched"].([]string)[1]
		if setname.StringInArray(onRegexMessage, []string{"Lucy", "笨蛋", "老公", "猪", "夹子", "主人"}) {
			ctx.Send(message.Text("这些名字可不好哦(敲)"))
			return
		}

		if onRegexMessage == "" {
			ctx.Send(message.Text("好哦~ 那~咱该叫你什么呢ww"))
			nextStep := event.WaitForNextMessage(ctx)
			if nextStep.String() == "" {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("aw咱不知道你要叫什么（ 溜了x"))
				return
			} else {
				onRegexMessage = nextStep.String()
			}
		}
		if strings.Contains(onRegexMessage, "[CQ") {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("嗯哼~咱知道你在做什么哦"))
			return
		}
		userID := strconv.FormatInt(ctx.Event.UserID, 10)
		err := setname.StoreUserNickname(userID, onRegexMessage)
		if err != nil {
			ctx.Send(message.Text("呜呜呜，出现错误了...ERROR: ", err))
			return
		}
		ctx.Send(message.Text("好哦~ ", onRegexMessage, "是一个好听的名字！"))
	})
	engine.OnFullMatch("", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		nickname := zero.BotConfig.NickName[0]
		ctx.SendChain(message.Text(
			[]string{
				"这里是" + nickname + "(っ●ω●)っ",
				nickname + "不在呢~",
				"哼！" + nickname + "不想理你~",
			}[rand.Intn(3)],
		))
		process.SleepAbout1sTo2s()
		ctx.Send(message.Poke(ctx.Event.UserID))
	})
	engine.On("notice/notify/poke", zero.OnlyToMe).SetBlock(false).Handle(func(ctx *zero.Ctx) {
		nickname := zero.BotConfig.NickName[0]
		switch {
		case poke.Load(ctx.Event.GroupID).AcquireN(3):
			ctx.SendChain(message.Text([]string{"请不要戳" + nickname + " >_<", "再戳也不会理你的哦！", "别以为人家会搭理哦！",
				"呜…别戳了…", "别戳了！", "喵~", "有笨蛋在戳我，我不说是谁", "达咩呦，达咩达咩", "哼!不许戳啦 大笨蛋", "别戳啦！", "有笨蛋~让咱看看是谁"}[rand.Intn(11)]), message.Image([]string{img + "2941750127783.jpg", img + "C(185HMG2G0FY`3~2_[_H)W.gif", img + "file_3491851.jpg", img + "file_3492326.jpg", img + "file_3492330.jpg", img + "load.jpg"}[rand.Intn(6)]))
		case poke.Load(ctx.Event.GroupID).Acquire():
			happyFew := fmt.Sprintf("（好感 - %d）", rand.Intn(100)+1)
			ctx.SendChain(message.Text("喂(#`O′) 戳", nickname, "干嘛！", happyFew))
			process.SleepAbout1sTo2s()
			ctx.Send(message.Poke(ctx.Event.UserID))
		}
	})
	engine.OnFullMatch("敲我", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(randImage("797198491dc98e4f.jpg", "file_3492319.jpg", "file_3492325.jpg"))
		process.SleepAbout1sTo2s()
		ctx.Send(randImage("6170420371656.jpg"))
		process.SleepAbout1sTo2s()
		ctx.SetGroupBan(
			ctx.Event.GroupID,
			ctx.Event.UserID,
			1*60)
		process.SleepAbout1sTo2s()
		ctx.SetGroupBan(
			ctx.Event.GroupID,
			ctx.Event.UserID,
			0)
	})
	engine.OnFullMatchGroup([]string{"戳我", "戳戳"}, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		trueOrNot := rand.Intn(100)
		if trueOrNot >= 50 {
			ctx.SendChain(message.At(ctx.Event.UserID), message.Text("好哦w"))
			process.SleepAbout1sTo2s()
			ctx.Send(message.Poke(ctx.Event.UserID))
		} else {
			process.SleepAbout1sTo2s()
			ctx.Send(message.Text("哼！Lucy才不想戳"))
		}
	})
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(img + file[rand.Intn(len(file))])
}
