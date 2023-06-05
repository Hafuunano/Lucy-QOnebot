package wife

import (
	fcext "github.com/FloatTech/floatbox/ctxext"
	sql "github.com/FloatTech/sqlite"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"math/rand"
	"strconv"
	"time"
)

var GlobalTimeManager = rate.NewManager[int64](time.Hour*12, 6)

// GetUserListAndChooseOne choose people.
func GetUserListAndChooseOne(ctx *zero.Ctx) int64 {
	getThisGroupList := ctx.GetGroupMemberList(ctx.Event.GroupID).Array()
	getThisUserID := getThisGroupList[rand.Intn(len(getThisGroupList))].Get("user_id").Int()
	return getThisUserID
}

// GlobalCDModel cd timeManager
func GlobalCDModel(ctx *zero.Ctx) bool {
	// 12h 6times.
	UserKeyTag := ctx.Event.UserID + ctx.Event.GroupID
	if !GlobalTimeManager.Load(UserKeyTag).Acquire() {
		return false
	}
	return true
}

// CheckTheUserIsTargetOrUser check the status.
func CheckTheUserIsTargetOrUser(db *sql.Sqlite, ctx *zero.Ctx) int64 {
	// -1 --> not found | 0 --> Target | 1 --> User
	marryLocker.Lock()
	defer marryLocker.Unlock()
	userID := ctx.Event.UserID
	groupID := ctx.Event.GroupID
	err := db.Find("grouplist_"+strconv.FormatInt(groupID, 10), &GlobalDataStruct{}, "Where userid is "+strconv.FormatInt(userID, 10))
	if err != nil {
		err := db.Find("grouplist_"+strconv.FormatInt(groupID, 10), &GlobalDataStruct{}, "Where targetid is "+strconv.FormatInt(userID, 10))
		if err != nil {
			return -1
		}
		return 0
	}
	return 1
}

// CheckTheUserIsInBlackListOrGroupList Check this user is in list?
func CheckTheUserIsInBlackListOrGroupList(userID int64, targetID int64, groupID int64) bool {
	/* -1 --> both is null
	1 --> the user random the person that he don't want (Other is in his blocklist) | or in blocklist(others)
	*/
	// first check the blocklist
	if CheckTheBlackListIsExistedToThisPerson(marryList, userID, targetID) || CheckTheBlackListIsExistedToThisPerson(marryList, targetID, userID) {
		return true
	}
	if CheckDisabledListIsExistedInThisGroup(marryList, targetID, groupID) {
		return true
	}
	return false
}

func GetSomeRanDomChoiceProps(ctx *zero.Ctx) int64 {
	// get Random numbers.
	randNum := fcext.RandSenderPerDayN(ctx.Event.UserID, 100)
	if randNum > 89 {
		getOtherRand := rand.Intn(9)
		switch {
		case getOtherRand < 3:
			return 2
		case getOtherRand > 3 && getOtherRand < 6:
			return 3
		case getOtherRand > 6:
			return 6
		}
	}
	return 1
}
