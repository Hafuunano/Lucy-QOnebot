package steam // Package steam listener by using corn

import (
	"fmt"
	"github.com/FloatTech/floatbox/binary"
	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/fumiama/cron"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	URL               = "https://steam-api.impart.icu/"                          // steam API 调用地址
	StatusURL         = "ISteamUser/GetPlayerSummaries/v2/?key=%+v&steamids=%+v" // 根据用户steamID获取用户状态
	TableListenPlayer = "listen_player"
)

var (
	database = &streamDB{
		db: &sql.Sqlite{},
	}
	// 开启并检查数据库链接
	getDB = fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		database.db.DBPath = engine.DataFolder() + "steam.db"
		err := database.db.Open(time.Hour * 24)
		if err != nil {
			ctx.SendChain(message.Text("[steam] ERROR: ", err))
			return false
		}
		if err = database.db.Create(TableListenPlayer, &player{}); err != nil {
			ctx.SendChain(message.Text("[steam] ERROR: ", err))
			return false
		}
		// 校验密钥是否初始化
		if apiKey == "" {
			ctx.SendChain(message.Text("[steam] ERROR: Steam API 密钥尚未初始化"))
			return false
		}
		return true
	})
	apiKey       string
	steamKeyFile = engine.DataFolder() + "steamapikey.txt"
	engine       = control.Register("steam", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "- Steam Listener \n",
		PrivateDataFolder: "steam",
	})
)

// streamDB 继承方法的存储结构
type streamDB struct {
	db *sql.Sqlite
	sync.RWMutex
}

// player 用户状态存储结构体
type player struct {
	SteamID       string `json:"steam_id"`        // 绑定用户标识ID
	PersonaName   string `json:"persona_name"`    // 用户昵称
	Target        string `json:"target"`          // 信息推送群组
	GameID        string `json:"game_id"`         // 游戏ID
	GameExtraInfo string `json:"game_extra_info"` // 游戏信息
	LastUpdate    int64  `json:"last_update"`     // 更新时间
}

func init() {
	go func() {
		if file.IsNotExist(steamKeyFile) {
			_, err := os.Create(steamKeyFile)
			if err != nil {
				panic(err)
			}
		}
		apikey, err := os.ReadFile(steamKeyFile)
		if err != nil {
			panic(err)
		}
		apiKey = binary.BytesToString(apikey)
		c := cron.New(cron.WithSeconds())
		spec := "*/30 * * * * *"
		var ctx *zero.Ctx
		c.AddFunc(spec, func() {
			listenUserChange(ctx)
		})
		c.Start()
	}()
	engine.OnFullMatch("拉取steam绑定用户状态", getDB).SetBlock(true).Handle(listenUserChange)
	engine.OnRegex(`^/steam\sbind\s*(\d+)$`, zero.OnlyGroup, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		steamID := ctx.State["regex_matched"].([]string)[1]
		// 获取用户状态
		playerStatus, err := getPlayerStatus([]string{steamID})
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err, "\nEXP: cannot get user info("))
			return
		}
		if len(playerStatus) == 0 {
			ctx.SendChain(message.Text("Couldn't find the steam user?"))
			return
		}
		playerData := playerStatus[0]
		// 判断用户是否已经初始化：若未初始化，通过用户的steamID获取当前状态并初始化；若已经初始化则更新用户信息
		info, err := database.find(steamID)
		if err != nil {
			ctx.SendChain(message.Text("[steam] ERROR: ", err, "\nEXP: Unexpected err."))
			return
		}
		// 处理数据
		groupID := strconv.FormatInt(ctx.Event.GroupID, 10)
		if info.Target == "" {
			info = player{
				SteamID:       steamID,
				PersonaName:   playerData.PersonaName,
				Target:        groupID,
				GameID:        playerData.GameID,
				GameExtraInfo: playerData.GameExtraInfo,
				LastUpdate:    time.Now().Unix(),
			}
		} else if !strings.Contains(info.Target, groupID) {
			info.Target = strings.Join([]string{info.Target, groupID}, ",")
		}
		// 更新数据库
		if err = database.update(info); err != nil {
			ctx.SendChain(message.Text("ERROR: ", err, "\n: cannnot update the current sqlite."))
			return
		}
		ctx.SendChain(message.Text("oki!"))
	})
	engine.OnRegex(`^/steam\sunbind\s*(\d+)$`, zero.OnlyGroup, getDB).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		steamID := ctx.State["regex_matched"].([]string)[1]
		groupID := strconv.FormatInt(ctx.Event.GroupID, 10)
		// 判断是否已经绑定该steamID，若已绑定就将群列表从推送群列表钟去除
		info, err := database.findWithGroupID(steamID, groupID)
		if err != nil {
			ctx.SendChain(message.Text("ERROR: ", err, "\n"))
			return
		}
		if info.SteamID == "" {
			ctx.SendChain(message.Text("ERROR: cannot find the id?"))
			return
		}
		// 从绑定列表中剔除需要删除的对象
		targets := strings.Split(info.Target, ",")
		newTargets := make([]string, 0)
		for _, target := range targets {
			if target == groupID {
				continue
			}
			newTargets = append(newTargets, target)
		}
		if len(newTargets) == 0 {
			if err = database.del(steamID); err != nil {
				ctx.SendChain(message.Text("ERROR: ", err, "\nEXP: Database err"))
				return
			}
		} else {
			info.Target = strings.Join(newTargets, ",")
			if err = database.update(info); err != nil {
				ctx.SendChain(message.Text("ERROR: ", err, "\nEXP: Database err"))
				return
			}
		}
		ctx.SendChain(message.Text("设置成功"))
	})
}

// listenUserChange 用于监听用户的信息变化
func listenUserChange(ctx *zero.Ctx) {
	su := zero.BotConfig.SuperUsers[0]
	// 获取所有处于监听状态的用户信息
	infos, err := database.findAll()
	if err != nil {
		return
	}
	if len(infos) == 0 {
		return
	}
	// 收集这波用户的streamId，然后查当前的状态，并建立信息映射表
	streamIds := make([]string, len(infos))
	localPlayerMap := make(map[string]player)
	for i, info := range infos {
		streamIds[i] = info.SteamID
		localPlayerMap[info.SteamID] = info
	}
	// 将所有用户状态查一遍
	playerStatus, err := getPlayerStatus(streamIds)
	if err != nil {
		// 出错就发消息
		ctx.SendPrivateMessage(su, message.Text("ERROR: ", err))
		return
	}
	// 遍历返回的信息做对比，假如信息有变化则发消息
	now := time.Now()
	for _, playerInfo := range playerStatus {
		var msg message.Message
		localInfo := localPlayerMap[playerInfo.SteamID]
		// 排除不需要处理的情况
		if localInfo.GameID == "" && playerInfo.GameID == "" {
			continue
		}
		// 打开游戏
		if localInfo.GameID == "" && playerInfo.GameID != "" {
			msg = append(msg, message.Text(playerInfo.PersonaName, "正在玩", playerInfo.GameExtraInfo))
			localInfo.LastUpdate = now.Unix()
		}
		// 更换游戏
		if localInfo.GameID != "" && playerInfo.GameID != localInfo.GameID && playerInfo.GameID != "" {
			msg = append(msg, message.Text(playerInfo.PersonaName, "玩了", (now.Unix()-localInfo.LastUpdate)/60, "分钟后, 丢下了", localInfo.GameExtraInfo, ", 转头去玩", playerInfo.GameExtraInfo))
			localInfo.LastUpdate = now.Unix()
		}
		// 关闭游戏
		if playerInfo.GameID != localInfo.GameID && playerInfo.GameID == "" {
			msg = append(msg, message.Text(playerInfo.PersonaName, "玩了", (now.Unix()-localInfo.LastUpdate)/60, "分钟后, 关掉了", localInfo.GameExtraInfo))
			localInfo.LastUpdate = 0
		}
		if msg != nil {
			groups := strings.Split(localInfo.Target, ",")
			for _, groupString := range groups {
				group, err := strconv.ParseInt(groupString, 10, 64)
				if err != nil {
					ctx.SendPrivateMessage(su, message.Text("ERROR: ", err, "\nOTHER: SteamID ", localInfo.SteamID))
					continue
				}
				ctx.SendGroupMessage(group, msg)
			}
		}
		// 更新数据
		localInfo.GameID = playerInfo.GameID
		localInfo.GameExtraInfo = playerInfo.GameExtraInfo
		if err = database.update(localInfo); err != nil {
			ctx.SendPrivateMessage(su, message.Text("ERROR: ", err, "\nEXP: 更新数据失败\nOTHER: SteamID ", localInfo.SteamID))
		}
	}
}

// getPlayerStatus 获取用户状态
func getPlayerStatus(streamIds []string) ([]*player, error) {
	players := make([]*player, 0)
	// 拼接请求地址
	url := fmt.Sprintf(URL+StatusURL, apiKey, strings.Join(streamIds, ","))
	// 拉取并解析数据
	data, err := web.GetData(url)
	if err != nil {
		return players, err
	}
	dataStr := binary.BytesToString(data)
	index := gjson.Get(dataStr, "response.players.#").Uint()
	for i := uint64(0); i < index; i++ {
		players = append(players, &player{
			SteamID:       gjson.Get(dataStr, fmt.Sprintf("response.players.%d.steamid", i)).String(),
			PersonaName:   gjson.Get(dataStr, fmt.Sprintf("response.players.%d.personaname", i)).String(),
			GameID:        gjson.Get(dataStr, fmt.Sprintf("response.players.%d.gameid", i)).String(),
			GameExtraInfo: gjson.Get(dataStr, fmt.Sprintf("response.players.%d.gameextrainfo", i)).String(),
		})
	}
	return players, nil
}

// update 如果主键不存在则插入一条新的数据，如果主键存在直接复写
func (sql *streamDB) update(dbInfo player) error {
	sql.Lock()
	defer sql.Unlock()
	return sql.db.Insert(TableListenPlayer, &dbInfo)
}

// find 根据主键查信息
func (sql *streamDB) find(steamID string) (dbInfo player, err error) {
	sql.Lock()
	defer sql.Unlock()
	condition := "where steam_id = " + steamID
	if !sql.db.CanFind(TableListenPlayer, condition) {
		return player{}, nil // 规避没有该用户数据的报错
	}
	err = sql.db.Find(TableListenPlayer, &dbInfo, condition)
	return
}

// findWithGroupID 根据用户steamID和groupID查询信息
func (sql *streamDB) findWithGroupID(steamID string, groupID string) (dbInfo player, err error) {
	sql.Lock()
	defer sql.Unlock()
	condition := "where steam_id = " + steamID + " AND target LIKE '%" + groupID + "%'"
	if !sql.db.CanFind(TableListenPlayer, condition) {
		return player{}, nil // 规避没有该用户数据的报错
	}
	err = sql.db.Find(TableListenPlayer, &dbInfo, condition)
	return
}

// findAll 查询所有库信息
func (sql *streamDB) findAll() (dbInfos []player, err error) {
	sql.Lock()
	defer sql.Unlock()
	info := player{}
	num, err := sql.db.Count(TableListenPlayer)
	if err != nil {
		return
	}
	if num == 0 {
		return []player{}, nil
	}
	err = sql.db.FindFor(TableListenPlayer, &info, "", func() error {
		if info.SteamID != "" {
			dbInfos = append(dbInfos, info)
		}
		return nil
	})
	return
}

// del 删除指定数据
func (sql *streamDB) del(steamID string) error {
	sql.Lock()
	defer sql.Unlock()
	return sql.db.Del(TableListenPlayer, "where steam_id = "+steamID)
}
