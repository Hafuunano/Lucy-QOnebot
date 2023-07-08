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

type DataHostPic struct {
	Pic string `db:"base64"`
	QQ  int64  `db:"user_qq"` // qq nums
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

func FormatUserDataBasePic(qq int64, pic string) *DataHostPic {
	return &DataHostPic{Pic: pic, QQ: qq}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	err := maiDatabase.Create("userpic", &DataHostPic{})
	err = maiDatabase.Create("userinfo", &DataHostSQL{})
	return err
}

func GetUserInfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Plate
}

func GetUserInfoFromDatabaseForPic(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostPic
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userpic", &infosql, "where user_qq is "+userIDStr)
	return infosql.Pic
}

func (info *DataHostSQL) BindUserDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userinfo", info)
}

func (info *DataHostPic) BindUserDataBaseForPic() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userpic", info)
}

func DeleteUserDataPic(userID int64) error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	userIDStr := strconv.FormatInt(userID, 10)
	return maiDatabase.Del("userpic", "where user_qq is "+userIDStr)
}
