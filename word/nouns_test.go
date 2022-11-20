package word

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeat(t *testing.T) {
	m := map[string]bool{}
	n := 0
	for _, word := range nouns {
		if _, ok := m[word]; !ok {
			m[word] = true
		} else {
			fmt.Printf("duplicate word `%s` found\n", word)
			n += 1
		}
	}
	assert.Equal(t, 0, n)
}
