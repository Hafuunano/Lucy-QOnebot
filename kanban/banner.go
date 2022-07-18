package kanban

import (
	"fmt"
	"strings"

	"github.com/fumiama/go-registry"
)

var (
	info = [...]string{
		"* OneBot + Zerobot + Nonebot2 with ZeroBot-Plugin Project.",
		"* 基于 ZeroBot-Plugin 1.50 Beta4 修改的 Lucy",
		"* Made By MoeMagicMango With ❤",
		"* Project: https://github.com/FloatTech/ZeroBot-Plugin",
		"* 说明书: https://manual-lucy.himoyo.cn",
		"* Copyright © 2021-2022 FloatTech. All Rights Reserved.",
	}
	// Banner ...
	Banner = strings.Join(info[:], "\n")
	reg    = registry.NewRegReader("reilia.fumiama.top:32664", "fumiama")
)

// PrintBanner ...
func PrintBanner() {
	fmt.Print(
		"\n======================[ZeroBot-Plugin]======================",
		"\n", Banner, "\n",
		"----------------------[ZeroBot-公告栏]----------------------",
		"\n", Kanban(), "\n",
		"============================================================\n\n",
	)
}

// Kanban ...
func Kanban() string {
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
