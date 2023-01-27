package builtin

import (
	"fmt"
	"testing"
)

func TestRepeat(t *testing.T) {
	wordMap := map[string]string{}
	for _, group := range Groups {
		for _, word := range group.Words {
			if groupName, ok := wordMap[word.Key()]; ok {
				fmt.Printf("%s [%s] [%s]\n", word.Key(), group.Name, groupName)
			} else {
				wordMap[word.Key()] = group.Name
			}
		}
	}
}
