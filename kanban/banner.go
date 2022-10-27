package kanban // Package kanban package kanban 初始化

import (
	"strings"
)

var (
	info = [...]string{
		"* OneBot + Zerobot + Nonebot2 with ZeroBot-Plugin Project.",
		"* Hosted On Tencent LightCloudServer in Nanjing.",
		"* Made By MoeMagicMango and FloatTech Project With ❤",
		"* Project: https://github.com/FloatTech/ZeroBot-Plugin",
		"* 说明书: https://manual-lucy.himoyo.cn",
		"* Copyright © 2021-2022 FloatTech. All Rights Reserved.",
	}
	// Banner ...
	Banner = strings.Join(info[:], "\n")
)
