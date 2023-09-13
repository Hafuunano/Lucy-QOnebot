package quote

import (
	sql "github.com/FloatTech/sqlite"
	"strconv"
	"sync"
	"time"
)

type Corruption struct {
	TrackID   int64  `db:"trackid"`
	QQ        int64  `db:"qq"`
	Msg       string `db:"msg"`
	HandlerQQ int64  `db:"handler"`
	Time      int64  `db:"time"` // tip : time means save time here.
}

var CorruptionSide = &sql.Sqlite{}
var CorruptLocker = sync.Mutex{}

func init() {
	CorruptionSide.DBPath = engine.DataFolder() + "quote.db"
	err := CorruptionSide.Open(time.Hour * 24)
	if err != nil {
		panic(err)
	}
}

/*
// CreateQuoteChannelForGroup For Group, Group Type Should Be = groupquote_{qid}
func CreateQuoteChannelForGroup(id int64) error {
	CorruptLocker.Lock()
	defer CorruptLocker.Unlock()
	return CorruptionSide.Create("groupquote_"+strconv.Itoa(int(id)), Corruption{})
}

*/

// InsertUserCorruptionTarget Insert Target
func InsertUserCorruptionTarget(InsertType Corruption, id int64) {
	CorruptLocker.Lock()
	defer CorruptLocker.Unlock()

	err := CorruptionSide.Insert("groupquote_"+strconv.Itoa(int(id)), &InsertType)
	if err != nil {
		// create new target insertaion.
		_ = CorruptionSide.Create("groupquote_"+strconv.Itoa(int(id)), &Corruption{})
		_ = CorruptionSide.Insert("groupquote_"+strconv.Itoa(int(id)), &InsertType)
	}
}

// IndexOfDataCorruption Index Of
func IndexOfDataCorruption(id int64) Corruption {
	CorruptLocker.Lock()
	defer CorruptLocker.Unlock()
	var tamper Corruption
	err := CorruptionSide.Pick("groupquote_"+strconv.Itoa(int(id)), &tamper)
	if err != nil {
		return Corruption{TrackID: 0}
	}
	return tamper
}

// ReferDataCorruption Refer A Trackerid
func ReferDataCorruption(id int64, trackerid string) Corruption {
	CorruptLocker.Lock()
	defer CorruptLocker.Unlock()
	var tamper Corruption
	err := CorruptionSide.Find("groupquote_"+strconv.Itoa(int(id)), &tamper, "WHERE trackid is "+trackerid)
	if err != nil {
		return Corruption{TrackID: 0}
	}
	return tamper
}

// RemoveIndexData Index Of
func RemoveIndexData(trackid string, id int64) error {
	CorruptLocker.Lock()
	defer CorruptLocker.Unlock()
	return CorruptionSide.Del("groupquote_"+strconv.Itoa(int(id)), "Where trackid = "+trackid)
}

func RemoveAllIndexData(id int64) error {
	CorruptLocker.Lock()
	defer CorruptLocker.Unlock()
	return CorruptionSide.Del("groupquote_"+strconv.Itoa(int(id)), "WHERE id = 0")
}
