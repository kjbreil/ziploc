package main

import (
	"io/ioutil"

	"github.com/kjbreil/ziploc/macro"
)

// writeInstall is a hack for now until a macro module is fleshed out
func (o *Option) writeInstall() error {
	s := macro.LineS("@FMT(CMP,@WIZGET(INSTALL_STEP)<3,'®FMT(CHR,26)');")
	s += macro.LineS("@WIZRESET;")
	s += macro.LineS("@EXEC(SQL=REFRESH_MENU);")

	return ioutil.WriteFile(o.makePath("", "INSTALL.SQL"), []byte(s), 0666)
}
