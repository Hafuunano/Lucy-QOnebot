// Package main, Lucy, are you here?
package main

import (
	"github.com/HafuuNano/Lucy-QOnebot/core"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	core.Init()
	zero.RunAndBlock(&zero.Config{
		NickName:      []string{""},
		CommandPrefix: "/",
		SuperUsers:    []int64{},
		Driver: []zero.Driver{
			driver.NewWebSocketServer(16, "ws://127.0.0.1:6700", ""),
		},
	}, nil)
}
