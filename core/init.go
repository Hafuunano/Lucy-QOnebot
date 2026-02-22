package core

import (
	"slices"
	"strconv"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/zerobot"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// dispatch builds a single handler that runs all plugins in sequence.
func dispatch(plugins []protocol.Handler) protocol.Handler {
	return func(ctx protocol.Context) {
		for _, h := range plugins {
			h(ctx)
		}
	}
}

// Init registers the global message handler: converts zero.Ctx to protocol.Context and runs the plugin chain.
func Init() {
	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		pc := zerobot.NewContext(ctx)
		// Set IsSuperAdmin check from ZeroBot config (SuperUsers).
		pc.IsSuperAdminFunc = func(userID string) bool {
			uid, err := strconv.ParseInt(userID, 10, 64)
			if err != nil {
				return false
			}
			return slices.Contains(zero.BotConfig.SuperUsers, uid)
		}
		dispatch(protocol.Chain())(pc)
	})
}
