package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/FloatTech/ZeroBot-Plugin/kanban"           // 在最前打印 banner
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/bottle"  // 漂流瓶
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/fortune" // fortune
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/manager" // 群管
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nsfw"    // nsfw
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/tools"   // 工具

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/atri" // atri

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/action" // action For Lucy触发

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/funwork" // 好玩的整合工具

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/score" // 签到

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/arc" // arcaea

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chat" // 回复

	"github.com/FloatTech/floatbox/process"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type zbpcfg struct {
	Z zero.Config        `json:"zero"`
	W []*driver.WSClient `json:"ws"`
}

var config zbpcfg

func init() {
	// 解析命令行参数
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// make it easy to write config rather than change the source code.
	sus := make([]int64, 0, 16)
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	w := flag.Bool("w", false, "Enable warning level log and higher.")
	h := flag.Bool("h", false, "Display this help.")
	// 直接写死 AccessToken 时，请更改下面第二个参数
	token := flag.String("t", "", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	url := flag.String("u", "ws://127.0.0.1:6700", "Set Url of WSClient.")
	// 默认昵称
	adana := flag.String("n", os.Getenv("name"), "Set default nickname.")
	prefix := flag.String("p", "/", "Set command prefix.")
	late := flag.Uint("l", 1000, "Response latency (ms).")
	rsz := flag.Uint("r", 4096, "Receiving buffer ring size.")
	maxpt := flag.Uint("x", 4, "Max process time (min).")
	flag.Parse()

	if *h {
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
		NickName:       append([]string{*adana}, "Lucy", "lucy", "Lucy酱"),
		CommandPrefix:  *prefix,
		SuperUsers:     sus,
		RingLen:        *rsz,
		Latency:        time.Duration(*late) * time.Millisecond,
		MaxProcessTime: time.Duration(*maxpt) * time.Minute,
		Driver:         []zero.Driver{config.W[0]},
	}
}

func main() {
	zero.OnFullMatchGroup([]string{".help", "帮助", "/help"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(kanban.Banner))
		})
	zero.RunAndBlock(&config.Z, process.GlobalInitMutex.Unlock)
}
