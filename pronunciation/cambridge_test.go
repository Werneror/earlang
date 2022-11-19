package pronunciation

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseWordPron(t *testing.T) {
	pron, err := WordPron("bus", "uk")
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(pron, ".mp3"))
	fmt.Println(pron)
}

func TestWordPron(t *testing.T) {
	c := &CambridgeDictionary{}
	pron, err := c.WordPron("bus", "uk")
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(pron, ".mp3"))
	fmt.Println(pron)
}
