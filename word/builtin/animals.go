package builtin

var Animals = Group{
	Name: "animals(builtin)",
	Words: []Word{
		{English: "fish", Chinese: "鱼"},
		{English: "bird", Chinese: "鸟"},
		{English: "bug", Chinese: "虫子"},
		{English: "beast", Chinese: "野兽", Query: "野兽"},
	},
}
