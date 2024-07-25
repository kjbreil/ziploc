package option

import (
	"fmt"
	"github.com/kjbreil/ziploc/dss"
	"github.com/kjbreil/ziploc/iniUpdater"
	"github.com/kjbreil/ziploc/macro"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (o *Option) GetBaseFiles() error {
	baseWithRoot := filepath.Join(o.root, o.BaseFolder)

	if err := statBase(baseWithRoot); err != nil {
		return err
	}

	if err := o.Walk(baseWithRoot); err != nil {
		return err
	}
	return nil
}

func statBase(baseWithRoot string) error {
	info, err := os.Stat(baseWithRoot)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("base directory not found %v", err)
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
	// 	err := o.dss.add(path)
	// 	if err != nil {
	// 		return fmt.Errorf("error adding %s to dss: %v\n", path, err)
	// 	}
	// }

	return nil
}

func (o *Option) FromBase(root string, keepTemp bool) error {
	if root == "" && o.root != "" {
		root = o.root
	}
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
	err = o.dss.WalkDir(filepath.Join(o.root, o.tempSubDir()))
	if err != nil {
		return err
	}

	// Write the DSS to the temp directory
	err = o.dss.Write(filepath.Join(root, o.tempSubDir(), o.getType(), o.Name))
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
		// remove anything to the left of basefolder in each path
		_, folder, _ := strings.Cut(filepath.Dir(ep), o.BaseFolder)

		// folder := strings.Replace(filepath.Dir(ep), o.BaseFolder, "", 1)
		if len(folder) > 0 && folder[0] == os.PathSeparator {
			folder = folder[1:]
		}

		// folder := filepath.Base(filepath.Dir(ep))
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
