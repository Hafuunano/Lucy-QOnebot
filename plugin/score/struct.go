package score

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func init() {
	go func() {
		sdb = initialize(engine.DataFolder() + "score.db")
		log.Println("[score]加载score数据库")
	}()
}

// scoredb 分数数据库
type scoredb gorm.DB

// scoretable 分数结构体
type scoretable struct {
	UID   int64 `gorm:"column:uid;primary_key"`
	Score int   `gorm:"column:score;default:0"`
}

// signintable 签到结构体
type signintable struct {
	UID       int64 `gorm:"column:uid;primary_key"`
	Count     int   `gorm:"column:count;default:0"`
	Coins     int   `gorm:"column:coins;default:0"`
	UpdatedAt time.Time
}

// globalTable 总体结构体
type globaltable struct {
	Counttime int `gorm:"column:counttime;default:0"`
	Times     string
}

// TableName ...
func (globaltable) TableName() string {
	return "global"
}

// TableName ...
func (scoretable) TableName() string {
	return "score"
}

// TableName ...
func (signintable) TableName() string {
	return "sign_in"
}

// initialize 初始化ScoreDB数据库
func initialize(dbpath string) *scoredb {
	var err error
	if _, err = os.Stat(dbpath); err != nil || os.IsNotExist(err) {
		// 生成文件
		f, err := os.Create(dbpath)
		if err != nil {
			return nil
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)
	}
	gdb, err := gorm.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
	gdb.AutoMigrate(&scoretable{}).AutoMigrate(&signintable{}).AutoMigrate(&globaltable{})
	return (*scoredb)(gdb)
}

// Close ...
func (sdb *scoredb) Close() error {
	db := (*gorm.DB)(sdb)
	return db.Close()
}

// GetScoreByUID 取得分数
func (sdb *scoredb) GetScoreByUID(uid int64) (s scoretable) {
	db := (*gorm.DB)(sdb)
	db.Debug().Model(&scoretable{}).FirstOrCreate(&s, "uid = ? ", uid)
	return s
}

// InsertOrUpdateScoreByUID 插入或更新打卡累计 初始化更新临时保存
func (sdb *scoredb) InsertOrUpdateScoreByUID(uid int64, score int) (err error) {
	db := (*gorm.DB)(sdb)
	s := scoretable{
		UID:   uid,
		Score: score,
	}
	if err = db.Debug().Model(&scoretable{}).First(&s, "uid = ? ", uid).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			err = db.Debug().Model(&scoretable{}).Create(&s).Error // newUser not user
		}
	} else {
		err = db.Debug().Model(&scoretable{}).Where("uid = ? ", uid).Update(
			map[string]interface{}{
				"score": score,
			}).Error
	}
	return
}

// GetSignInByUID 取得签到次数
func (sdb *scoredb) GetSignInByUID(uid int64) (si signintable) {
	db := (*gorm.DB)(sdb)
	db.Debug().Model(&signintable{}).FirstOrCreate(&si, "uid = ? ", uid)
	return si
}

// GetCurrentCount 取得现在的人数
func (sdb *scoredb) GetCurrentCount(times string) (si globaltable) {
	db := (*gorm.DB)(sdb)
	db.Debug().Model(&globaltable{}).FirstOrCreate(&si, "times = ? ", times)
	return si
}

// InsertOrUpdateSignInCountByUID 插入或更新签到次数
func (sdb *scoredb) InsertOrUpdateSignInCountByUID(uid int64, count int) (err error) {
	db := (*gorm.DB)(sdb)
	si := signintable{
		UID:   uid,
		Count: count,
	}
	if err = db.Debug().Model(&signintable{}).First(&si, "uid = ? ", uid).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Debug().Model(&signintable{}).Create(&si) // newUser not user
		}
	} else {
		err = db.Debug().Model(&signintable{}).Where("uid = ? ", uid).Update(
			map[string]interface{}{
				"count": count,
			}).Error
	}
	return
}

func (sdb *scoredb) InsertUserCoins(uid int64, coins int) (err error) { // 修改金币数值
	db := (*gorm.DB)(sdb)
	si := signintable{
		UID:   uid,
		Coins: coins,
	}
	if err = db.Debug().Model(&signintable{}).First(&si, "uid = ? ", uid).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			db.Debug().Model(&signintable{}).Create(&si) // newUser not user
		}
	} else {
		err = db.Debug().Model(&signintable{}).Where("uid = ? ", uid).Update(
			map[string]interface{}{
				"coins": coins,
			}).Error
	}
	return
}

func (sdb *scoredb) UpdateUserTime(counttime int, times string) (err error) {
	db := (*gorm.DB)(sdb)
	si := globaltable{
		Counttime: counttime,
		Times:     times,
	}
	if err = db.Debug().Model(&globaltable{}).First(&si, "times = ?", times).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			db.Debug().Model(&globaltable{}).Create(&si) // judge if the table remain.
		}
	} else {
		// do your favor
		err = db.Debug().Model(&globaltable{}).Where("times = ?", times).Update(map[string]interface{}{
			"counttime": counttime,
		}).Error
	}
	return err
}

func checkUserCoins(coins int) bool { // 参与一次15个柠檬片
	if coins-50 < 0 {
		return false
	} else {
		return true
	}
}

func getHourWord(t time.Time) (reply string) {
	switch {
	case 5 <= t.Hour() && t.Hour() < 12:
		reply = "今天又是新的一天呢ww"
	case 12 <= t.Hour() && t.Hour() < 14:
		reply = "吃饭了嘛w~如果没有快去吃饭哦w"
	case 14 <= t.Hour() && t.Hour() < 19:
		reply = "记得多去补点水呢~ww,辛苦了哦w"
	case 19 <= t.Hour() && t.Hour() < 24:
		reply = "晚上好吖w 今天过的开心嘛ww"
	case 0 <= t.Hour() && t.Hour() < 5:
		reply = "快去睡觉~已经很晚了w"
	default:
	}
	return
}

func getLevel(count int) int {
	for k, v := range levelArray {
		if count == v {
			return k
		} else if count < v {
			return k - 1
		}
	}
	return -1
}
