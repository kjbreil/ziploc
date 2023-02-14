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
	checkFor    map[string]map[string]string
	folder      string
}

func New() *U {
	var u U
	u.D = make(map[string]map[string]string)
	u.sectionKeys = make(map[string][]string)
	u.checkFor = make(map[string]map[string]string)
	return &u
}

func readIni(name string) (*ini.File, error) {
	if !strings.EqualFold(filepath.Ext(name), ".ini") {
		return nil, fmt.Errorf("%s does not have an INI extension", name)
	}
	o := ini.LoadOptions{
		Loose:                       false,
		Insensitive:                 false,
		InsensitiveSections:         true,
		InsensitiveKeys:             true,
		IgnoreContinuation:          false,
		IgnoreInlineComment:         true,
		SkipUnrecognizableLines:     false,
		ShortCircuit:                false,
		AllowBooleanKeys:            false,
		AllowShadows:                false,
		AllowNestedValues:           false,
		AllowPythonMultilineValues:  false,
		SpaceBeforeInlineComment:    false,
		UnescapeValueDoubleQuotes:   false,
		UnescapeValueCommentSymbols: false,
		UnparseableSections:         nil,
		KeyValueDelimiters:          "",
		KeyValueDelimiterOnWrite:    "",
		ChildSectionDelimiter:       "",
		PreserveSurroundedQuote:     false,
		DebugFunc:                   nil,
		ReaderBufferSize:            0,
		AllowNonUniqueSections:      false,
	}

	i := ini.Empty(o)
	i.NameMapper = ini.SnackCase
	err := i.Append(name)
	if err != nil {
		return nil, err
	}

	return i, nil
}

// CompareTwo takes the original file and custom file, returning the settings that are new in the custom file
func CompareTwo(original, custom, folder string) (*U, error) {
	u := New()
	u.folder = folder
	_, u.file = filepath.Split(original)
	o, err := readIni(original)
	if err != nil {
		return nil, err
	}
	c, err := readIni(custom)
	if err != nil {
		return nil, err
	}

	// loop over custom sections finding matching one in original, then compare the section
	// if a setting does not exist in custom then add the default, if they match do nothing and if updated in custom add
	for _, customSection := range c.Sections() {
		originalSection := o.Section(customSection.Name())
		if originalSection != nil {
			u.compareSection(originalSection, customSection)
		}
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

func (u *U) setRight(section, key, value string) {
	// they don't match so set the key to the right
	_, ok := u.D[section]
	if !ok {
		u.D[section] = make(map[string]string)
		u.sections = append(u.sections, section)
	}
	u.D[section][key] = value
	u.sectionKeys[section] = append(u.sectionKeys[section], key)
}

func (u *U) setDefault(section, key, value string) {
	_, ok := u.checkFor[section]
	if !ok {
		u.checkFor[section] = make(map[string]string)
	}
	u.checkFor[section][key] = value
	u.setRight(section, key, value)

}

func (u *U) compareSection(left, right *ini.Section) {

	// add anything that is not in the original at the "top"
	for _, k := range right.Keys() {
		if k.Name() == "ftppassword" {
			fmt.Println("pw:", k.Value())
		}
		// if the key is already there continue
		_, ok := u.D[left.Name()][k.Name()]
		if ok {
			continue
		}
		lv := left.Key(k.Name()).Value()

		// get the value of the right key
		rv := right.Key(k.Name()).Value()
		// if they match move on
		if lv == rv {
			continue
		}
		u.setRight(left.Name(), k.Name(), rv)
	}

	// loop over original keys
	for _, k := range left.Keys() {
		// if the key is already there continue
		_, ok := u.D[left.Name()][k.Name()]
		if ok {
			continue
		}
		lv := left.Key(k.Name()).Value()
		if !right.HasKey(k.Name()) {
			// add to the check for map
			// Check and create key with default value
			u.setDefault(left.Name(), k.Name(), lv)
			continue
		}
		// get the value of the right key
		rv := right.Key(k.Name()).Value()
		// if they match move on
		if lv == rv {
			continue
		}
		u.setRight(left.Name(), k.Name(), rv)
	}
}

func (u *U) Bytes() []byte {
	var b []byte

	for _, section := range u.sections {
		for _, key := range u.sectionKeys[section] {
			_, ok := u.checkFor[section]
			if ok {
				_, ok = u.checkFor[section][key]
				if ok {
					value := u.checkFor[section][key]
					// @FMT(CMP,@dbHOT(INI,SYSTEM.INI,SMS,TARGETMULTI)
					t := fmt.Sprintf("@FMT(CMP,@dbHOT(INI,%s,%s,%s),,'®WIZSET(%s[%s]%s=%s)');\r\n",
						strings.ToUpper(u.file),
						strings.ToUpper(section),
						strings.ToUpper(key),
						strings.ToUpper(u.file),
						strings.ToUpper(section),
						strings.ToUpper(key),
						value,
					)
					b = append(b, []byte(t)...)
					continue
				}
			}
			value := u.D[section][key]
			t := fmt.Sprintf("@WIZSET(%s[%s]%s=%s);\r\n",
				strings.ToUpper(u.file),
				strings.ToUpper(section),
				strings.ToUpper(key),
				value,
			)
			b = append(b, []byte(t)...)
		}
	}
	var pre string
	if !strings.EqualFold(strings.Trim(u.folder, "//"), "OFFICE") {
		pre = fmt.Sprintf("../%s/", u.folder)
	}
	s := fmt.Sprintf("@dbHot(%s%s,SET,%s[*);",
		pre,
		strings.ToUpper(u.file),
		strings.ToUpper(u.file))
	b = append(b, []byte(s)...)
	return b
}

func CreateManual(file string, folder string, section string, key string, value string) *U {

	u := New()
	u.file = file
	u.folder = folder
	u.sections = append(u.sections, section)
	u.sectionKeys[section] = []string{key}
	u.D[section] = make(map[string]string)
	u.D[section][key] = value
	return u
}
