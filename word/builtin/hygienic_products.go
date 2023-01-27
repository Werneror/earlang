package builtin

var HygienicProducts = Group{
	Name: "hygienic products(builtin)",
	Words: []Word{
		{English: "perfume", Chinese: "香水"},
		{English: "razor", Chinese: "刮胡刀"},
		{English: "shampoo", Chinese: "洗发水", Query: "洗发水"},
		{English: "tampons", Chinese: "卫生棉条"},
		{English: "bleach", Chinese: "漂白水", Query: "漂白剂"},
		{English: "detergent", Chinese: "洗衣粉", Query: "洗衣粉"},
		{English: "duster", Chinese: "抹布", Query: "抹布"},
		{English: "soap", Chinese: "肥皂"},
		{English: "broom", Chinese: "扫帚"},
		{English: "bucket", Chinese: "桶子"},
		{English: "mop", Chinese: "拖把"},
		{English: "sponge", Chinese: "海绵"},
	},
}
