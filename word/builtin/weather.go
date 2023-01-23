package builtin

var Weather = Group{
	Name: "weather(builtin)",
	Words: []Word{
		{English: "sunshine", Chinese: "阳光"},
		{English: "rain", Chinese: "雨"},
		{English: "snow", Chinese: "雪"},
		{English: "hail", Chinese: "冰雹"},
		{English: "sleet", Chinese: "冰雨"},
		{English: "fog", Chinese: "雾"},
		{English: "cloud", Chinese: "云"},
		{English: "rainbow", Chinese: "彩虹"},
		{English: "wind", Chinese: "大风"},
		{English: "lightning", Chinese: "闪电", Query: "闪电"},
		{English: "storm", Chinese: "暴风雨", Query: "暴风雨"},
		{English: "tornado", Chinese: "龙卷风"},
		{English: "hurricane", Chinese: "飓风"},
		{English: "flood", Chinese: "洪水"},
		{English: "frost", Chinese: "霜"},
		{English: "ice", Chinese: "冰"},
		{English: "drought", Chinese: "干旱"},
		{English: "hot", Chinese: "热", Query: "炎热的"},
		{English: "cold", Chinese: "冷", Query: "冷的"},
		{English: "thermometer", Chinese: "温度计"},
		{English: "barometer", Chinese: "气压表", Query: "气压表"},
	},
}
