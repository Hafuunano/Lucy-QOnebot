package daily

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	engine = control.Register("daily", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Daily Refresh == > Score / Fortune. ",
		PrivateDataFolder: "daily",
	})
)
