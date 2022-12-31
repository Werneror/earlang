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
}

type Group struct {
	Name  string
	Words []Word
}

func SaveToDisk(dir string) error {
	for _, group := range Groups {
		path := filepath.Join(dir, fmt.Sprintf("%s.list", group.Name))
		// TODO: 已经存在的单词组如果有新增或删除单词要怎么办？
		if _, err := os.Stat(path); os.IsNotExist(err) {
			s := make([]string, 0, len(group.Words))
			for _, w := range group.Words {
				s = append(s, fmt.Sprintf("%s,%s", w.English, w.Chinese))
			}
			err := os.WriteFile(path, []byte(strings.Join(s, "\n")), os.ModePerm)
			if err != nil {
				return errors.Wrapf(err, "failed to save built-in group %s", group.Name)
			}
		}
	}
	return nil
}
