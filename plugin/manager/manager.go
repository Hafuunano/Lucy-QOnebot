// Package manager 群管
package manager

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/floatbox/math"
	sql "github.com/FloatTech/sqlite"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"

	"github.com/FloatTech/ZeroBot-Plugin/plugin/manager/timer"
)

const (
	hint = "====群管====\n" +
		"- 禁言@QQ 1分钟\n" +
		"- 解除禁言 @QQ\n" +
		"- 我要自闭 1分钟\n" +
		"- 开启全员禁言\n" +
		"- 解除全员禁言\n" +
		"- 升为管理@QQ\n" +
		"- 取消管理@QQ\n" +
		"- 修改名片@QQ XXX\n" +
		"- 修改头衔@QQ XXX\n" +
		"- 申请头衔 XXX\n" +
		"- 踢出群聊@QQ\n" +
		"- 退出群聊 1234@bot\n" +
		"- 群聊转发 1234 XXX\n" +
		"- 私聊转发 0000 XXX\n" +
		"- 在MM月dd日的hh点mm分时(用https://url)提醒大家XXX\n" +
		"- 在MM月[每周 | 周几]的hh点mm分时(用https://url)提醒大家XXX\n" +
		"- 取消在MM月dd日的hh点mm分的提醒\n" +
		"- 取消在MM月[每周 | 周几]的hh点mm分的提醒\n" +
		"- 在\"cron\"时(用[url])提醒大家[xxx]\n" +
		"- 取消在\"cron\"的提醒\n" +
		"- 列出所有提醒\n" +
		"- 翻牌\n" +
		"- 设置欢迎语XXX 可选添加 [{at}] [{nickname}] [{avatar}] [{id}] {at}可在发送时艾特被欢迎者 {nickname}是被欢迎者名字 {avatar}是被欢迎者头像 {id}是被欢迎者QQ号\n" +
		"- 测试欢迎语\n" +
		"- 设置告别辞 参数同设置欢迎语\n" +
		"- 测试告别辞\n" +
		"- [开启 | 关闭]入群验证"
)

var (
	db    = &sql.Sqlite{}
	clock timer.Clock
)

func init() { // 插件主体
	engine := control.Register("manager", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              hint,
		PrivateDataFolder: "manager",
	})

	go func() {
		db.DBPath = engine.DataFolder() + "config.db"
		err := db.Open(time.Hour * 24)
		if err != nil {
			panic(err)
		}
		clock = timer.NewClock(db)
		err = db.Create("welcome", &welcome{})
		if err != nil {
			panic(err)
		}
		err = db.Create("member", &member{})
		if err != nil {
			panic(err)
		}
		err = db.Create("farewell", &welcome{})
		if err != nil {
			panic(err)
		}
	}()

	// 升为管理
	engine.OnRegex(`^升为管理.*?(\d+)`, zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupAdmin(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被升为管理的人的qq
				true,
			)
			nickname := ctx.GetGroupMemberInfo( // 被升为管理的人的昵称
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被升为管理的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text(nickname + " 升为了管理~"))
		})

	engine.OnFullMatch("打我", zero.OnlyToMe).SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			time.Sleep(time.Second * 2)
			ctx.SendChain(message.Text("给你一拳~"))
			time.Sleep(time.Second * 1)
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				1*60)
			time.Sleep(time.Second * 1)
			ctx.SendChain(message.Text("两拳~"))
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				4*60)
			time.Sleep(time.Second * 1)
			ctx.SendChain(message.Text("三拳~"))
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				8*60)
			time.Sleep(time.Second * 1)
			ctx.SendChain(message.Text("够了嘛 bakaꐦ≖ ≖"))
		})

	// 取消管理
	engine.OnRegex(`^取消管理.*?(\d+)`, zero.OnlyGroup, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupAdmin(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被取消管理的人的qq
				false,
			)
			nickname := ctx.GetGroupMemberInfo( // 被取消管理的人的昵称
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被取消管理的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text("残念~ " + nickname + " 暂时失去了管理员的资格"))
		})
	// 踢出群聊
	engine.OnRegex(`^踢出群聊.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupKick(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			)
			nickname := ctx.GetGroupMemberInfo( // 被踢出群聊的人的昵称
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被踢出群聊的人的qq
				false,
			).Get("nickname").Str
			ctx.SendChain(message.Text("残念~ " + nickname + " 被放逐"))
		})
	// 退出群聊
	engine.OnRegex(`^退出群聊.*?(\d+)`, zero.OnlyToMe, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupLeave(
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要退出的群的群号
				true,
			)
		})
	// 开启全体禁言
	engine.OnRegex(`^开启全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupWholeBan(
				ctx.Event.GroupID,
				true,
			)
			ctx.SendChain(message.Text("全员自闭开始~"))
		})
	// 解除全员禁言
	engine.OnRegex(`^解除全员禁言$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupWholeBan(
				ctx.Event.GroupID,
				false,
			)
			ctx.SendChain(message.Text("全员自闭结束~"))
		})
	// 禁言
	engine.OnRegex(`^禁言.*?(\d+).*?\s(\d+)(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			duration := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch ctx.State["regex_matched"].([]string)[3] {
			case "分钟":
				//
			case "小时":
				duration *= 60
			case "天":
				duration *= 60 * 24
			default:
				//
			}
			if duration >= 43200 {
				duration = 43199 // qq禁言最大时长为一个月
			}
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要禁言的人的qq
				duration*60, // 要禁言的时间（分钟）
			)
			ctx.SendChain(message.Text("小黑屋收留成功~"))
		})
	// 解除禁言
	engine.OnRegex(`^解除禁言.*?(\d+)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 要解除禁言的人的qq
				0,
			)
			ctx.SendChain(message.Text("小黑屋释放成功~"))
		})
	// 自闭禁言
	engine.OnRegex(`^(我要自闭|禅定).*?(\d+)(.*)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			duration := math.Str2Int64(ctx.State["regex_matched"].([]string)[2])
			switch ctx.State["regex_matched"].([]string)[3] {
			case "分钟", "min", "mins", "m":
				duration *= 1
			case "小时", "hour", "hours", "h":
				duration *= 60
			case "天", "day", "days", "d":
				duration *= 60 * 24
			default:
				duration *= 0
			}
			if duration == 114514 {
				ctx.Send(message.Text("好臭的数字a.a 不许用."))
			}
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				duration*60, // 要自闭的时间（分钟）
			)
			ctx.SendChain(message.Text("那我就不手下留情了~"))
		})
	// 修改名片
	engine.OnRegex(`^修改名片.*?(\d+).*?\s(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupCard(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被修改群名片的人
				ctx.State["regex_matched"].([]string)[2],                 // 修改成的群名片
			)
			ctx.SendChain(message.Text("嗯！已经修改了"))
		})
	// 修改头衔
	engine.OnRegex(`^修改头衔.*?(\d+).*?\s(.*)`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupSpecialTitle(
				ctx.Event.GroupID,
				math.Str2Int64(ctx.State["regex_matched"].([]string)[1]), // 被修改群头衔的人
				ctx.State["regex_matched"].([]string)[2],                 // 修改成的群头衔
			)
			ctx.SendChain(message.Text("嗯！已经修改了"))
		})
	// 申请头衔
	engine.OnRegex(`^申请头衔(.*)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SetGroupSpecialTitle(
				ctx.Event.GroupID,
				ctx.Event.UserID,                         // 被修改群头衔的人
				ctx.State["regex_matched"].([]string)[1], // 修改成的群头衔
			)
			ctx.SendChain(message.Text("嗯！不错的头衔呢~"))
		})
	// 定时提醒
	engine.OnRegex(`^在(.{1,2})月(.{1,3}日|每?周.?)的(.{1,3})点(.{1,3})分时(用.+)?提醒大家(.*)`, zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dateStrs := ctx.State["regex_matched"].([]string)
			ts := timer.GetFilledTimer(dateStrs, ctx.Event.SelfID, ctx.Event.GroupID, false)
			if ts.En() {
				go clock.RegisterTimer(ts, true, false)
				ctx.SendChain(message.Text("记住了~"))
			} else {
				ctx.SendChain(message.Text("参数非法:" + ts.Alert))
			}
		})
	// 定时 cron 提醒
	engine.OnRegex(`^在"(.*)"时(用.+)?提醒大家(.*)`, zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dateStrs := ctx.State["regex_matched"].([]string)
			var url, alert string
			switch len(dateStrs) {
			case 4:
				url = strings.TrimPrefix(dateStrs[2], "用")
				alert = dateStrs[3]
			case 3:
				alert = dateStrs[2]
			default:
				ctx.SendChain(message.Text("参数非法!"))
				return
			}
			logrus.Debugln("[manager] cron:", dateStrs[1])
			ts := timer.GetFilledCronTimer(dateStrs[1], alert, url, ctx.Event.SelfID, ctx.Event.GroupID)
			if clock.RegisterTimer(ts, true, false) {
				ctx.SendChain(message.Text("记住了~"))
			} else {
				ctx.SendChain(message.Text("参数非法:" + ts.Alert))
			}
		})
	// 取消定时
	engine.OnRegex(`^取消在(.{1,2})月(.{1,3}日|每?周.?)的(.{1,3})点(.{1,3})分的提醒`, zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dateStrs := ctx.State["regex_matched"].([]string)
			ts := timer.GetFilledTimer(dateStrs, ctx.Event.SelfID, ctx.Event.GroupID, true)
			ti := ts.GetTimerID()
			ok := clock.CancelTimer(ti)
			if ok {
				ctx.SendChain(message.Text("取消成功~"))
			} else {
				ctx.SendChain(message.Text("没有这个定时器哦~"))
			}
		})
	// 取消 cron 定时
	engine.OnRegex(`^取消在"(.*)"的提醒`, zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dateStrs := ctx.State["regex_matched"].([]string)
			ts := timer.Timer{Cron: dateStrs[1], GrpID: ctx.Event.GroupID}
			ti := ts.GetTimerID()
			ok := clock.CancelTimer(ti)
			if ok {
				ctx.SendChain(message.Text("取消成功~"))
			} else {
				ctx.SendChain(message.Text("没有这个定时器哦~"))
			}
		})
	// 列出本群所有定时
	engine.OnFullMatch("列出所有提醒", zero.AdminPermission, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(clock.ListTimers(ctx.Event.GroupID)))
		})
	// 随机点名
	engine.OnFullMatchGroup([]string{"翻牌"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			// 无缓存获取群员列表
			list := ctx.CallAction("get_group_member_list", zero.Params{
				"group_id": ctx.Event.GroupID,
				"no_cache": true,
			}).Data
			temp := list.Array()
			sort.SliceStable(temp, func(i, j int) bool {
				return temp[i].Get("last_sent_time").Int() < temp[j].Get("last_sent_time").Int()
			})
			temp = temp[math.Max(0, len(temp)-10):]
			who := temp[rand.Intn(len(temp))]
			if who.Get("user_id").Int() == ctx.Event.SelfID {
				ctx.SendChain(message.Text("幸运儿居然是我自己"))
				return
			}
			if who.Get("user_id").Int() == ctx.Event.UserID {
				ctx.SendChain(message.Text("哎呀，就是你自己了"))
				return
			}
			nick := who.Get("card").Str
			if nick == "" {
				nick = who.Get("nickname").Str
			}
			ctx.SendChain(
				message.Text(
					nick,
					" 就是你啦！",
				),
			)
		})
	// 入群欢迎
	engine.OnNotice().SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.NoticeType == "group_increase" && ctx.Event.SelfID != ctx.Event.UserID {
				var w welcome
				err := db.Find("welcome", &w, "where gid = "+strconv.FormatInt(ctx.Event.GroupID, 10))
				if err == nil {
					ctx.SendGroupMessage(ctx.Event.GroupID, message.ParseMessageFromString(welcometocq(ctx, w.Msg)))
				} else {
					ctx.SendChain(message.Text("欢迎~"))
				}
				c, ok := control.Lookup("manager")
				if ok {
					enable := c.GetData(ctx.Event.GroupID)&1 == 1
					if enable {
						uid := ctx.Event.UserID
						a := rand.Intn(100)
						b := rand.Intn(100)
						r := a + b
						ctx.SendChain(message.At(uid), message.Text(fmt.Sprintf("考你一道题：%d+%d=?\n如果60秒之内答不上来，%s就要把你踢出去了哦~", a, b, zero.BotConfig.NickName[0])))
						// 匹配发送者进行验证
						rule := func(ctx *zero.Ctx) bool {
							for _, elem := range ctx.Event.Message {
								if elem.Type == "text" {
									text := strings.ReplaceAll(elem.Data["text"], " ", "")
									ans, err := strconv.Atoi(text)
									if err == nil {
										if ans != r {
											ctx.SendChain(message.Text("答案不对哦，再想想吧~"))
											return false
										}
										return true
									}
								}
							}
							return false
						}
						next := zero.NewFutureEvent("message", 999, false, ctx.CheckSession(), rule)
						recv, cancel := next.Repeat()
						select {
						case <-time.After(time.Minute):
							cancel()
							ctx.SendChain(message.Text("拜拜啦~"))
							ctx.SetGroupKick(ctx.Event.GroupID, uid, false)
						case <-recv:
							cancel()
							ctx.SendChain(message.Text("答对啦~"))
						}
					}
				}
			}
		})
	// 退群提醒
	engine.OnNotice().SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.NoticeType == "group_decrease" {
				var w welcome
				err := db.Find("farewell", &w, "where gid = "+strconv.FormatInt(ctx.Event.GroupID, 10))
				if err == nil {
					ctx.SendGroupMessage(ctx.Event.GroupID, message.ParseMessageFromString(welcometocq(ctx, w.Msg)))
				} else {
					userid := ctx.Event.UserID
					ctx.SendChain(message.Text(ctx.CardOrNickName(userid), "(", userid, ")", "暂时离开了~\n"))
				}
			}
		})

	// 设置欢迎语
	engine.OnRegex(`^设置欢迎语([\s\S]*)$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			w := &welcome{
				GrpID: ctx.Event.GroupID,
				Msg:   ctx.State["regex_matched"].([]string)[1],
			}
			err := db.Insert("welcome", w)
			if err == nil {
				ctx.SendChain(message.Text("记住啦!"))
			} else {
				ctx.SendChain(message.Text("出错啦: ", err))
			}
		})
	// 测试欢迎语
	engine.OnFullMatch("测试欢迎语", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var w welcome
			err := db.Find("welcome", &w, "where gid = "+strconv.FormatInt(ctx.Event.GroupID, 10))
			if err == nil {
				ctx.SendGroupMessage(ctx.Event.GroupID, message.ParseMessageFromString(welcometocq(ctx, w.Msg)))
			} else {
				ctx.SendChain(message.Text("欢迎~"))
			}
		})
	// 设置告别辞
	engine.OnRegex(`^设置告别辞([\s\S]*)$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			w := &welcome{
				GrpID: ctx.Event.GroupID,
				Msg:   ctx.State["regex_matched"].([]string)[1],
			}
			err := db.Insert("farewell", w)
			if err == nil {
				ctx.SendChain(message.Text("记住啦!"))
			} else {
				ctx.SendChain(message.Text("出错啦: ", err))
			}
		})
	// 测试告别辞
	engine.OnFullMatch("测试告别辞", zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var w welcome
			err := db.Find("farewell", &w, "where gid = "+strconv.FormatInt(ctx.Event.GroupID, 10))
			if err == nil {
				ctx.SendGroupMessage(ctx.Event.GroupID, message.ParseMessageFromString(welcometocq(ctx, w.Msg)))
			} else {
				userid := ctx.Event.UserID
				ctx.SendChain(message.Text(ctx.CardOrNickName(userid), "(", userid, ")", "离开了我们..."))
			}
		})
	// 入群后验证开关
	engine.OnRegex(`^(.*)入群验证$`, zero.OnlyGroup, zero.AdminPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			option := ctx.State["regex_matched"].([]string)[1]
			c, ok := control.Lookup("manager")
			if ok {
				data := c.GetData(ctx.Event.GroupID)
				switch option {
				case "开启", "打开", "启用":
					data |= 1
				case "关闭", "关掉", "禁用":
					data &= 0x7fffffff_fffffffe
				default:
					return
				}
				err := c.SetData(ctx.Event.GroupID, data)
				if err == nil {
					ctx.SendChain(message.Text("已", option))
					return
				}
				ctx.SendChain(message.Text("出错啦: ", err))
				return
			}
			ctx.SendChain(message.Text("找不到服务!"))
		})

	engine.OnFullMatch("优质睡眠", zero.OnlyToMe, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			dyn := time.Now().Hour()
			if dyn < 23 && dyn >= 4 {
				ctx.SendChain(message.Text("还没有到时间哦~ 请11点后再来w"))
				return
			}
			if zero.AdminPermission(ctx) {
				ctx.Send(message.Text("怪哦 咱可不能禁言诶("))
				return
			} else {
				ctx.SetGroupBan(
					ctx.Event.GroupID,
					ctx.Event.UserID,
					8*3600,
				)
				ctx.SendChain(message.Text("一键提供优质睡眠服务~"))
			}
		})
}

// 传入 ctx 和 welcome格式string 返回cq格式string  使用方法:welcometocq(ctx,w.Msg)
func welcometocq(ctx *zero.Ctx, welcome string) string {
	at := "[CQ:at,qq=" + strconv.FormatInt(ctx.Event.UserID, 10) + "]"
	avatar := "[CQ:image,file=" + "https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640]"
	id := strconv.FormatInt(ctx.Event.UserID, 10)
	nickname := ctx.CardOrNickName(ctx.Event.UserID)
	cqstring := strings.ReplaceAll(welcome, "{at}", at)
	cqstring = strings.ReplaceAll(cqstring, "{nickname}", nickname)
	cqstring = strings.ReplaceAll(cqstring, "{avatar}", avatar)
	cqstring = strings.ReplaceAll(cqstring, "{id}", id)
	return cqstring
}
