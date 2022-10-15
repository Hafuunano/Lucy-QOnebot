// Package chat basicchat
package chat

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"

	name "github.com/FloatTech/ZeroBot-Plugin/dependence/name"
)

type kimo = map[string]*[]string

var (
	poke   = rate.NewManager[int64](time.Minute*5, 8)  // 戳一戳
	limit  = rate.NewManager[int64](time.Minute*3, 28) // 回复限制
	img    = "file:///root/Lucy_Project/memes/"
	engine = control.Register("chat", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "chat\n- [BOT名字]\n- [戳一戳BOT]\n- 空调开\n- 空调关\n- 群温度\n- 设置温度[正整数]",
	})
)

func init() { // 插件主体
	go func() {
		data, err := os.ReadFile(engine.DataFolder() + "kimoi.json")
		if err != nil {
			panic(err)
		}
		kimomap := make(kimo, 256)
		err = json.Unmarshal(data, &kimomap)
		if err != nil {
			panic(err)
		}
		chatList := make([]string, 0, 256)
		for k := range kimomap {
			chatList = append(chatList, k)
		}
		engine.OnRegex(`叫我.*?(.*)`, zero.OnlyToMe).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
			texts := ctx.State["regex_matched"].([]string)[1]
			if name.StringInArray(texts, []string{"Lucy", "笨蛋", "老公", "猪", "夹子", "主人"}) {
				ctx.Send(message.Text("这些名字可不好哦(敲)"))
				return
			}
			if texts == "" {
				ctx.Send(message.Text("好哦~ 那~咱该叫你什么呢ww"))
				nextstep := ctx.FutureEvent("message", ctx.CheckSession())
				recv, cancel := nextstep.Repeat()
				for i := range recv {
					texts := i.MessageString()
					if texts != "" {
						cancel()
						continue
					} else {
						return
					}
				}
			}
			userID := strconv.FormatInt(ctx.Event.UserID, 10)
			err = name.StoreUserNickname(userID, texts)
			if err != nil {
				ctx.Send(message.Text("发生了一些不可预料的问题 请稍后再试,ERR: ", err))
			}
			ctx.Send(message.Text("好哦~ ", texts, " ちゃん~~~"))
		})
		engine.OnFullMatchGroup(chatList, zero.OnlyToMe).SetBlock(true).Handle(
			func(ctx *zero.Ctx) {
				switch {
				case limit.Load(ctx.Event.UserID).AcquireN(3):
					key := ctx.MessageString()
					val := *kimomap[key]
					text := val[rand.Intn(len(val))]
					userID := strconv.FormatInt(ctx.Event.UserID, 10)
					userNickName := name.LoadUserNickname(userID)
					result := strings.ReplaceAll(text, "你", userNickName)
					process.SleepAbout1sTo2s()
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(result)) // 来自于 https://github.com/Kyomotoi/AnimeThesaurus 的回复 经过二次修改
				case limit.Load(ctx.Event.UserID).Acquire():
					process.SleepAbout1sTo2s()
					ctx.Send(message.Text("咱不想说话~好累qwq"))
					return
				default:
				}
			})
	}()
	// 被喊名字
	engine.OnFullMatch("", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var nickname = zero.BotConfig.NickName[0]
			time.Sleep(time.Second * 1)
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

	// 戳一戳
	engine.On("notice/notify/poke", zero.OnlyToMe).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			var nickname = zero.BotConfig.NickName[0]
			switch {
			case poke.Load(ctx.Event.GroupID).AcquireN(3):
				// 5分钟共8块命令牌 一次消耗3块命令牌
				time.Sleep(time.Second * 1)
				ctx.SendChain(message.Text([]string{"请不要戳" + nickname + " >_<", "再戳也不会理你的哦！", "别以为人家会搭理哦！",
					"呜…别戳了…", "别戳了！", "喵~", "有笨蛋在戳我，我不说是谁", "达咩呦，达咩达咩", "好怪..你不要过来啊啊啊啊啊", "别戳啦！"}[rand.Intn(10)]))
			case poke.Load(ctx.Event.GroupID).Acquire():
				// 5分钟共8块命令牌 一次消耗1块命令牌
				time.Sleep(time.Second * 1)
				happyFew := fmt.Sprintf("（好感 - %d）", rand.Intn(100)+1)
				ctx.SendChain(message.Text("喂(#`O′) 戳", nickname, "干嘛！", happyFew))
				process.SleepAbout1sTo2s()
				ctx.Send(message.Poke(ctx.Event.UserID))
			default:
				// 频繁触发，不回复
			}
		})
	// 戳我
	engine.OnFullMatchGroup([]string{"戳我", "戳戳"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			trueOrNot := rand.Intn(100)
			if trueOrNot >= 50 {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("好哦w"))
				process.SleepAbout1sTo2s()
				ctx.Send(message.Poke(ctx.Event.UserID))
			} else {
				process.SleepAbout1sTo2s()
				ctx.Send(message.Text("哼！Lucy才不想戳你"))
			}
		})

	// 群空调
	var AirConditTemp = map[int64]int{}
	var AirConditSwitch = map[int64]bool{}
	engine.OnFullMatch("空调开").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			AirConditSwitch[ctx.Event.GroupID] = true
			ctx.SendChain(message.Text("❄️哔~"))
		})
	engine.OnFullMatch("空调关").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			AirConditSwitch[ctx.Event.GroupID] = false
			delete(AirConditTemp, ctx.Event.GroupID)
			ctx.SendChain(message.Text("💤哔~"))
		})
	engine.OnRegex(`设置温度(\d+)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if _, exist := AirConditTemp[ctx.Event.GroupID]; !exist {
				AirConditTemp[ctx.Event.GroupID] = 26
			}
			if AirConditSwitch[ctx.Event.GroupID] {
				temp := ctx.State["regex_matched"].([]string)[1]
				AirConditTemp[ctx.Event.GroupID], _ = strconv.Atoi(temp)
				ctx.SendChain(message.Text(
					"❄️风速中", "\n",
					"群温度 ", AirConditTemp[ctx.Event.GroupID], "℃",
				))
			} else {
				ctx.SendChain(message.Text(
					"💤", "\n",
					"群温度 ", AirConditTemp[ctx.Event.GroupID], "℃",
				))
			}
		})
	engine.OnFullMatch(`群温度`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if _, exist := AirConditTemp[ctx.Event.GroupID]; !exist {
				AirConditTemp[ctx.Event.GroupID] = 26
			}
			if AirConditSwitch[ctx.Event.GroupID] {
				ctx.SendChain(message.Text(
					"❄️风速中", "\n",
					"群温度 ", AirConditTemp[ctx.Event.GroupID], "℃",
				))
			} else {
				ctx.SendChain(message.Text(
					"💤", "\n",
					"群温度 ", AirConditTemp[ctx.Event.GroupID], "℃",
				))
			}
		})
	// 针对于自己的添加w 我很懒不要打我a.a
	engine.OnKeywordGroup([]string{"会什么"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText(
				"麻烦看一下说明书惹 群内发送lucy.help即可w",
			))
		})
	engine.OnFullMatch("捏脸", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText("大笨蛋！不许捏٩(๑`^´๑)۶", "疼....不许这样！哼！"), randImage("26329371069850.jpg"))
		})
	engine.OnFullMatch("摸头", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("6126814446620.jpg", "kawaii.jpg"), randText("咱超可爱的w"))
		})
	engine.OnFullMatch("敲我", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("797198491dc98e4f.jpg"))
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
	engine.OnFullMatch("酱", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText("Lucy酱在这边~", "略略略~这边是Lucy(*/ω＼*)", "Lucy在忙哦w 有什么事情嘛"))
		})
	engine.OnFullMatch("摸摸", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText("啾啾~", "呼呼~", "摸摸~"), randImage("22b530369f3c0fdd.jpg"))
		})
	engine.OnFullMatchGroup([]string{"呼呼", "抱抱"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("26329502616465.jpg", "sleep.jpg"))
		})
	engine.OnFullMatch("抱住", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("22b530369f3c0fdd.jpg", "6126814446620.jpg", "wwwss.jpg"))
		})
	engine.OnFullMatch("举高高", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			process.SleepAbout1sTo2s()
			ctx.SendChain(message.Text("哼！才不让举高高呢"), randImage("dcf07a381f30e9240bf68c845b086e061c95f72a.jpg"))
		})
	engine.OnKeywordGroup([]string{"MoYoez", "MoeMagicMango"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randText("这些都是主人的名字~", "这些是架子(划掉)使用的名字w"))
		})
	engine.OnFullMatchGroup([]string{"是笨蛋"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("sasda.jpg", "55D0B4A5E335FE55A924E71469F35AC7.png"))
		})
}
func randText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(img + file[rand.Intn(len(file))])
}
