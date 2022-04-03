package reborn

import (
	wr "github.com/mroth/weightedrand"
)

var (
	gender, _ = wr.NewChooser(
		wr.Choice{Item: "男孩子", Weight: 33707},
		wr.Choice{Item: "女孩子", Weight: 39292},
		wr.Choice{Item: "雌雄同体", Weight: 1001},
		wr.Choice{Item: "猫猫!", Weight: 10000},
		wr.Choice{Item: "狗狗!", Weight: 10000},
		wr.Choice{Item: "🐉!", Weight: 3000},
		wr.Choice{Item: "龙猫~", Weight: 3000},
	)
)

func randcoun() string {
	return areac.Pick().(string)
}

func randgen() string {
	return gender.Pick().(string)
}
