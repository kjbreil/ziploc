package iniUpdater

import (
	"fmt"
	"github.com/go-ini/ini"
	"path/filepath"
	"strings"
)

type U struct {
	D           map[string]map[string]string
	file        string
	sections    []string
	sectionKeys map[string][]string
}

func New() *U {
	var u U
	u.D = make(map[string]map[string]string)
	u.sectionKeys = make(map[string][]string)
	return &u
}

func readIni(name string) (*ini.File, error) {
	if filepath.Ext(name) != ".ini" {
		return nil, fmt.Errorf("%s does not have an INI extension", name)
	}
	return ini.Load(name)
}

// CompareTwo takes the original file and custom file, returning the settings that are new in the custom file
func CompareTwo(original, custom string) (*U, error) {
	u := New()
	_, u.file = filepath.Split(original)
	o, err := readIni(original)
	if err != nil {
		return nil, err
	}
	c, err := readIni(custom)
	if err != nil {
		return nil, err
	}

	// loop over original sections finding matching one in custom, then compare the section
	// if a setting does not exist in custom then add the default, if they match do nothing and if updated in custom add
	for _, originalSection := range o.Sections() {
		customSection := c.Section(originalSection.Name())
		if customSection != nil {
			u.compareSection(originalSection, customSection)
		}
	}
	return u, nil
}

func (u *U) compareSection(left, right *ini.Section) {

	for k, lv := range left.KeysHash() {
		rv := right.Key(k).Value()
		if lv == rv {
			continue
		}
		_, ok := u.D[left.Name()]
		if !ok {
			u.D[left.Name()] = make(map[string]string)
			u.sections = append(u.sections, left.Name())
		}
		u.D[left.Name()][k] = rv
		u.sectionKeys[left.Name()] = append(u.sectionKeys[left.Name()], k)

	}

	// reverso
	for k, rv := range right.KeysHash() {
		lv := right.Key(k).Value()
		if lv == rv {
			continue
		}

		_, ok := u.D[left.Name()]
		if !ok {
			u.D[left.Name()] = make(map[string]string)
		}
		u.D[left.Name()][k] = rv

	}
}

func (u *U) Bytes() []byte {
	var b []byte

	for _, section := range u.sections {
		for _, key := range u.sectionKeys[section] {
			value := u.D[section][key]
			t := fmt.Sprintf("@WIZSET(%s[%s]%s=%s);\r\n", strings.ToUpper(u.file), section, key, value)
			b = append(b, []byte(t)...)
		}
	}
	// TODO: handle top level ini with .\\..\\ (SMSStart.ini)
	s := fmt.Sprintf("@dbHot(%s,SET,%s[*);", strings.ToUpper(u.file), strings.ToUpper(u.file))
	b = append(b, []byte(s)...)
	return b
}
