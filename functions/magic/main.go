package magic

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var engine = control.Register("magic", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault: false,
})

func init() {
	// let's roll with Lucy's Magic, use datapack mode.
}
