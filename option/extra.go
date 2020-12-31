package option

import (
	"io/ioutil"

	"github.com/kjbreil/ziploc/macro"
)

// WriteInstall is a hack for now until a macro module is fleshed out
func (o *Option) WriteInstall(root string) error {
	s := macro.LineS("@FMT(CMP,@WIZGET(INSTALL_STEP)<3,'®FMT(CHR,26)');")
	s += macro.LineS("@WIZRESET;")
	s += macro.LineS("@EXEC(SQL=REFRESH_MENU);")

	return ioutil.WriteFile(o.makePath(root, "", "INSTALL.SQL"), []byte(s), 0666)
}
