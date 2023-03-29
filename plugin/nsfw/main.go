// Package nsfw for nsfw (
package nsfw

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	engine = control.Register("nsfw", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "Hi NekoPachi!\n",
	})
)
