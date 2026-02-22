package core

import (
	"slices"
	"strconv"

	"github.com/Hafuunano/Protocol-ConvertTool/protocol"
	"github.com/Hafuunano/Protocol-ConvertTool/protocol/zerobot"
	zero "github.com/wdvxdr1123/ZeroBot"
)

// dispatch builds a single handler that runs all plugins in sequence. Stops after a handler that called BlockNext().
func dispatch(plugins []protocol.Handler) protocol.Handler {
	return func(ctx protocol.Context) {
		for _, h := range plugins {
			h(ctx)
			if ctx.ShouldBlockNext() {
				break
			}
		}
	}
}

// makeContext builds protocol.Context from zero.Ctx with IsSuperAdmin, IsAdmin, and OnlyToMe wired.
// onlyToMe is true when handling zero.OnMessage(zero.OnlyToMe), false for zero.OnMessage().
func makeContext(ctx *zero.Ctx, onlyToMe bool) *zerobot.Context {
	pc := zerobot.NewContext(ctx)
	pc.IsSuperAdminFunc = func(userID string) bool {
		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return false
		}
		return slices.Contains(zero.BotConfig.SuperUsers, uid)
	}
	pc.IsAdminFunc = func() bool { return zero.AdminPermission(ctx) }
	pc.OnlyToMe = onlyToMe
	return pc
}

// Init registers message handlers: OnMessage (all messages) and OnMessageReply (reply/@ bot only).
func Init() {
	// OnMessage: all messages -> default chain (HookMessage).
	zero.OnMessage().Handle(func(ctx *zero.Ctx) {
		dispatch(protocol.Chain())(makeContext(ctx, false))
	})
	// OnMessageReply: only reply to bot or @ bot -> HookMessageReply chain.
	zero.OnMessage(zero.OnlyToMe).Handle(func(ctx *zero.Ctx) {
		chain := protocol.ChainOn(protocol.HookMessageReply)
		if len(chain) > 0 {
			dispatch(chain)(makeContext(ctx, true))
		}
	})
}
