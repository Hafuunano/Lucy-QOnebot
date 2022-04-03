package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	// 注：以下插件均可通过前面加 // 注释，注释后停用并不加载插件
	// 下列插件可与 wdvxdr1123/ZeroBot v1.1.2 以上配合单独使用

	// 词库类
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_atri" // ATRI词库
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_chat" // 基础词库
	// 实用类

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_github"       // 搜索GitHub仓库
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_manager"      // 群管
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_runcode"      // 在线运行代码
	// 娱乐类

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_ai_false"  // 服务器监控
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_choose"    // 选择困难症帮手
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_curse"     // 骂人
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_gif"       // 制图
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_hs"        // 炉石
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_music"     // 点歌
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_reborn"    // 投胎
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_score"     // 分数
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_shadiao"   // 沙雕app
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_wangyiyun" // 网易云音乐热评

	// b站相关
//	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_bilibili"       // 查询b站用户信息
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_bilibili_parse" // b站视频链接解析
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_diana" // 嘉心糖发病

	// 二次元图片
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_aiwife"     // 随机老婆
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_danbooru"   // DeepDanbooru二次元图标签识别
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_lolicon"    // lolicon 随机图片
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_setu" // 本地涩图
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_nativewife" // 本地老婆
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_saucenao"   // 以图搜图

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_tracemoe"      // 搜番
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_vtb_quotation" // vtb语录

	// 定制化设计
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_groupforme" // 我自己+botinfo
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_action"   //Keyword Action反触发.
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin_ramdomimg" //随机图片

	// 以下为内置依赖，勿动
	"github.com/FloatTech/ZeroBot-Plugin/order"
	"github.com/fumiama/go-registry"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	contents = []string{
		"* Hana基于 Onebot + ZeroBot + Golang 原生实现的bot",
		"* 目前运营为 Morning Call(MoYoez) ",
		"* 本项目经本人基于1.2.5 beta3 fixed.编译修改",
		"* Project: https://github.com/FloatTech/ZeroBot-Plugin",
		"* 说明书: https://hana.himoyo.cn",
		"* Copyright © 2021-2022 FloatTech. All Rights Reserved.",
	}
	nicks  = []string{"hana", "Hana酱", "Hana"}
	banner = strings.Join(contents, "\n")
	token  *string
	url    *string
	adana  *string
	prefix *string
	reg    = registry.NewRegReader("reilia.fumiama.top:32664", "fumiama")
)

func init() {
	// 解析命令行参数
	d := flag.Bool("d", false, "Enable debug level log and higher.")
	w := flag.Bool("w", false, "Enable warning level log and higher.")
	h := flag.Bool("h", false, "Display this help.")
	// 解析命令行参数，输入 `-g 监听地址:端口` 指定 gui 访问地址，默认 127.0.0.1:3000
	// g := flag.String("g", "127.0.0.1:3000", "Set web gui listening address.")

	// 直接写死 AccessToken 时，请更改下面第二个参数
	token = flag.String("t", "", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	url = flag.String("u", "ws://127.0.0.1:6700", "Set Url of WSClient.")
	// 默认昵称
	adana = flag.String("n", "Hana", "Set default nickname.")
	prefix = flag.String("p", "/", "Set command prefix.")

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

}

func getKanban() string {
	err := reg.Connect()
	if err != nil {
		return err.Error()
	}
	defer reg.Close()
	text, err := reg.Get("ZeroBot-Plugin/kanban")
	if err != nil {
		return err.Error()
	}
	return text
}

func main() {
	order.Wait()

	// 帮助
	zero.OnFullMatchGroup([]string{".gohelp"},zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(banner))
		})

	zero.RunAndBlock(
		zero.Config{
			NickName:      append([]string{*adana}, nicks...),
			CommandPrefix: *prefix,
			// SuperUsers 某些功能需要主人权限，可通过以下两种方式修改
			SuperUsers: []string{"1292581422"}, // 通过代码写死的方式添加主人账号
			//	SuperUsers: flag.Args(), // 通过命令行参数的方式添加主人账号
			Driver: []zero.Driver{driver.NewWebSocketClient(*url, *token)},
		},
	)
}
