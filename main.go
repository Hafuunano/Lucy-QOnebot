package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/FloatTech/ZeroBot-Plugin/kanban"        // 在最前打印 banner
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/atri" // atri
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chat" // 回复

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/manager" // 群管

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/action" // action For Lucy触发

	_ "github.com/FloatTech/zbputils/job" // 定时指令触发器

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/bilibili"       // 查询b站用户信息
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/funwork" // 好玩的整合工具
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/tools"   // 工具

	"github.com/FloatTech/zbputils/process"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
	// -----------------------以上为内置依赖，勿动------------------------ //
)

func init() {
	sus := make([]int64, 0, 16)
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	w := flag.Bool("w", false, "Enable warning level log and higher.")
	h := flag.Bool("h", false, "Display this help.")
	token := flag.String("t", "", "Set AccessToken of WSClient.")
	url := flag.String("u", "ws://127.0.0.1:6700", "Set Url of WSClient.")
	// 默认昵称
	adana := flag.String("n", "Lucy", "Set default nickname.")
	prefix := flag.String("p", "/", "Set command prefix.")

	flag.Parse()

	if *h {
		kanban.PrintBanner()
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	} else {
		if *d && !*w {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if *w {
			logrus.SetLevel(logrus.WarnLevel)
		}
	}

	for _, s := range flag.Args() {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, i)
	}

	// 通过代码写死的方式添加主人账号
	sus = append(sus, 1292581422)
	config.W = []*driver.WSClient{driver.NewWebSocketClient(*url, *token)}
	config.Z = zero.Config{
		NickName:      append([]string{*adana}, "Lucy", "lucy", "Lucy酱"),
		CommandPrefix: *prefix,
		SuperUsers:    sus,
		Driver:        []zero.Driver{config.W[0]},
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // 全局 seed，其他插件无需再 seed
	zero.OnFullMatchGroup([]string{".help"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(kanban.Banner))
		})
	zero.RunAndBlock(config.Z, process.GlobalInitMutex.Unlock)
}
