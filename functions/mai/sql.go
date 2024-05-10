package mai

import (
	"strconv"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
)

type UserIDToQQ struct {
	QQ     int64  `db:"user_qq"` // qq nums
	Userid string `db:"user_id"` // user_id
}

type DataHostSQL struct {
	QQ         int64  `db:"user_qq"` // qq nums
	Plate      string `db:"plate"`   // plate
	Background string `db:"bg"`      // bg
}

type UserSwitcherService struct {
	QQ     int64 `db:"user_qq"`
	IsUsed bool  `db:"isused"` // true == lxns service \ false == Diving Fish.
}

type UserIDToToken struct {
	// refer usertoken ==> qqid
	UserID string `db:"user_id"` // THIS IS QQ
	Token  string `db:"user_token"`
}

type UserAutoSync struct {
	QQ     int64 `db:"user_qq"`
	UserID int64 `db:"user_id"`
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
	InitDataBase()
}

func FormatUserDataBase(qq int64, plate string, bg string) *DataHostSQL {
	return &DataHostSQL{QQ: qq, Plate: plate, Background: bg}
}

func FormatUserSwitcher(qq int64, isSwitcher bool) *UserSwitcherService {
	return &UserSwitcherService{QQ: qq, IsUsed: isSwitcher}
}

func FormatUserIDDatabase(qq int64, userid string) *UserIDToQQ {
	return &UserIDToQQ{QQ: qq, Userid: userid}
}

func FormatUserToken(qq string, token string) *UserIDToToken {
	return &UserIDToToken{Token: token, UserID: qq}
}

func InitDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	maiDatabase.Create("userinfo", &DataHostSQL{})
	maiDatabase.Create("useridinfo", &UserIDToQQ{})
	maiDatabase.Create("userswitcherinfo", &UserSwitcherService{})
	maiDatabase.Create("usertokenid", &UserIDToToken{})
	maiDatabase.Create("userautoupdate", &UserAutoSync{})
	return nil
}

// AutoUpdate : ReMaimai Project

// STEP: RECEIVE USER UPDATER BODY, Query User QQ ==> Check USERTOKEN

func InsertQueryParamsToStruct(autoSyncStruct UserAutoSync) error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	err := maiDatabase.Insert("userautoupdate", &autoSyncStruct)
	return err
}

func QueryUserTokenData(userID int64) UserAutoSync {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var returnData UserAutoSync
	maiDatabase.Find("userautoupdate", &returnData, "WHERE user_id is "+strconv.FormatInt(userID, 10))
	return returnData
}

func QueryUserTokenDataByQQ(QQ int64) UserAutoSync {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var returnData UserAutoSync
	maiDatabase.Find("userautoupdate", &returnData, "WHERE user_qq is "+strconv.FormatInt(QQ, 10))
	return returnData
}

func RemoveUserTokenData(qq int64) error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Del("userautoupdate", "WHERE user_qq is "+strconv.FormatInt(qq, 10))
}

func GetUserSwitcherInfoFromDatabase(userid int64) bool {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var info UserSwitcherService
	userIDStr := strconv.FormatInt(userid, 10)
	err := maiDatabase.Find("userswitcherinfo", &info, "where user_qq is "+userIDStr)
	if err != nil {
		return false
	}
	return info.IsUsed
}

func (info *UserSwitcherService) ChangeUserSwitchInfoFromDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userswitcherinfo", info)
}

// GetUserIDFromDatabase Params: user qq id ==> user maimai id.
func GetUserIDFromDatabase(userQQ int64) UserIDToQQ {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql UserIDToQQ
	userIDStr := strconv.FormatInt(userQQ, 10)
	_ = maiDatabase.Find("useridinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql
}

// GetUserQQFromDatabase Params: user qq id ==> user maimai id.
func GetUserQQFromDatabase(userid int64) UserIDToQQ {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql UserIDToQQ
	userIDStr := strconv.FormatInt(userid, 10)
	_ = maiDatabase.Find("useridinfo", &infosql, "where user_id is "+userIDStr)
	return infosql
}

func RemoveUserIdFromDatabase(qq int64) {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	userIDStr := strconv.FormatInt(qq, 10)
	_ = maiDatabase.Del("useridinfo", "where user_qq is "+userIDStr)
}

func (info *UserIDToQQ) BindUserIDDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("useridinfo", info)
}

// maimai render b50

func GetUserInfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Plate
}

func GetUserDefaultinfoFromDatabase(userID int64) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql DataHostSQL
	userIDStr := strconv.FormatInt(userID, 10)
	_ = maiDatabase.Find("userinfo", &infosql, "where user_qq is "+userIDStr)
	return infosql.Background
}

func (info *DataHostSQL) BindUserDataBase() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("userinfo", info)
}

func GetUserToken(userid string) string {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	var infosql UserIDToToken // Here is User_id ==> userQQ ,?
	maiDatabase.Find("usertokenid", &infosql, "where user_id is '"+userid+"'")
	if infosql.Token == "" {
		return ""
	}
	return infosql.Token
}

func (info *UserIDToToken) BindUserToken() error {
	maiLocker.Lock()
	defer maiLocker.Unlock()
	return maiDatabase.Insert("usertokenid", info)
}
