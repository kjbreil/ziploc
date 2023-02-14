package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/iniUpdater"
	"github.com/kjbreil/ziploc/macro"
	"os"
	"path/filepath"
)

func (o *Option) makeVersionIni(root string) error {
	f, _ := iniUpdater.ManualFile("NCBP_VERSION.INI", "OFFICE", "VERSION", o.Name, o.Version)
	dest := o.makePath(root, "INBOX", "NCBP_VERSION_INI.SQL")
	if _, err := os.Stat(filepath.Dir(dest)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(dest), 0777)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	err = macro.Correct(out, f, true)
	if err != nil {
		return err
	}

	_ = f.Close()
	_ = out.Close()

	return nil
}
