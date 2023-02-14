package catmorecat

import (
	"sync"

	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

var (
	engine = control.Register("catmorecat", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "Cat More Cat!",
	})
	catdata = &catdb{
		db: &sql.Sqlite{},
	}
)

type catdb struct {
	db *sql.Sqlite
	sync.RWMutex
}

func init() {
	engine.OnFullMatch("领养猫猫").SetBlock(true).Handle(func(ctx *zero.Ctx) { //	领养猫猫是不需要钱的(

	})
}
