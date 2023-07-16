// Package kanban package kanban 初始化
package kanban

import (
	"strings"
)

var (
	info = [...]string{
		"* OneBot + Zerobot + Nonebot2 with Zerobot-Plugin Project.",
		"* Made By MoeMagicMango and FloatTech Project With ❤",
		"* Project: https://github.com/FloatTech/ZeroBot-Plugin",
		"* 说明书: https://lucy.impart.icu",
		"* 认领: https://moe.himoyo.cn/archives/110/",
		"* Copyright © 2021-2023 FloatTech. All Rights Reserved.",
	}
	// Banner ...
	Banner = strings.Join(info[:], "\n")
)
