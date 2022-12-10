package word

import (
	"os"

	"github.com/sirupsen/logrus"
)

type pointer struct {
	pointer     int
	persistFile string
}

func (p *pointer) addOne() int {
	p.pointer = p.pointer + 1
	err := SavePointerToFile(p.persistFile, p.pointer)
	if err != nil {
		logrus.Errorf("failed to save pointer to file: %v", err)
	}
	return p.pointer
}

func (p *pointer) minusOne() int {
	if p.pointer == 0 {
		return p.pointer
	}
	p.pointer = p.pointer - 1
	err := SavePointerToFile(p.persistFile, p.pointer)
	if err != nil {
		logrus.Errorf("failed to save pointer to file: %v", err)
	}
	return p.pointer
}

func (p *pointer) getValue() int {
	return p.pointer
}

func (p *pointer) loadValue() error {
	var err error
	p.pointer, err = LoadPointerFromFile(p.persistFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (p *pointer) setValue(v int) {
	p.pointer = v
	err := SavePointerToFile(p.persistFile, p.pointer)
	if err != nil {
		logrus.Errorf("failed to save pointer to file: %v", err)
	}
}

func newPointer(persistFile string) *pointer {
	p := &pointer{persistFile: persistFile}
	err := p.loadValue()
	if err != nil {
		logrus.Errorf("failed to load pointer from file: %v", err)
	}
	return p
}
