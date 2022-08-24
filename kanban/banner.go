package kanban

import (
	"strings"

	"github.com/fumiama/go-registry"
)

var (
	info = [...]string{
		"* OneBot + Zerobot + Nonebot2 with ZeroBot-Plugin Project.",
		"* Made By MoeMagicMango With ❤",
		"* Project: https://github.com/FloatTech/ZeroBot-Plugin",
		"* 说明书: https://manual-lucy.himoyo.cn",
		"* Copyright © 2021-2022 FloatTech. All Rights Reserved.",
	}
	// Banner ...
	Banner = strings.Join(info[:], "\n")
	reg    = registry.NewRegReader("reilia.fumiama.top:32664", "fumiama")
)

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
