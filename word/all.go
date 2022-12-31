package word

import (
	"earlang/config"
	"os"
	"strings"
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
		if !strings.HasSuffix(name, ".list") {
			continue
		}
		g, err := NewGroup(strings.TrimSuffix(name, ".list"))
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}
