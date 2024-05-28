// Package main for Lucy(HiMoYo Version)
package main

import (
	"github.com/MoYoez/Lucy-QOnebot/box/whitelist"
	"time"

	"github.com/MoYoez/Lucy-QOnebot/box/notify"

	"github.com/joho/godotenv"

	_ "github.com/MoYoez/Lucy-QOnebot/functions/choose"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/daily"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/group"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/interaction"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/mai"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/manager"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/nsfw"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/pgr"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/score"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/setu"

	_ "github.com/MoYoez/Lucy-QOnebot/functions/tools"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/whitelist"
	_ "github.com/MoYoez/Lucy-QOnebot/functions/wife"

	_ "github.com/MoYoez/Lucy-QOnebot/functions/slash"

	_ "github.com/MoYoez/Lucy-QOnebot/functions/simai"

	"github.com/FloatTech/floatbox/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var WhiteListMap []int64

func init() {
	// 解析命令行参数
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	WhiteListMap = make([]int64, 0)
	WhiteListMap = whitelist.WhiteListHandler()
}

func main() {
	zero.OnMessage().SetBlock(false).Handle(func(ctx *zero.Ctx) {
		var newGroup int64
		newGroup = 0
		for _, data := range WhiteListMap {
			if data == ctx.Event.GroupID {
				newGroup = data
			}
		}
		if newGroup == 0 {
			ctx.Block()
		}
	})
	zero.OnFullMatchGroup([]string{".help", "/help"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(notify.Banner))
		})
	zero.RunAndBlock(&zero.Config{
		NickName:       append([]string{"Lucy"}),
		CommandPrefix:  "/",
		SuperUsers:     []int64{1292581422},
		Driver:         []zero.Driver{driver.NewWebSocketClient("ws://127.0.0.1:6700", "")},
		MaxProcessTime: time.Minute * 4,
		RingLen:        0,
	}, process.GlobalInitMutex.Unlock)
}
