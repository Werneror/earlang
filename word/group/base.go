package group

type Word struct {
	English string
	Chinese string
}

type Group struct {
	Name  string
	Words []Word
}
