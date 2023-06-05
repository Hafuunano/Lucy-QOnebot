package wife

import (
	sql "github.com/FloatTech/sqlite"
	"strconv"
	"sync"
	"time"
)

// BlackListStruct (Example: blacklist_1292581422)
type BlackListStruct struct {
	BlackList int64 `db:"blacklist"`
}

// DisabledListStruct (Example: disabled_1292581422)
type DisabledListStruct struct {
	DisabledList int64 `db:"disabledlist"`
}

// OrderListStruct (Example: orderlist_1145141919180)
type OrderListStruct struct {
	OrderPerson  int64  `db:"order"`
	TargerPerson int64  `db:"target"`
	Time         string `db:"time"`
}

// GlobalDataStruct (Example: grouplist_1292581422 )
type GlobalDataStruct struct {
	PairKey  string `db:"pairkey"`
	UserID   int64  `db:"userid"`
	TargetID int64  `db:"targetid"`
	time     string `db:"time"`
}

// PairKeyStruct pairkey is used to check the list.
type PairKeyStruct struct {
	PairKey  string `db:"pairkey"`
	statusID int64  `db:"statusid"`
}

var (
	marryList   = &sql.Sqlite{}
	marryLocker = sync.Mutex{}
)

func init() {
	// init UserDataBase
	marryList.DBPath = engine.DataFolder() + "data.db"
	err := marryList.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

/*
SQL HANDLER
*/

// InitTable init table if not existed.
func InitTable(table string, db *sql.Sqlite, structList interface{}) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	return db.Create(table, &structList)
}

func FormatInsertUserGlobalMarryList(UserID int64, targetID int64, PairKeyRaw string) *GlobalDataStruct {
	return &GlobalDataStruct{PairKey: PairKeyRaw, UserID: UserID, TargetID: targetID, time: strconv.FormatInt(time.Now().Unix(), 10)}
}

func FormatPairKey(PairKeyRaw string, statusID int64) *PairKeyStruct {
	return &PairKeyStruct{PairKey: PairKeyRaw, statusID: statusID}
}

func FormatBlackList(blockID int64) *BlackListStruct {
	return &BlackListStruct{BlackList: blockID}
}

func FormatDisabledList(disabledID int64) *DisabledListStruct {
	return &DisabledListStruct{DisabledList: disabledID}
}

func FormatOrderList(orderPersonal int64, targetPersonal int64, time string) *OrderListStruct {
	return &OrderListStruct{OrderPerson: orderPersonal, TargerPerson: targetPersonal, Time: time}
}

// InsertUserGlobalMarryList no check,use it with caution // Only insert and do nothing.
func InsertUserGlobalMarryList(db *sql.Sqlite, groupID int64, UserID int64, targetID int64, StatusID int64, PairKeyRaw string) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	formatList := FormatInsertUserGlobalMarryList(UserID, targetID, PairKeyRaw)
	err := db.Insert("grouplist_"+strconv.FormatInt(groupID, 10), formatList)
	if err != nil {
		_ = InitTable("grouplist_"+strconv.FormatInt(groupID, 10), db, &GlobalDataStruct{})
		err = db.Insert("grouplist_"+strconv.FormatInt(groupID, 10), formatList)
	}
	// throw key
	err = db.Insert("pairkey_"+strconv.FormatInt(groupID, 10), FormatPairKey(PairKeyRaw, StatusID))
	if err != nil {
		_ = InitTable("pairkey_"+strconv.FormatInt(groupID, 10), db, &PairKeyStruct{})
		_ = db.Insert("pairkey_"+strconv.FormatInt(groupID, 10), FormatPairKey(PairKeyRaw, StatusID))
	}
	return err
}

// RemoveUserGlobalMarryList just remove it,it will persist the key and change it to Type4.
func RemoveUserGlobalMarryList(db *sql.Sqlite, pairKey string, groupID int64) bool {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	// first find the key list.
	var pairKeyNeed PairKeyStruct
	err := db.Find("pairkey_"+strconv.FormatInt(groupID, 10), &pairKeyNeed, "where pairkey is "+pairKey)
	if err != nil {
		// cannnot find, don't need to remove.
		return false
	}
	// if found.
	getThisKey := pairKeyNeed.PairKey
	err = db.Del("pairkey_"+strconv.FormatInt(groupID, 10), "where pairkey is "+getThisKey)
	err = db.Del("grouplist_"+strconv.FormatInt(groupID, 10), "where pairkey is "+getThisKey)
	// store? || persist this key and check the next time.
	err = db.Insert("pairkey_"+strconv.FormatInt(groupID, 10), FormatPairKey(pairKey, 4))
	if err != nil {
		return false
	}
	return true
}

// CustomRemoveUserGlobalMarryList just remove it,it will persist the key (referKeyStatus)
func CustomRemoveUserGlobalMarryList(db *sql.Sqlite, pairKey string, groupID int64, statusID int64) bool {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	// first find the key list.
	var pairKeyNeed PairKeyStruct
	err := db.Find("pairkey_"+strconv.FormatInt(groupID, 10), &pairKeyNeed, "where pairkey is "+pairKey)
	if err != nil {
		// cannnot find, don't need to remove.
		return false
	}
	// if found.
	getThisKey := pairKeyNeed.PairKey
	err = db.Del("pairkey_"+strconv.FormatInt(groupID, 10), "where pairkey is "+getThisKey)
	err = db.Del("grouplist_"+strconv.FormatInt(groupID, 10), "where pairkey is "+getThisKey)
	// store? || persist this key and check the next time.
	err = db.Insert("pairkey_"+strconv.FormatInt(groupID, 10), FormatPairKey(pairKey, statusID))
	if err != nil {
		return false
	}
	return true
}

// CheckThisKeyStatus check this key status.
func CheckThisKeyStatus(db *sql.Sqlite, pairKey string, groupID int64) int64 {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	// first find the key list.
	var pairKeyNeed PairKeyStruct
	err := db.Find("pairkey_"+strconv.FormatInt(groupID, 10), &pairKeyNeed, "where pairkey is "+pairKey)
	if err != nil {
		// cannnot find, don't need to remove.
		return -1
	}
	return pairKeyNeed.statusID
}

// AddBlackList add blacklist
func AddBlackList(db *sql.Sqlite, userID int64, targetID int64) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var blackListNeed BlackListStruct
	err := db.Find("blacklist_"+strconv.FormatInt(userID, 10), &blackListNeed, "where blacklist is "+strconv.FormatInt(targetID, 10))
	if err != nil {
		// add it, not sure then init this and add.
		_ = InitTable("blacklist_"+strconv.FormatInt(userID, 10), db, BlackListStruct{})
		err = db.Insert("blacklist_"+strconv.FormatInt(userID, 10), FormatBlackList(targetID))
		return err
	}
	// find this so don't do it.
	return err
}

// DeleteBlackList delete blacklist.
func DeleteBlackList(db *sql.Sqlite, userID int64, targetID int64) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var blackListNeed BlackListStruct
	err := db.Find("blacklist_"+strconv.FormatInt(userID, 10), &blackListNeed, "where blacklist is "+strconv.FormatInt(targetID, 10))
	if err != nil {
		// not init or didn't find.
		return err
	}
	_ = db.Del("blacklist_"+strconv.FormatInt(userID, 10), "where blacklist is "+strconv.FormatInt(targetID, 10))
	return err
}

// CheckTheBlackListIsExistedToThisPerson check the person is blocked.
func CheckTheBlackListIsExistedToThisPerson(db *sql.Sqlite, userID int64, targetID int64) bool {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var blackListNeed BlackListStruct
	err := db.Find("blacklist_"+strconv.FormatInt(userID, 10), &blackListNeed, "where blacklist is "+strconv.FormatInt(targetID, 10))
	if err != nil {
		// not init or didn't find.
		return false
	}
	return true
}

// AddDisabledList add disabledList
func AddDisabledList(db *sql.Sqlite, userID int64, groupID int64) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var disabledListNeed DisabledListStruct
	err := db.Find("disabled_"+strconv.FormatInt(userID, 10), &disabledListNeed, "where disabledlist is "+strconv.FormatInt(groupID, 10))
	if err != nil {
		// add it, not sure then init this and add.
		_ = InitTable("disabled_"+strconv.FormatInt(userID, 10), db, DisabledListStruct{})
		err = db.Insert("disabled_"+strconv.FormatInt(userID, 10), FormatDisabledList(groupID))
		return err
	}
	// find this so don't do it.
	return err
}

// DeleteDisabledList delete disabledlist.
func DeleteDisabledList(db *sql.Sqlite, userID int64, groupID int64) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var disabledListNeed DisabledListStruct
	err := db.Find("disabled_"+strconv.FormatInt(userID, 10), &disabledListNeed, "where disabledlist is "+strconv.FormatInt(groupID, 10))
	if err != nil {
		// not init or didn't find.
		return err
	}
	_ = db.Del("disabled_"+strconv.FormatInt(userID, 10), "where disabledlist is "+strconv.FormatInt(groupID, 10))
	return err
}

// CheckDisabledListIsExistedInThisGroup Check this group.
func CheckDisabledListIsExistedInThisGroup(db *sql.Sqlite, userID int64, groupID int64) bool {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var disabledListNeed DisabledListStruct
	err := db.Find("disabled_"+strconv.FormatInt(userID, 10), &disabledListNeed, "where disabledlist is "+strconv.FormatInt(groupID, 10))
	if err != nil {
		// not init or didn't find.
		return false
	}
	return true
}

// AddOrderToList add order.
func AddOrderToList(db *sql.Sqlite, userID int64, targetID int64, time string, groupID int64) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var addOrderListNeed OrderListStruct
	err := db.Find("orderlist_"+strconv.FormatInt(groupID, 10), &addOrderListNeed, "where order is "+strconv.FormatInt(userID, 10))
	if err != nil {
		// create and insert.
		_ = InitTable("orderlist_"+strconv.FormatInt(groupID, 10), db, OrderListStruct{})
		_ = db.Insert("orderlist_"+strconv.FormatInt(groupID, 10), FormatOrderList(userID, targetID, time))
		return err
	}
	// find it.
	return err
}

// RemoveOrderToList remove it.
func RemoveOrderToList(db *sql.Sqlite, userID int64, groupID int64) error {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var addOrderListNeed OrderListStruct
	err := db.Find("orderlist_"+strconv.FormatInt(groupID, 10), &addOrderListNeed, "where order is "+strconv.FormatInt(userID, 10))
	if err != nil {
		return err
	}
	err = db.Del("orderlist_"+strconv.FormatInt(groupID, 10), "where order is "+strconv.FormatInt(userID, 10))
	return err
}

// CheckThisOrderList getList
func CheckThisOrderList(db *sql.Sqlite, userID int64, groupID int64) (OrderUser int64, TargetUSer int64, time string) {
	marryLocker.Lock()
	defer marryLocker.Unlock()
	var addOrderListNeed OrderListStruct
	err := db.Find("orderlist_"+strconv.FormatInt(groupID, 10), &addOrderListNeed, "where order is "+strconv.FormatInt(userID, 10))
	if err != nil {
		return -1, -1, ""
	}
	OrderUser = addOrderListNeed.OrderPerson
	TargetUSer = addOrderListNeed.OrderPerson
	time = addOrderListNeed.Time
	return OrderUser, TargetUSer, time
}

/*
This Path is to use to find the list and delete it.
*/

// FindInList find the list is existed.
func FindInList(list []int, target int) (bool, int) {
	for i, num := range list {
		if num == target {
			return true, i
		}
	}
	return false, -1
}

// FindAndDeleteInThisList the list and delete it.
func FindAndDeleteInThisList(list []int, target int) ([]int, bool) {
	for i, num := range list {
		if num == target {
			list = append(list[:i], list[i+1:]...)
			return list, true
		}
	}
	return list, false
}
