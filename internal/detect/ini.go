package detect

import (
	"github.com/kjbreil/ziploc/internal/iniUpdater"
	"path/filepath"
	"strings"
)

type IniMap struct {
	Base      string   `json:"base"`
	Originals []string `json:"original"`
}

func (d *Detect) GetINIs() error {
	err := d.Detect()
	if err != nil {
		return err
	}
	iniMaps := make(map[string]IniMap)

	for _, e := range d.customDSS.Data {
		if strings.EqualFold(filepath.Ext(e.Script), ".INI") {

			entries := d.originalDSS.Get(e.Script)

			m := IniMap{
				Base:      BuildFilePath(filepath.Dir(e.Author), e.Destination, e.Script, false),
				Originals: make([]string, 0, len(entries)*2),
			}
			for _, ee := range entries {
				m.Originals = append(m.Originals, BuildFilePath(ee.Author, ee.Destination, ee.Script, false))
			}
			for _, ee := range entries {
				m.Originals = append(m.Originals, BuildFilePath(ee.Author, ee.Destination, ee.Script, true))
			}
			iniMaps[e.Script] = m

		}
	}
	d.inis = make(map[string][]*iniUpdater.U)
	for k, v := range iniMaps {

		for _, o := range v.Originals {
			u, err := iniUpdater.CompareTwo(v.Base, o, "OFFICE")
			if err != nil {
				continue
			}
			if !u.Empty() {
				d.inis[k] = append(d.inis[k], u)
			}
		}
		d.customDSS.Delete(k)
	}

	return nil
}

func BuildFilePath(storemanDir, destination, fileName string, firstLoad bool) string {
	lowerD := strings.ToLower(destination)
	if _, ok := ValidDirsInStoreman[lowerD]; ok {
		if firstLoad {
			return filepath.Join(storemanDir, "FirstLoad", destination, fileName)
		}
		return filepath.Join(storemanDir, destination, fileName)

	}
	if _, ok := ValidDirsInOffice[lowerD]; ok {
		if firstLoad {
			return filepath.Join(storemanDir, "FirstLoad", destination, fileName)
		}
		return filepath.Join(storemanDir, "Office", destination, fileName)
	}
	return ""
}

// filepath.Join("FirstLoad", ee.Destination)
