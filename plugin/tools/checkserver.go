// Package aifalse 暂时只有服务器监控
package tools

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/FloatTech/zbputils/control"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() { // 插件主体
	engine.OnFullMatchGroup([]string{"检查身体", "自检", "启动自检", "系统状态"}, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(
				"* HiMoYo ECS Status,running on AliYun\n",
				"* CPU Usage: ", cpuPercent(), "%\n",
				"* RAM Usage: ", memPercent(), "%\n",
				"* DiskInfo Usage Check: ", diskPercent(), "\n",
				"  Lucyは、高性能ですから！",
			),
			)
		})

		// 作为信息推送工具使用
	engine.OnRegex(`^【PUSH】.*`, zero.SuperUserPermission).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			m, ok := control.Lookup("tools")
			if ok {
				msg := ctx.Event.Message
				zero.RangeBot(func(id int64, ctx *zero.Ctx) bool {
					for _, g := range ctx.GetGroupList().Array() {
						gid := g.Get("group_id").Int()
						if m.IsEnabledIn(gid) {
							ctx.SendGroupMessage(gid, msg)
						}
					}
					return true
				})
			}
		})

	engine.OnRegex(`给夹子留话.*?(.*)`, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			su := zero.BotConfig.SuperUsers[0]
			now := time.Unix(ctx.Event.Time, 0).Format("2006-01-02 15:04:05")
			uid := ctx.Event.UserID
			gid := ctx.Event.GroupID
			username := ctx.CardOrNickName(uid)
			botid := ctx.Event.SelfID
			botname := zero.BotConfig.NickName[0]
			rawmsg := ctx.State["regex_matched"].([]string)[1]
			rawmsg = message.UnescapeCQCodeText(rawmsg)
			msg := make(message.Message, 10)
			msg = append(msg, message.CustomNode(botname, botid, "有人留言哦w\n在"+now))
			if gid != 0 {
				groupname := ctx.GetGroupInfo(gid, true).Name
				msg = append(msg, message.CustomNode(botname, botid, "来自群聊:["+groupname+"]("+strconv.FormatInt(gid, 10)+")\n来自群成员:["+username+"]("+strconv.FormatInt(uid, 10)+")\n以下是留言内容"))
			} else {
				msg = append(msg, message.CustomNode(botname, botid, "来自私聊:["+username+"]("+strconv.FormatInt(uid, 10)+")\n以下是留言内容:"))
			}
			msg = append(msg, message.CustomNode(username, uid, rawmsg))
			ctx.SendPrivateForwardMessage(su, msg)
		})
}

func cpuPercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return -1
	}
	return math.Round(percent[0])
}

func memPercent() float64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return -1
	}
	return math.Round(memInfo.UsedPercent)
}

func diskPercent() string {
	parts, err := disk.Partitions(true)
	if err != nil {
		return err.Error()
	}
	msg := ""
	for _, p := range parts {
		diskInfo, err := disk.Usage(p.Mountpoint)
		if err != nil {
			msg += "\n  - " + err.Error()
			continue
		}
		pc := uint(math.Round(diskInfo.UsedPercent))
		if pc > 0 {
			msg += fmt.Sprintf("\n  - %s(%dM) %d%%", p.Mountpoint, diskInfo.Total/1024/1024, pc)
		}
	}
	return msg
}
