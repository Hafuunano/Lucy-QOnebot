package arc

import (
	sql "github.com/FloatTech/sqlite"
	"sync"
	"time"
)

// Arcinfosql use nonebot-arcaeabot's config.
type Arcinfosql struct {
	QQ    int64 `db:"user_qq"`   // qq nums
	Arcid int64 `db:"arcaea_id"` // arcid nums
}

var (
	arcAcc    = &sql.Sqlite{}
	arcLocker = sync.RWMutex{}
)

func init() {
	arcAcc.DBPath = engine.DataFolder() + "arcacc.db"
	err := arcAcc.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	_ = InitalizeSqlite(arcAcc)
}

// FormatInfo FormatUserInfo and prepare to send it to sql.
func FormatInfo(qqnum int64, arcid int64) *Arcinfosql {
	return &Arcinfosql{Arcid: arcid, QQ: qqnum}
}

// BindUserArcaeaInfo Bind user's acc.
func (info *Arcinfosql) BindUserArcaeaInfo(db *sql.Sqlite) error {
	arcLocker.Lock()
	defer arcLocker.Unlock()
	return db.Insert("userinfo", info)
}

// InitalizeSqlite Initalize sqlite.
func InitalizeSqlite(db *sql.Sqlite) error {
	arcLocker.Lock()
	defer arcLocker.Unlock()
	return db.Create("userinfo", &Arcinfosql{})
}
