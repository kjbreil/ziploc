package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/iniUpdater"
	"github.com/kjbreil/ziploc/macro"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (o *Option) GetBaseFiles() error {
	baseWithRoot := filepath.Join(o.root, o.BaseFolder)

	info, err := os.Stat(baseWithRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("base directory not found %v", err)
	}
	err = o.Walk(baseWithRoot)
	if err != nil {
		return err
	}
	return nil
}

func (o *Option) GetBaseDSS() error {
	// Create a new DSS for this run
	o.dss = dss.New(o.Name, o.Priority)
	// loop over the paths and add to the DSS
	// for path, info := range o.files {
	// 	switch {
	// 	// ignore root files
	// 	case o.Folder(info.Name()) == "ROOT":
	// 		continue
	// 	// ignore unix hidden files
	// 	case info.Name()[0:1] == ".":
	// 		continue
	// 	}
	// 	if o.iniHasOriginal(info.Name()) {
	// 		continue
	// 	}
	//
	// 	err := o.dss.Add(path)
	// 	if err != nil {
	// 		return fmt.Errorf("error adding %s to dss: %v\n", path, err)
	// 	}
	// }

	return nil
}

// filepath.Join(o.tempSubDir())
func (o *Option) dssWalkTempDir(folder string) error {

	stat, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", folder)
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, eachFile := range files {
		p := filepath.Join(folder, eachFile.Name())
		if eachFile.IsDir() {
			err = o.dssWalkTempDir(p)
			if err != nil {
				return err
			}
			continue
		}
		// log.Printf("adding file: %s\n", p)
		err := o.dss.Add(p)
		if err != nil {
			return fmt.Errorf("error adding %s to dss: %v\n", p, err)
		}
	}
	return nil
}

func (o *Option) FromBase(root string, keepTemp bool) error {
	o.DeleteTemp(root)
	var err error
	// making the temp directory

	tempPath := o.makePath(root, "", "")

	log.Println(tempPath)
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		err = os.MkdirAll(tempPath, 0777)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}
	log.Println("Walking Path", o.BaseFolder)
	// "walk" the BaseFolder for files, adding them to files if the match
	// TODO: REGEX for filename matching (optional?)

	err = o.WriteInstall(root)
	if err != nil {
		return fmt.Errorf("error writing install file: %v\n", err)
	}

	err = o.CopyBase(root)
	if err != nil {
		return err
	}

	// make the DSS here
	err = o.dssWalkTempDir(filepath.Join(o.tempSubDir()))
	if err != nil {
		return err
	}

	// Write the DSS to the temp directory
	err = o.dss.Write(filepath.Join(root, o.tempSubDir(), o.getType(), o.dss.Name))
	if err != nil {
		return err
	}

	err = o.MakeZip(root)
	if err != nil {
		return err
	}

	if !keepTemp {
		o.DeleteTemp(root)
	}

	return nil
}

func (o *Option) CopyBase(root string) error {
	log.Println("-->", root)
	// b := filepath.Join(root, o.BaseFolder)
	for ep, ef := range o.files {

		if ef.Name()[0:1] == "." {
			continue
		}
		var force bool
		f, err := os.Open(ep)
		if err != nil {
			return err
		}
		// folder := ep[len(b) : len(ep)-len(ef.Name())]
		folder := strings.Replace(filepath.Dir(ep), o.BaseFolder, "", 1)
		if len(folder) > 0 && folder[0] == os.PathSeparator {
			folder = folder[1:]
		}

		//folder := filepath.Base(filepath.Dir(ep))
		// log.Println(filepath.Base(o.BaseFolder), folder)
		if filepath.Base(o.BaseFolder) == folder {
			folder = ""
		}
		dest := o.makePath(root, folder, ef.Name())

		// is the file a part of the ini map, maybe do something then
		noExt := strings.ToUpper(ef.Name()[:len(ef.Name())-len(filepath.Ext(ef.Name()))])
		ini, ok := o.IniMaps[noExt]
		if ok {
			if ini.Original != nil {
				f, err = iniUpdater.File(*ini.Original, *ini.Base, folder)
				if err != nil {
					return err
				}
				dest = o.makePath(root, "INBOX", noExt+"_INI.SQL")
				force = true
			}
		}

		if strings.EqualFold(filepath.Ext(ef.Name()), ".sql") {
			force = true
		}

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

		err = macro.Correct(out, f, force)
		if err != nil {
			return err
		}

		_ = f.Close()
		_ = out.Close()
	}
	if o.Version != "" {
		return o.makeVersionIni(root)
	}

	return nil
}
