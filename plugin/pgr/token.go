package pgr

import (
	"github.com/tidwall/gjson"
	"strconv"
)

// may we use a way to store data?

func init() {
	//
}

func RawJsonParse(raw string) (qq int64, Session string) {
	getQQ := gjson.Get(raw, "qq").String()
	strToInt, err := strconv.ParseInt(getQQ, 10, 64)
	if err != nil {
		return 0, ""
	}
	getSession := gjson.Get(raw, "session").String()
	return strToInt, getSession
}
