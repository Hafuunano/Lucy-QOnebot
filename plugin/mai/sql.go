package mai

import (
	sql "github.com/FloatTech/sqlite"
	"strconv"
	"sync"
	"time"
)

type DataHostSQL struct {
	QQ    int64  `db:"user_qq"` // qq nums
	Plate string `db:"plate"`   // plate
}

var (
	maiDatabase = &sql.Sqlite{}
	maiLocker   = sync.Mutex{}
)

func init() {
	maiDatabase.DBPath = engine.DataFolder() + "maisql.db"
	err := maiDatabase.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	_ = InitDataBase()
}

func FormatUserDataBase(qq int64, plate string) *DataHostSQL {
	return &DataHostSQL{QQ: qq, Plate: plate}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Create("userinfo", &DataHostSQL{})
}

func GetUserInfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Plate
}

func (info *DataHostSQL) BindUserDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userinfo", info)
}
