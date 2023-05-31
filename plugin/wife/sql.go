package wife

import (
	sql "github.com/FloatTech/sqlite"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"sync"
	"time"
)

// JointList CheckThisGroup first | example: JointList+groupid | warning: be careful this need clean when the next day go. | EveryUser use keys when paired.
type JointList struct {
	Key int64 `db:"key"`
}

// UserGroup table for every Group. | Example: groupid
type UserGroup struct {
	HandleID int64 `db:"handleid"` // user,set 1, sometimes it can be 0. ww.
	Time     int64 `db:"time"`     // use timestamp (HH:MM:SS)
	TargetID int64 `db:"targatid"` // should existed.,
	Key      int64 `db:"key"`      // pair use key to choose.
}

// BlackListGenerator Give User Permissions to let them block someone. | BlackList+QQID
type BlackListGenerator struct {
	QQ        int64    `db:"user_qq"`
	blacklist []string `db:"blacklist"`
}

// UserPermissionStatus Give user permission if user didn't want to be chosen by bot. | PermissionList+QQID
type UserPermissionStatus struct {
	QQ           int64    `db:"user_qq"`
	CantBeChosen []string `db:"cantbechosen"` // input groupid.
}

var (
	userlistLocker   = sync.Mutex{}
	userlistDatabase = &sql.Sqlite{}
)

func init() {
	userlistDatabase.DBPath = engine.DataFolder() + "list.db"
	err := userlistDatabase.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
	err = InitBlackListAndUserPermission(userlistDatabase)
	if err != nil {
		panic(err)
	}

}

// UserPermissionStatusRemoveDisabled Disabled Permission status.
func UserPermissionStatusRemoveDisabled(EventUser int64, groupID int64, sql *sql.Sqlite, ctx *zero.Ctx) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	var permissionUser UserPermissionStatus
	err := sql.Find("PermissionList"+strconv.FormatInt(EventUser, 10), &permissionUser, "where user_qq is "+strconv.FormatInt(EventUser, 10))
	if err != nil {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("NO NEED TO REMOVE"))
		return err
	}
	getList := permissionUser.CantBeChosen
	index := -1
	for i, str := range getList {
		if str == strconv.FormatInt(groupID, 10) {
			index = i
			break
		}
	}
	if index != -1 {
		getList = append(getList[:index], getList[index+1:]...)
	}
	prepareToThrow := &UserPermissionStatus{QQ: EventUser, CantBeChosen: getList}
	err = sql.Insert("PermissionList"+strconv.FormatInt(EventUser, 10), &prepareToThrow)
	if err != nil {
		panic(err)
	}
	return nil
}

// UserPermissionStatusInsertEnabled insert group to make it disabled.
func UserPermissionStatusInsertEnabled(EventUser int64, groupID int64, sql *sql.Sqlite) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	var permissionUser UserPermissionStatus
	err := sql.Find("PermissionList"+strconv.FormatInt(EventUser, 10), &permissionUser, "where user_qq is "+strconv.FormatInt(EventUser, 10))
	if err != nil {
		// if not existed.
		err := sql.Create("PermissionList"+strconv.FormatInt(EventUser, 10), &UserPermissionStatus{})
		if err != nil {
			panic(err)
		}
		disabledGroup := []string{strconv.FormatInt(groupID, 10)}
		prepareToThrow := &UserPermissionStatus{QQ: EventUser, CantBeChosen: disabledGroup}
		err = sql.Insert("PermissionList"+strconv.FormatInt(EventUser, 10), &prepareToThrow)
		if err != nil {
			panic(err)
		}
		return nil
	}
	getList := permissionUser.CantBeChosen
	getTargetIsExisted := CheckTheListIsExistedTheTarget(getList, strconv.FormatInt(groupID, 10))
	if !getTargetIsExisted {
		addList := append(getList, strconv.FormatInt(groupID, 10))
		prepareToThrow := &UserPermissionStatus{QQ: EventUser, CantBeChosen: addList}
		err = sql.Insert("PermissionList"+strconv.FormatInt(EventUser, 10), &prepareToThrow)
		if err != nil {
			panic(err)
		}
		return nil
	}
	// existed.
	return err
}

// InsertUserBlackList Insert BlackList
func InsertUserBlackList(EventUserId int64, blockID int64, sql *sql.Sqlite) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	// first get []string from user.
	var blackList BlackListGenerator
	err := sql.Find("BlackList"+strconv.FormatInt(EventUserId, 10), &blackList, "where user_qq is "+strconv.FormatInt(EventUserId, 10))
	if err != nil {
		// if not existed.
		err := sql.Create("BlackList"+strconv.FormatInt(EventUserId, 10), &BlackListGenerator{})
		if err != nil {
			panic(err)
		}
		blackList := []string{strconv.FormatInt(blockID, 10)}
		prePareThrowList := &BlackListGenerator{QQ: EventUserId, blacklist: blackList}
		err = sql.Insert("BlackList"+strconv.FormatInt(EventUserId, 10), &prePareThrowList)
		if err != nil {
			panic(err)
		}
		return nil
	}
	getList := blackList.blacklist
	getTargetIsExisted := CheckTheListIsExistedTheTarget(getList, strconv.FormatInt(blockID, 10))
	if !getTargetIsExisted {
		addList := append(getList, strconv.FormatInt(blockID, 10))
		prePareThrowList := &BlackListGenerator{QQ: EventUserId, blacklist: addList}
		err = sql.Insert("BlackList"+strconv.FormatInt(EventUserId, 10), &prePareThrowList)
		if err != nil {
			panic(err)
		}
		return nil
	}
	// then stop this case.(existed.)
	return err
}

// RemoveUserBlackList Remove BlackList
func RemoveUserBlackList(EventUserId int64, blockID int64, sql *sql.Sqlite, ctx *zero.Ctx) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	// first get []string from user.
	var blackList BlackListGenerator
	err := sql.Find("BlackList"+strconv.FormatInt(EventUserId, 10), &blackList, "where user_qq is "+strconv.FormatInt(EventUserId, 10))
	if err != nil {
		// if not existed, no need to remove
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("NO NEED TO REMOVE."))
		return err
	}
	getList := blackList.blacklist
	index := -1
	for i, str := range getList {
		if str == strconv.FormatInt(blockID, 10) {
			index = i
			break
		}
	}
	if index != -1 {
		getList = append(getList[:index], getList[index+1:]...)
	}
	prePareThrowList := &BlackListGenerator{QQ: EventUserId, blacklist: getList}
	err = sql.Insert("BlackList"+strconv.FormatInt(EventUserId, 10), &prePareThrowList)
	if err != nil {
		panic(err)
	}
	return nil
}

// InsertUserMarriedMember tips: you need to make sure that all checked.
func InsertUserMarriedMember(handleID int64, targetID int64, ctx *zero.Ctx, sql *sql.Sqlite, key int64) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	getThisGroupID := ctx.Event.GroupID
	prepareThrownList := &UserGroup{HandleID: handleID, TargetID: targetID, Time: time.Now().Unix(), Key: key}
	err := sql.Insert(strconv.FormatInt(getThisGroupID, 10), prepareThrownList)
	// ensure the data can be thrown.
	if err != nil {
		err := sql.Create(strconv.FormatInt(getThisGroupID, 10), &UserGroup{})
		err = sql.Insert(strconv.FormatInt(getThisGroupID, 10), prepareThrownList)
		if err != nil {
			panic(err)
		}
	}
	return err
}

// InitBlackListAndUserPermission Initalize sqlite data.
func InitBlackListAndUserPermission(sqlite *sql.Sqlite) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	err := sqlite.Create("blacklist", &BlackListGenerator{})
	err = sqlite.Create("permissionList", &UserPermissionStatus{})
	return err
}

// CheckTheListIsExistedTheTarget feature.
func CheckTheListIsExistedTheTarget(list []string, target string) bool {
	index := -1
	for i, str := range list {
		if str == target {
			index = i
			break
		}
	}
	if index != -1 {
		return true
	}
	return false
}

// JointListIfFind find it. if not existed then created.
func JointListIfFind(sql *sql.Sqlite, eventKey int64, groupID int64) bool {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	var JointListHandle JointList
	err := sql.Find("JointList"+strconv.FormatInt(groupID, 10), &JointListHandle, "where key is "+strconv.FormatInt(eventKey, 10))
	if err != nil {
		// not existed? or cannot find.
		_ = sql.Create("JointList"+strconv.FormatInt(groupID, 10), &JointListHandle)
		return false
	}
	return true
}

// JointListAdd add it
func JointListAdd(sql *sql.Sqlite, eventKey int64, groupID int64) error {
	getStatus := JointListIfFind(sql, eventKey, groupID)
	if !getStatus {
		// add.
		err := sql.Insert("JointList"+strconv.FormatInt(groupID, 10), &JointList{Key: eventKey})
		if err != nil {
			return err
		}
	}
	return nil
}

// JointListRemove Remove it
func JointListRemove(sql *sql.Sqlite, eventKey int64, groupID int64) error {
	userlistLocker.Lock()
	defer userlistLocker.Unlock()
	getStatus := JointListIfFind(sql, eventKey, groupID)
	if getStatus {
		// add.
		err := sql.Del("JointList"+strconv.FormatInt(groupID, 10), "where key is "+strconv.FormatInt(eventKey, 10))
		if err != nil {
			return err
		}
	}
	return nil
}
