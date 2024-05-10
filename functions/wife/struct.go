package wife

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	fcext "github.com/FloatTech/floatbox/ctxext"
	sql "github.com/FloatTech/sqlite"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var GlobalTimeManager = rate.NewManager[int64](time.Hour*12, 6)
var LeaveTimeManager = rate.NewManager[int64](time.Hour*12, 4)

func init() {
	timer := time.NewTimer(time.Until(getNextExecutionTime()))
	go func() {
		for range timer.C {
			ResetToInitalizeMode()
			timer.Reset(time.Until(getNextExecutionTime()))
		}
	}()
}

// GetUserListAndChooseOne choose people.
func GetUserListAndChooseOne(ctx *zero.Ctx) int64 {
	getThisGroupList := ctx.GetGroupMemberList(ctx.Event.GroupID).Array()
	getThisUserID := getThisGroupList[rand.Intn(len(getThisGroupList))].Get("user_id").Int()
	return getThisUserID
}

// GlobalCDModelCost cd timeManager
func GlobalCDModelCost(ctx *zero.Ctx) bool {
	// 12h 6times.
	UserKeyTag := ctx.Event.UserID + ctx.Event.GroupID
	return GlobalTimeManager.Load(UserKeyTag).Acquire()
}

// GlobalCDModelCostLeastReply Get the existed Token.
func GlobalCDModelCostLeastReply(ctx *zero.Ctx) int {
	UserKeyTag := ctx.Event.UserID + ctx.Event.GroupID
	return int(GlobalTimeManager.Load(UserKeyTag).Tokens())
}

// LeaveCDModelCost cd timeManager
func LeaveCDModelCost(ctx *zero.Ctx) bool {
	// 12h 6times.
	UserKeyTag := ctx.Event.UserID + ctx.Event.GroupID
	return LeaveTimeManager.Load(UserKeyTag).Acquire()
}

// LeaveCDModelCostLeastReply Get the existed Token.
func LeaveCDModelCostLeastReply(ctx *zero.Ctx) int {
	UserKeyTag := ctx.Event.UserID + ctx.Event.GroupID
	return int(LeaveTimeManager.Load(UserKeyTag).Tokens())
}

// CheckTheUserIsTargetOrUser check the status.
func CheckTheUserIsTargetOrUser(db *sql.Sqlite, ctx *zero.Ctx, userID int64) (statuscode int64, targetID int64) {
	// -1 --> not found | 0 --> Target | 1 --> User
	marryLocker.Lock()
	defer marryLocker.Unlock()
	groupID := ctx.Event.GroupID
	var globalDataStructList GlobalDataStruct
	err := db.Find("grouplist_"+strconv.FormatInt(groupID, 10), &globalDataStructList, "where userid is '"+strconv.FormatInt(userID, 10)+"'")
	if err != nil {
		err = db.Find("grouplist_"+strconv.FormatInt(groupID, 10), &globalDataStructList, "where targetid is '"+strconv.FormatInt(userID, 10)+"'")
		if err != nil {
			return -1, -1
		}
		return 0, globalDataStructList.UserID
	}
	if globalDataStructList.TargetID == globalDataStructList.UserID {
		return 10, globalDataStructList.UserID
	}
	return 1, globalDataStructList.TargetID
}

// CheckTheUserIsInBlackListOrGroupList Check this user is in list?
func CheckTheUserIsInBlackListOrGroupList(userID int64, targetID int64, groupID int64) bool {
	/* -1 --> both is null
	1 --> the user random the person that he don't want (Other is in his blocklist) | or in blocklist(others)
	*/
	// first check the blocklist
	if !CheckTheBlackListIsExistedToThisPerson(marryList, userID, targetID) || !CheckTheBlackListIsExistedToThisPerson(marryList, targetID, userID) {
		return true
	}
	// check the target is disabled this group
	if !CheckDisabledListIsExistedInThisGroup(marryList, userID, groupID) {
		return true
	}
	return false
}

// GetSomeRanDomChoiceProps get some props chances to make it random.
func GetSomeRanDomChoiceProps(ctx *zero.Ctx) int64 {
	// get Random numbers.
	randNum := rand.Intn(90) + fcext.RandSenderPerDayN(ctx.Event.UserID, 30)
	if randNum > 110 {
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

// ReplyMeantMode format the reply and clear.
func ReplyMeantMode(header string, referTarget int64, statusCodeToPerson int64, ctx *zero.Ctx) message.MessageID {
	msg := header
	var replyTarget string
	if statusCodeToPerson == 1 {
		replyTarget = "老婆"
	} else {
		replyTarget = "老公"
	}
	aheader := msg + "\n今天你的群" + replyTarget + "是\n"
	formatAvatar := GenerateUserImageLink(referTarget)
	formatReply := "[ " + ctx.CardOrNickName(referTarget) + " ] " + "( " + strconv.FormatInt(referTarget, 10) + " ) 哦w～"
	return ctx.SendChain(message.Text(aheader), message.Image(formatAvatar), message.Text(formatReply))
}

// GenerateMD5 Generate MD5
func GenerateMD5(userID int64, targetID int64, groupID int64) string {
	input := strconv.FormatInt(userID+targetID+groupID, 10)
	hash := md5.New()
	hash.Write([]byte(input))
	hashValue := hash.Sum(nil)
	hashString := hex.EncodeToString(hashValue)
	return hashString
}

// CheckTheUserStatusAndDoRepeat If ture, means it no others (Only Refer to current user.)
func CheckTheUserStatusAndDoRepeat(ctx *zero.Ctx) bool {
	getStatusCode, getOtherUserData := CheckTheUserIsTargetOrUser(marryList, ctx, ctx.Event.UserID) // 判断这个user是否已经和别人在一起了，同时判断Type3
	switch {
	case getStatusCode == 0:
		// case target mode (0)
		ctx.Send(ReplyMeantMode("貌似你今天已经有了哦～", getOtherUserData, 0, ctx))
		return false
	case getStatusCode == 1:
		ctx.Send(ReplyMeantMode("貌似你今天已经有了哦～", getOtherUserData, 1, ctx))
		// case user mode (1)
		return false
	case getStatusCode == 10:
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("啾啾～今天的对象是你自己哦w"))
		return false
	}
	return true
}

// CheckTheTargetUserStatusAndDoRepeat Check the target status and do repeats.
func CheckTheTargetUserStatusAndDoRepeat(ctx *zero.Ctx, ChooseAPerson int64) bool {
	getTargetStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, ChooseAPerson) // 判断这个target是否已经和别人在一起了，同时判断Type3
	switch {
	case getTargetStatusCode == 1 || getTargetStatusCode == 0:
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("对方已经有人了哦w～算是运气不好的一次呢,Lucy多给一次机会呢w"))
		return false
	case getTargetStatusCode == 10:
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("啾啾～今天的对方是单身贵族哦（笑~ Lucy再给你一次机会哦w"))
		return false
	}
	return true
}

// ResuitTheReferUserAndMakeIt Result For Married || be married,
func ResuitTheReferUserAndMakeIt(ctx *zero.Ctx, dict map[string][]string, EventUser int64, TargetUser int64) {
	// get failed possibility.
	props := rand.Intn(100)
	if props < 50 {
		// failed,lost chance.
		getFailedMsg := dict["failed"][rand.Intn(len(dict["failed"]))]
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getFailedMsg))
		return
	}
	returnNumber := GetSomeRanDomChoiceProps(ctx)
	switch {
	case returnNumber == 1:
		GlobalCDModelCost(ctx)
		getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
		// normal mode. nothing happened.
		ctx.Send(ReplyMeantMode(getSuccessMsg, TargetUser, 1, ctx))
		generatePairKey := GenerateMD5(EventUser, TargetUser, ctx.Event.GroupID)
		err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, EventUser, TargetUser, 1, generatePairKey)
		if err != nil {
			fmt.Print(err)
			return
		}
	case returnNumber == 2:
		GlobalCDModelCost(ctx)
		ctx.Send(ReplyMeantMode("貌似很奇怪哦～因为某种奇怪的原因～1变成了0,0变成了1", TargetUser, 0, ctx))
		generatePairKey := GenerateMD5(TargetUser, EventUser, ctx.Event.GroupID)
		err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, TargetUser, EventUser, 2, generatePairKey)
		if err != nil {
			panic(err)
		}
	// reverse Target Mode
	case returnNumber == 3:
		GlobalCDModelCost(ctx)
		// drop target pls.
		// status 3 turns to be case 1 ,for it cannot check it again. (With 2 same person, so it can be panic.)
		getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
		// normal mode. nothing happened.
		ctx.Send(ReplyMeantMode(getSuccessMsg, TargetUser, 1, ctx))
		generatePairKey := GenerateMD5(EventUser, TargetUser, ctx.Event.GroupID)
		err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, EventUser, TargetUser, 1, generatePairKey)
		if err != nil {
			fmt.Print(err)
			return
		}
	// you became your own target
	case returnNumber == 6:
		GlobalCDModelCost(ctx)
		// now no wife mode.
		getHideMsg := dict["hide_mode"][rand.Intn(len(dict["hide_mode"]))]
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getHideMsg, "\n貌似没有任何反馈～"))
		generatePairKey := GenerateMD5(EventUser, TargetUser, ctx.Event.GroupID)
		err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, EventUser, TargetUser, 6, generatePairKey)
		if err != nil {
			panic(err)
		}
	}
}

// CheckThePairKey Check this pair key
func CheckThePairKey(db *sql.Sqlite, uid int64, groupID int64) string {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var globalDataStructList GlobalDataStruct
	err := db.Find("grouplist_"+strconv.FormatInt(groupID, 10), &globalDataStructList, "where userid is '"+strconv.FormatInt(uid, 10)+"'")
	if err != nil {
		err = db.Find("grouplist_"+strconv.FormatInt(groupID, 10), &globalDataStructList, "where targetid is '"+strconv.FormatInt(uid, 10)+"'")
		if err != nil {
			return ""
		}
		return globalDataStructList.PairKey
	}
	return globalDataStructList.PairKey
}

// GenerateUserImageLink Generate Format Link.
func GenerateUserImageLink(uid int64) string {
	index := message.UnescapeCQCodeText("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(uid, 10) + "&s=640")
	return index
}

// ResetToInitalizeMode delete marrylist (pairkey | grouplist)
func ResetToInitalizeMode() {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	getFullList, err := marryList.ListTables()
	if err != nil {
		panic(err)
	}
	// find all the list named: grouplist | pairkey
	getFilteredList := FindStrings(getFullList, "grouplist")
	getPairKeyFilteredList := FindStrings(getFullList, "pairkey")
	getFullFilteredList := append(getFilteredList, getPairKeyFilteredList...)
	getLength := len(getFullFilteredList)
	if getLength == 0 {
		return
	}
	for i := 0; i < getLength; i++ {
		err := marryList.Drop(getFullFilteredList[i])
		if err != nil {
			panic(err)
		}
	}
}

// CheckTheOrderListAndBackDetailed Check this and reply some details
func CheckTheOrderListAndBackDetailed(userID int64, groupID int64) (chance int, target int64, time string) {
	var orderListStructFinder OrderListStruct
	err := marryList.Find("orderlist_"+strconv.FormatInt(groupID, 10), &orderListStructFinder, "where order is '"+strconv.FormatInt(userID, 10)+"'")
	if err != nil {
		return 0, 0, ""
	}
	getTheTarget := orderListStructFinder.TargerPerson
	getTheTimer := orderListStructFinder.Time
	getRandomMoreChance := rand.Intn(30)
	return getRandomMoreChance, getTheTarget, getTheTimer
}

// FindStrings find the strings in list.
func FindStrings(list []string, searchString string) []string {
	result := make([]string, 0)
	for _, item := range list {
		if strings.Contains(item, searchString) {
			result = append(result, item)
		}
	}
	return result
}

// getNextEcecutionTime to this 23:00
func getNextExecutionTime() time.Time {
	now := time.Now()
	nextExecutionTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
	if nextExecutionTime.Before(now) {
		nextExecutionTime = nextExecutionTime.Add(24 * time.Hour)
	}
	return nextExecutionTime
}
