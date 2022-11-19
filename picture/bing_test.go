package picture

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWordPictures(t *testing.T) {
	bis := &BingImageSearch{}
	pictures, err := bis.WordPictures("bus", 5)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(pictures))
	for _, picUrl := range pictures {
		fmt.Println(picUrl)
	}
}
