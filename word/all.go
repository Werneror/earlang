package word

import (
	"os"
	"strings"

	"github.com/werneror/earlang/config"
)

func AllGroups() ([]*Group, error) {
	dir, err := os.ReadDir(config.WordDir)
	if err != nil {
		return nil, err
	}
	groups := make([]*Group, 0)
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !strings.HasSuffix(name, config.WordListFileExtension) {
			continue
		}
		g, err := NewGroup(strings.TrimSuffix(name, config.WordListFileExtension))
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
