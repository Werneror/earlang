package builtin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Word struct {
	English string
	Chinese string
	// 在搜索引擎搜索图片时搜索什么，如果 Query 为空则把 English 作为 query
	Query string
}

func (w Word) Key() string {
	return fmt.Sprintf("%s,%s", w.English, w.Chinese)
}

type Group struct {
	Name  string
	Words []Word
}

func SaveToDisk(dir string) error {
	for _, group := range Groups {
		path := filepath.Join(dir, fmt.Sprintf("%s.txt", group.Name))
		s := make([]string, 0, len(group.Words))
		for _, w := range group.Words {
			s = append(s, fmt.Sprintf("%s,%s,%s", w.English, w.Chinese, w.Query))
		}
		err := os.WriteFile(path, []byte(strings.Join(s, "\n")), os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "failed to save built-in group %s", group.Name)
		}
	}
	return nil
}
