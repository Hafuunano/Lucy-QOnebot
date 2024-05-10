package ticket

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"time"
)

var limit = rate.NewManager[int64](time.Minute*3, 28)

func GetTiredToken(ctx *zero.Ctx) float64 {
	return limit.Load(ctx.Event.UserID).Tokens()
}

func GetCostTiredToken(ctx *zero.Ctx) bool {
	return limit.Load(ctx.Event.UserID).AcquireN(3)
}
