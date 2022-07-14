// Package aifalse 暂时只有服务器监控
package tools

import (
	"fmt"
	"math"

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

	engine.OnRegex(`^【.*】.*`, zero.SuperUserPermission).SetBlock(false).
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
