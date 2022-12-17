package Unfinished // Package tools only accessible to PC

import (
	"encoding/json"
	"github.com/FloatTech/ZeroBot-Plugin/plugin/tools"
	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"os"
	"strconv"
)

const (
	apexcheck = "https://api.mozambiquehe.re/"
)

func init() {
	accessload := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		apexToken := tools.engine.DataFolder() + "apexToken.txt"
		_, err := os.ReadFile(apexToken)
		if err != nil {
			ctx.SendChain(message.Text("还没有填写ApexToken~,请在本插件Floder下生成apexToken.txt文件，填入apexToken"))
			return false
		}
		return true
	})

	tools.engine.OnFullMatch("/apex", accessload, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		apexToken := tools.engine.DataFolder() + "apexToken.txt"
		apexTokenByte, _ := os.ReadFile(apexToken)
		token := helper.BytesToString(apexTokenByte)
		uid := LoadUserPlayID(ctx, strconv.FormatInt(ctx.Event.UserID, 10))
		playStatus, err := web.GetData(apexcheck + "bridge?auth=" + token + "&player=" + uid + "&platform=PC")
		if err != nil {
			ctx.SendChain(message.Text("Something went wrong,", err))
			return
		}
		formatUserStatus := helper.BytesToString(playStatus)
		playerRawName := gjson.Get(formatUserStatus, "global.name")
		playerRawLvl := gjson.Get(formatUserStatus, "global.level")
		playerRawIfbanned := gjson.Get(formatUserStatus, "global.bans.isActive").Bool()
		playerRawLastbannedReason := gjson.Get(formatUserStatus, "global.bans.last_banReason").String()
		playerRawBattleLevel := gjson.Get(formatUserStatus, "rank.name").String()
		playerRawBattleScore1 := gjson.Get(formatUserStatus, "rank.rankDiv")
		playerRawBattleScore2 := gjson.Get(formatUserStatus, "rank.rankScore")
		playerRawArenaName := gjson.Get(formatUserStatus, "arena.rankName")
		playerRawArenaScore1 := gjson.Get(formatUserStatus, "arena.rankDiv")
		playerRawArenaScore2 := gjson.Get(formatUserStatus, "arena.rankScore")
	})
}

func InsertPlayUser(userID string, playID string) (err error) { // 一次轮询，填入对应玩家信息
	var PlayerData map[string]interface{}
	playList := tools.engine.DataFolder() + "apex_playerlist.json"
	data, err := os.ReadFile(playList)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.WriteFile(playList, []byte("{}"), 0777)
		} else {
			return err
		}
	}
	err = json.Unmarshal(data, &PlayerData)
	if err != nil {
		return err
	}
	PlayerData[userID] = playID
	newData, err := json.Marshal(PlayerData)
	if err != nil {
		return err
	}
	_ = os.WriteFile(playList, newData, 0777)
	return nil
}

func LoadUserPlayID(ctx *zero.Ctx, userID string) (playUID string) {
	var data map[string]string
	playListJson := tools.engine.DataFolder() + "apex_playerlist.json"
	playList, err := os.ReadFile(playListJson)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(message.Text("发生了未知错误awa,一定是夹子的问题.jpg")))
		return
	}
	err = json.Unmarshal(playList, &data)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(message.Text("发生了未知错误awa,一定是夹子的问题.jpg")))
		return
	}
	result := data[userID]
	if result == "" {
		ctx.Send(message.ReplyWithMessage(message.Text("你还没有登记数据呢！")))
		return
	}
	return result
}

func QueryUserNameToID(ctx *zero.Ctx, name string) (UserName string, UserID string) {
	apexTokenByte, _ := os.ReadFile(tools.engine.DataFolder() + "apexToken.txt")
	apexToken := helper.BytesToString(apexTokenByte)
	url := apexcheck + "nametouid?auth=" + apexToken + "&player=" + name + "&platform=PC"
	resultByte, err := web.GetData(url)
	resultin := helper.BytesToString(resultByte)
	if err != nil {
		ctx.Send(message.ReplyWithMessage(message.Text("未找到UID")))
		return "", ""
	}
	UserName = gjson.Get(resultin, "name").String()
	UserID = gjson.Get(resultin, "uid").String()
	return UserName, UserID
}

func ReverseChart(name string) (result string) {
	// formatMap removed Level2 and Level3 Due to Same dicts
	formatMap := map[string]string{
		"Kings Canyon":  "诸王峡谷",
		"World\"s Edge": "世界尽头",
		"Olympus":       "奥林匹斯",
		"Storm Point":   "风暴点",
		"Broken Moon":   "破碎月亮",
		"Encore":        "再来一次",
		"Habitat":       "栖息地 4",
		"Overflow":      "熔岩流",
		"Phase runner":  "相位穿梭器",
		"Party crasher": "派对破坏者",
		// Barrels
		"barrel_stabilizer": "枪管稳定器",
		"laser_sight":       "激光瞄准镜",
		// Mags
		"extended_energy_mag": "加长式能量弹匣",
		"extended_heavy_mag":  "加长式重型弹匣",
		"extended_light_mag":  "加长式轻型弹匣",
		"extended_sniper_mag": "加长狙击弹匣",
		"shotgun_bolt":        "霰弹枪枪栓",
		// Optics
		"optic_hcog_classic":          "单倍全息衍射式瞄准镜`经典`",
		"optic_holo":                  "单倍幻影",
		"optic_variable_holo":         "单倍至 2 倍可调节式幻影瞄准镜",
		"optic_hcog_bruiser":          "2 倍全息衍射式瞄准镜`格斗家`",
		"optic_sniper":                "6 倍狙击手",
		"optic_variable_aog":          "2 倍至 4 倍可调节式高级光学瞄准镜",
		"optic_hcog_ranger":           "3 倍全息衍射式瞄准镜游侠",
		"optic_variable_sniper":       "4 倍至 8 倍可调节式狙击手",
		"optic_digital_threat":        "单倍数字化威胁",
		"optic_digital_sniper_threat": "4 倍至 10 倍数字化狙击威胁",
		// Hop_Ups
		"anvil_receiver":       "铁砧接收器",
		"double_tap_trigger":   "双发扳机",
		"skullpiercer_rifling": "穿颅器",
		"turbocharger":         "涡轮增压器",
		// Stocks
		"standard_stock": "标准枪托",
		"sniper_stock":   "狙击枪枪托",
		// Gear
		"helmet":                "头盔",
		"evo_shield":            "进化护盾",
		"knockdown_shield":      "击倒护盾",
		"backpack":              "背包",
		"survival":              "隔热板",
		"mobile_respawn_beacon": "移动重生信标",
		// Regen
		"shield_cell":         "小型护盾电池",
		"syringe":             "注射器",
		"large_shield_cell":   "护盾电池",
		"med_kit":             "医疗箱",
		"phoenix_kit":         "凤凰治疗包",
		"ultimate_accelerant": "绝招加速剂",
		// Ammo
		"special": "特殊弹药",
		"energy":  "能量弹药",
		"heavy":   "重型弹药",
		"light":   "轻型弹药",
		"shotgun": "霰弹枪弹药",
		"sniper":  "狙击弹药",
		// Other
		"evo_points": "进化点数",
		"Rare":       "稀有",
		"Epic":       "史诗",
		// Rank
		"Unranked":      "菜鸟",
		"Bronze":        "青铜",
		"Silver":        "白银",
		"Gold":          "黄金",
		"Platinum":      "白金",
		"Diamond":       "钻石",
		"Master":        "大师",
		"Apex Predator": "Apex 猎杀者",
		// Legends
		"Bloodhound": "寻血猎犬",
		"Gibraltar":  "直布罗陀",
		"Lifeline":   "命脉",
		"Pathfinder": "探路者",
		"Wraith":     "恶灵",
		"Bangalore":  "班加罗尔",
		"Caustic":    "侵蚀",
		"Mirage":     "幻象",
		"Octane":     "动力小子",
		"Wattson":    "沃特森",
		"Crypto":     "密客",
		"Revenant":   "亡灵",
		"Loba":       "罗芭",
		"Rampart":    "兰伯特",
		"Horizon":    "地平线",
		"Fuse":       "暴雷",
		"Valkyrie":   "瓦尔基里",
		"Seer":       "希尔",
		"Ash":        "艾许",
		"Mad Maggie": "疯玛吉",
		"Newcastle":  "纽卡斯尔",
		"Vantage":    "万蒂奇",
		"Catalyst":   "卡特莉丝",
		// Level
		"Common":    "等级1",
		"Legendary": "等级4",
		"Mythic":    "等级5",
		// Weapon
		"peacekeeper": "和平捍卫者",
		"spitfire":    "喷火轻机枪",
		// API
		"offline":                    "离线",
		"online":                     "在线",
		"0":                          "否",
		"1":                          "是",
		"invite":                     "邀请",
		"open":                       "打开",
		"inLobby":                    "在大厅",
		"inMatch":                    "比赛中",
		"true":                       "是",
		"false":                      "否",
		"COMPETITIVE_DODGE_COOLDOWN": "竞技逃跑冷却",
		"None":                       "无",
	}
	result, err := formatMap[name]
	if err != true {
		result = name
	}
	return result
}
