package arknights

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"sort"
	"strings"
	"time"

	"github.com/fogleman/gg"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"golang.org/x/image/font"
)

var tags = [...]string{"近卫干员", "狙击干员", "重装干员",
	"医疗干员", "辅助干员", "术师干员",
	"特种干员", "先锋干员", "近战位",
	"远程位", "高级资深干员", "控场",
	"爆发", "资深干员", "治疗", "支援",
	"新手", "费用回复", "输出", "生存",
	"群攻", "防护", "减速", "削弱",
	"快速复活", "位移", "召唤", "支援机械"}

var Fonts font.Face

type recruitChar struct {
	Name   string
	Rarity int8
	Tags   []string
}

type RecruitChars []recruitChar

func (r RecruitChars) Len() int           { return len(r) }
func (r RecruitChars) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RecruitChars) Less(i, j int) bool { return r[i].Rarity < r[j].Rarity }

var recruitTags = [...]recruitChar{{"Lancet-2", 0, []string{"医疗干员", "远程位", "治疗", "支援机械"}},
	{"Castle-3", 0, []string{"近卫干员", "近战位", "支援", "支援机械"}},
	{"夜刀", 1, []string{"先锋干员", "近战位", "新手"}},
	{"黑角", 1, []string{"重装干员", "近战位", "新手"}},
	{"巡林者", 1, []string{"狙击干员", "远程位", "新手"}},
	{"杜林", 1, []string{"术师干员", "远程位", "新手"}},
	{"12F", 1, []string{"术师干员", "远程位", "新手"}},
	{"芬", 2, []string{"先锋干员", "近战位", "费用回复"}},
	{"香草", 2, []string{"先锋干员", "近战位", "费用回复"}},
	{"翎羽", 2, []string{"先锋干员", "近战位", "输出", "费用回复"}},
	{"玫兰莎", 2, []string{"近卫干员", "近战位", "输出", "生存"}},
	{"米格鲁", 2, []string{"重装干员", "近战位", "防护"}},
	{"克洛丝", 2, []string{"狙击干员", "远程位", "输出"}},
	{"安德切尔", 2, []string{"狙击干员", "远程位", "输出"}},
	{"炎熔", 2, []string{"术师干员", "远程位", "群攻"}},
	{"芙蓉", 2, []string{"医疗干员", "远程位", "治疗"}},
	{"安赛尔", 2, []string{"医疗干员", "远程位", "治疗"}},
	{"史都华德", 2, []string{"术师干员", "远程位", "输出"}},
	{"梓兰", 2, []string{"辅助干员", "远程位", "减速"}},
	{"夜烟", 3, []string{"术师干员", "远程位", "输出", "削弱"}},
	{"远山", 3, []string{"术师干员", "远程位", "群攻"}},
	{"杰西卡", 3, []string{"狙击干员", "远程位", "输出", "生存"}},
	{"流星", 3, []string{"狙击干员", "远程位", "输出", "削弱"}},
	{"白雪", 3, []string{"狙击干员", "远程位", "群攻", "减速"}},
	{"清道夫", 3, []string{"先锋干员", "近战位", "费用回复", "输出"}},
	{"红豆", 3, []string{"先锋干员", "近战位", "输出", "费用回复"}},
	{"杜宾", 3, []string{"近卫干员", "近战位", "输出", "支援"}},
	{"缠丸", 3, []string{"近卫干员", "近战位", "生存", "输出"}},
	{"霜叶", 3, []string{"近卫干员", "近战位", "减速", "输出"}},
	{"艾丝黛尔", 3, []string{"近卫干员", "近战位", "群攻", "生存"}},
	{"慕斯", 3, []string{"近卫干员", "近战位", "输出"}},
	{"砾", 3, []string{"特种干员", "近战位", "快速复活", "防护"}},
	{"暗索", 3, []string{"特种干员", "近战位", "位移"}},
	{"末药", 3, []string{"医疗干员", "远程位", "治疗"}},
	{"嘉维尔", 3, []string{"医疗干员", "远程位", "治疗"}},
	{"调香师", 3, []string{"医疗干员", "远程位", "治疗"}},
	{"角峰", 3, []string{"重装干员", "近战位", "防护"}},
	{"蛇屠箱", 3, []string{"重装干员", "近战位", "防护"}},
	{"古米", 3, []string{"重装干员", "近战位", "防护", "治疗"}},
	{"地灵", 3, []string{"辅助干员", "远程位", "减速"}},
	{"阿消", 3, []string{"特种干员", "近战位", "位移"}},
	{"白面鸮", 4, []string{"医疗干员", "远程位", "治疗", "支援"}},
	{"凛冬", 4, []string{"先锋干员", "近战位", "费用回复", "支援"}},
	{"德克萨斯", 4, []string{"先锋干员", "近战位", "费用回复", "控场"}},
	{"因陀罗", 4, []string{"近卫干员", "近战位", "输出", "生存"}},
	{"幽灵鲨", 4, []string{"近卫干员", "近战位", "群攻", "生存"}},
	{"蓝毒", 4, []string{"狙击干员", "远程位", "输出"}},
	{"白金", 4, []string{"狙击干员", "远程位", "输出"}},
	{"陨星", 4, []string{"狙击干员", "远程位", "群攻", "削弱"}},
	{"梅尔", 4, []string{"辅助干员", "远程位", "召唤", "控场"}},
	{"赫默", 4, []string{"医疗干员", "远程位", "治疗"}},
	{"华法琳", 4, []string{"医疗干员", "远程位", "治疗", "支援"}},
	{"临光", 4, []string{"重装干员", "近战位", "防护", "治疗"}},
	{"红", 4, []string{"特种干员", "近战位", "快速复活", "控场"}},
	{"雷蛇", 4, []string{"重装干员", "近战位", "防护", "输出"}},
	{"可颂", 4, []string{"重装干员", "近战位", "防护", "位移"}},
	{"火神", 4, []string{"重装干员", "近战位", "生存", "防护", "输出"}},
	{"普罗旺斯", 4, []string{"狙击干员", "远程位", "输出"}},
	{"守林人", 4, []string{"狙击干员", "远程位", "输出", "爆发"}},
	{"崖心", 4, []string{"特种干员", "近战位", "位移", "输出"}},
	{"初雪", 4, []string{"辅助干员", "远程位", "削弱"}},
	{"真理", 4, []string{"辅助干员", "远程位", "减速", "输出"}},
	{"狮蝎", 4, []string{"特种干员", "近战位", "输出", "生存"}},
	{"食铁兽", 4, []string{"特种干员", "近战位", "位移", "减速"}},
	{"能天使", 5, []string{"狙击干员", "远程位", "输出"}},
	{"推进之王", 5, []string{"先锋干员", "近战位", "费用回复", "输出"}},
	{"伊芙利特", 5, []string{"术师干员", "远程位", "群攻", "削弱"}},
	{"闪灵", 5, []string{"医疗干员", "远程位", "治疗", "支援"}},
	{"夜莺", 5, []string{"医疗干员", "远程位", "治疗", "支援"}},
	{"星熊", 5, []string{"重装干员", "近战位", "防护", "输出"}},
	{"塞雷娅", 5, []string{"重装干员", "近战位", "防护", "治疗", "支援"}},
	{"银灰", 5, []string{"近卫干员", "近战位", "输出", "支援"}},
	{"空爆", 2, []string{"狙击干员", "远程位", "群攻"}},
	{"月见夜", 2, []string{"近卫干员", "近战位", "输出"}},
	{"猎蜂", 3, []string{"近卫干员", "近战位", "输出"}},
	{"夜魔", 4, []string{"术师干员", "远程位", "输出", "治疗", "减速"}},
	{"斯卡蒂", 5, []string{"近卫干员", "近战位", "输出", "生存"}},
	{"陈", 5, []string{"近卫干员", "近战位", "输出", "爆发"}},
	{"诗怀雅", 4, []string{"近卫干员", "近战位", "输出", "支援"}},
	{"格雷伊", 3, []string{"术师干员", "远程位", "群攻", "减速"}},
	{"泡普卡", 2, []string{"近卫干员", "近战位", "群攻", "生存"}},
	{"斑点", 2, []string{"重装干员", "近战位", "防护", "治疗"}},
	{"THRM-EX", 0, []string{"特种干员", "近战位", "爆发", "支援机械"}},
	{"黑", 5, []string{"狙击干员", "远程位", "输出"}},
	{"赫拉格", 5, []string{"近卫干员", "近战位", "输出", "生存"}},
	{"格劳克斯", 4, []string{"辅助干员", "远程位", "减速", "控场"}},
	{"星极", 4, []string{"近卫干员", "近战位", "输出", "防护"}},
	{"苏苏洛", 3, []string{"医疗干员", "远程位", "治疗"}},
	{"桃金娘", 3, []string{"先锋干员", "近战位", "费用回复", "治疗"}},
	{"麦哲伦", 5, []string{"辅助干员", "远程位", "支援", "减速", "输出"}},
	{"送葬人", 4, []string{"狙击干员", "远程位", "群攻"}},
	{"红云", 3, []string{"狙击干员", "远程位", "输出"}},
	{"莫斯提马", 5, []string{"术师干员", "远程位", "群攻", "支援", "控场"}},
	{"槐琥", 4, []string{"特种干员", "近战位", "快速复活", "削弱"}},
	{"清流", 3, []string{"医疗干员", "远程位", "治疗", "支援"}},
	{"梅", 3, []string{"狙击干员", "远程位", "输出", "减速"}},
	{"煌", 5, []string{"近卫干员", "近战位", "输出", "生存"}},
	{"灰喉", 4, []string{"狙击干员", "远程位", "输出"}},
	{"苇草", 4, []string{"先锋干员", "近战位", "费用回复", "输出"}},
	{"布洛卡", 4, []string{"近卫干员", "近战位", "群攻", "生存"}},
	{"安比尔", 3, []string{"狙击干员", "远程位", "输出", "减速"}},
	{"阿", 5, []string{"特种干员", "远程位", "支援", "输出"}},
	{"吽", 4, []string{"重装干员", "近战位", "防护", "治疗"}},
	{"正义骑士号", 0, []string{"狙击干员", "远程位", "支援", "支援机械"}},
}

type tagResult struct {
	Tags   []string
	Result RecruitChars
}

func combine(tags []string, nums int) (tagCombine [][]string) {
	var temp []int
	n := len(tags) - 1
	var dfs func(int)
	dfs = func(cur int) {
		if len(temp)+(n-cur+1) < nums {
			return
		}
		if len(temp) == nums {
			comb := make([]int, nums)
			copy(comb, temp)
			var tempString []string
			for _, value := range temp {
				tempString = append(tempString, tags[value])
			}
			tagCombine = append(tagCombine, tempString)
			return
		}
		temp = append(temp, cur)
		dfs(cur + 1)
		temp = temp[:len(temp)-1]
		dfs(cur + 1)
	}
	dfs(0)
	return
}

func isSubTags(tags []string, tags2 recruitChar) bool {
	has6tag := false
	for _, tag := range tags {
		flag := false
		for _, tag2 := range tags2.Tags {
			if tag == tag2 {
				flag = true
				break
			}
			if tag == "高级资深干员" && tags2.Rarity == 5 {
				flag = true
				has6tag = true
				break
			}
			if tag == "资深干员" && tags2.Rarity == 4 {
				flag = true
				break
			}
		}
		if !flag {
			return false
		}
	}
	if !has6tag && tags2.Rarity == 5 {
		return false
	}
	return true
}

func tagsAnalysis(tags []string) (tagResults []tagResult) {
	for nums := 1; nums <= 3; nums++ {
		for _, chooseTags := range combine(tags, nums) {
			var tmp tagResult
			var hasLowRarityChar bool = false
			for _, tags := range recruitTags {
				if isSubTags(chooseTags, tags) {
					if tags.Rarity == 2 {
						hasLowRarityChar = true
					}
					tmp.Tags = chooseTags
					tmp.Result = append(tmp.Result, tags)
				}
			}
			if len(tmp.Result) > 0 && !hasLowRarityChar {
				sort.Sort(tmp.Result)
				tagResults = append(tagResults, tmp)
			}
		}
	}
	return
}

func recruit(ctx *zero.Ctx) {
	ctx.SendChain(
		message.Text("请在60s内发送你的公招截图"),
	)
	rule := func(ctx *zero.Ctx) bool {
		for _, elem := range ctx.Event.Message {
			if elem.Type == "image" {
				ocrTags := make([]string, 0)
				ocrResult := ctx.OCRImage(elem.Data["file"]).Get("texts.#.text").Array()
				for _, text := range ocrResult {
					for _, tag := range tags {
						if tag == text.Str {
							ocrTags = append(ocrTags, tag)
							break
						}
					}
				}
				text := fmt.Sprintf("识别标签:%s\n", strings.Join(ocrTags, " "))
				for _, tmp := range tagsAnalysis(ocrTags) {
					var resultStrSlice []string
					for _, char := range tmp.Result {
						resultStrSlice = append(resultStrSlice, fmt.Sprintf("%s(%d)", char.Name, char.Rarity+1))
					}
					text += fmt.Sprintf("[%s]\n%s\n\n",
						strings.Join(tmp.Tags, " "),
						strings.Join(resultStrSlice, " "),
					)
				}
				img := drawStrImage(text, 500.0)
				ctx.SendChain(message.Image("base64://" + helper.BytesToString(Image2Base64(img))))
				return true
			}
		}
		ctx.SendChain(message.Text("没有发现图片"))
		return true
	}
	next := zero.NewFutureEvent("message", 999, false, zero.CheckUser(ctx.Event.UserID), rule)
	recv, cancel := next.Repeat()
	select {
	case <-time.After(time.Minute):
		ctx.SendChain(message.Text("已超时"))
		cancel()
	case <-recv:
		cancel()
	}
}

func drawStrImage(str string, W float64) image.Image {
	w := 0.0
	h := Fonts.Metrics().Height / 64
	totalW := W
	totalH := h
	for _, char := range str {
		charWeight := float64(font.MeasureString(Fonts, string(char)) >> 6)
		if w+charWeight > totalW || string(char) == "\n" {
			if string(char) == "\n" {
				w = 0
			} else {
				w = charWeight
			}
			totalH += h
		} else {
			w += charWeight
		}
	}
	var img *gg.Context
	img = gg.NewContext(int(totalW), int(totalH))
	img.SetHexColor("FFFFFF")
	img.Clear()
	img.SetFontFace(Fonts)
	img.SetHexColor("000000")
	w = 0.0
	totalH = h
	for _, char := range str {
		charWeight := float64(font.MeasureString(Fonts, string(char)) >> 6)
		if w+charWeight > totalW || string(char) == "\n" {
			totalH += h
			if string(char) != "\n" {
				w = charWeight
				img.DrawString(string(char), 0.0, float64(totalH))
			} else {
				w = 0
			}
		} else {
			img.DrawString(string(char), w, float64(totalH))
			w += charWeight
		}
	}
	return img.Image()
}

func Image2Base64(image image.Image) []byte {
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 95
	_ = jpeg.Encode(encoder, image, &opt)
	err := encoder.Close()
	if err != nil {
		return nil
	}
	return buffer.Bytes()
}
