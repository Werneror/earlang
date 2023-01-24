package builtin

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	wordMap := map[string]string{}
	for _, group := range Groups {
		for _, word := range group.Words {
			if groupName, ok := wordMap[word.English]; ok {
				fmt.Printf("%s [%s] [%s]\n", word.English, group.Name, groupName)
			} else {
				wordMap[word.English] = group.Name
			}
		}
	}
}
