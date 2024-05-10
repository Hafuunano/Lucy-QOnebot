// Package group For group commit here.
package group

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	engine = control.Register("group", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Group for mainly use reaction.",
		PrivateDataFolder: "group",
	})
)
