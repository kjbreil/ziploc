package option

import (
	"fmt"
	"path/filepath"
	"strings"
)

type IniMap struct {
	Name     string  `json:"name"`
	Base     *string `json:"base,omitempty"`
	Original *string `json:"original,omitempty"`
	Instance *string `json:"instance,omitempty"`
}

func (o *Option) FindINI() error {
	if o.dss == nil {
		return fmt.Errorf("dss has not been generated yet")
	}

	for p, f := range o.files {
		if filepath.Ext(p) == ".ini" {
			if o.IniMaps == nil {
				o.IniMaps = make(map[string]IniMap)
			}
			_, ok := o.IniMaps[f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))]]
			if ok {
				continue
			}

			i := IniMap{Name: f.Name()}

			if o.InstanceDir != nil {
				i.Instance = o.GetIniInstancePath(f.Name())
			} else {
				i.Base = &p
			}

			o.IniMaps[f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))]] = i
		}
	}

	return nil
}

func (o *Option) GetIniInstancePath(i string) *string {
	for p, f := range o.instanceFiles {
		if strings.EqualFold(i, f.Name()) {
			return &p
		}
	}
	return nil
}
