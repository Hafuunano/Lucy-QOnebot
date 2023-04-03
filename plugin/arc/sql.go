package arc

import (
	sql "github.com/FloatTech/sqlite"
	zero "github.com/wdvxdr1123/ZeroBot"
	"strconv"
	"sync"
	"time"
)

// Arcinfosql use nonebot-arcaeabot's config.
type Arcinfosql struct {
	QQ    int64  `db:"user_qq"`   // qq nums
	Arcid string `db:"arcaea_id"` // arcid nums
	name  string `db:"arcname"`   // arcname
}

var (
	arcAcc    = &sql.Sqlite{}
	arcLocker = sync.Mutex{}
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
func FormatInfo(qqnum int64, arcid string, getArcUserName string) *Arcinfosql {
	return &Arcinfosql{Arcid: arcid, QQ: qqnum, name: getArcUserName}
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

// GetUserArcaeaInfo check user info.
func GetUserArcaeaInfo(db *sql.Sqlite, ctx *zero.Ctx) (account string, err error) {
	arcLocker.Lock()
	defer arcLocker.Unlock()
	uidStr := strconv.FormatInt(ctx.Event.UserID, 10)
	var infosql Arcinfosql
	err = db.Find("userinfo", &infosql, "where user_qq is "+uidStr)
	if err != nil {
		return "", nil
	}
	return infosql.Arcid, err
}
