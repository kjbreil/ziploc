package main

import (
	"flag"
	"github.com/kjbreil/ziploc/extract"
	"github.com/kjbreil/ziploc/option"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	configLocation = flag.String("c", "", "Config file")
	makeTemplate   = flag.Bool("template", false, "output a template config file and exit")

	newConfig = flag.Bool("nc", false, "New config file to build from an original config and a new basepath")

	fromZip  = flag.String("zip", "", "read from a zip and build config based on said zip, this is the zip file ot extract")
	basePath = flag.String("base", "base", "base path to which extraction of the zip happens, if a config is provided this is ignored")

	outDir = flag.String("out", "", "override zip output directory defined in the json")
)

func main() {
	flag.Parse()

	if *makeTemplate {
		option.ConfigTemplate(true)
		return
	}

	// from zip argument passed
	if *fromZip != "" {
		doFromZip()
		return
	}

	if *newConfig && *configLocation != "" {
		copyConfig()
		return
	}

	// // no configuration passed so run GUI
	if *configLocation == "" {
		// gui.OpenGui()
		log.Println("no config specified")
		return
	}
	cl := filepath.Clean(*configLocation)

	if filepath.Base(cl) == "*" {
		folder, err := os.Stat(filepath.Dir(cl))
		if err != nil {
			panic(err)
		}
		if !folder.IsDir() {
			panic("something went wrong, not a folder")
		}
		log.Println(filepath.Dir(cl))
		filepath.Walk(filepath.Dir(cl), func(path string, info fs.FileInfo, err error) error {

			if filepath.Ext(path) != ".json" {
				return nil
			}
			doSingleConfig(path)

			return nil
		})
	} else {
		doSingleConfig(*configLocation)
	}
}

func doSingleConfig(configPath string) {
	o, err := option.ReadConfig(configPath)
	if err != nil {
		log.Panic(err)
	}

	if *outDir != "" {
		o.ZipOut = *outDir
	}

	if o.Git != nil {
		err = o.DoGitRepo()
		if err != nil {
			log.Panicln(err)
		}
	}

	// get the files
	err = o.GetBaseFiles()
	if err != nil {
		log.Panic(err)
	}
	// make dss object for base files
	err = o.GetBaseDSS()
	if err != nil {
		log.Panic(err)
	}
	// if instanceDir is set walk it
	if o.InstanceDir != nil {
		err = o.WalkInstance(*o.InstanceDir)
		if err != nil {
			log.Panic(err)
		}
	}

	// find any INI's and add them to the map
	err = o.FindINI()
	if err != nil {
		log.Panic(err)
	}

	if o.InstanceDir != nil {
		err = o.FromInstance()
		if err != nil {
			log.Panic(err)
		}
	}

	err = o.FromBase("")
	if err != nil {
		log.Panic(err)
	}

	o.WriteConfig(configPath)
}

func doFromZip() {
	fz := filepath.Clean(*fromZip)

	if filepath.Base(fz) == "*" {
		log.Println("making from zips")

		folder, err := os.Stat(filepath.Dir(fz))
		if err != nil {
			panic(err)
		}
		if !folder.IsDir() {
			panic("something went wrong, not a folder")
		}
		filepath.Walk(filepath.Dir(fz), func(path string, info fs.FileInfo, err error) error {
			log.Println(path)

			if filepath.Ext(path) != ".zip" {
				return nil
			}
			log.Println("doing zip", path)
			doSingleZip(path)

			return nil
		})

	} else {
		log.Println("making from zip")

		doSingleZip(*fromZip)
	}
	return
}

func doSingleZip(path string) {
	log.Println("doing single zip")
	var err error
	var fromOption *option.Option
	var newOption *option.Option

	// config exists so make o from that
	if *configLocation != "" {
		fromOption, err = option.ReadConfig(*configLocation)
		if err != nil {
			log.Panicln(err)
		}
	}
	newOption, err = extract.ReadZip(path, *basePath, true, fromOption)
	if err != nil {
		log.Panicln(err)
	}
	// get the folder that the fromOption is in
	configLocation := filepath.Dir(*configLocation)

	newOption.WriteConfig(filepath.Join(configLocation, newOption.Name) + ".json")
}

func copyConfig() {
	cleanBasePath := filepath.Clean(*basePath)
	cleanConfigLocation := filepath.Clean(*configLocation)

	log.Println(cleanBasePath)

	originalOption, err := option.ReadConfig(cleanConfigLocation)
	if err != nil {
		log.Panic(err)
	}

	newOptionName := filepath.Base(cleanBasePath)

	newOptionBaseFolder := filepath.Join("..", cleanBasePath)

	// log.Println(findMatchingTopLevelDir(originalOption.BaseFolder, cleanBasePath, cleanConfigLocation))

	newOption := option.Option{
		Name:        newOptionName,
		Priority:    originalOption.Priority,
		IsOption:    originalOption.IsOption,
		Include:     originalOption.Include,
		Exclude:     originalOption.Exclude,
		Ignore:      originalOption.Ignore,
		BaseFolder:  newOptionBaseFolder,
		TempDir:     originalOption.TempDir,
		ZipOut:      originalOption.ZipOut,
		Instance:    originalOption.Instance,
		InstanceDir: originalOption.InstanceDir,
		Git:         originalOption.Git,
		IniMaps:     originalOption.IniMaps,
	}
	cleanORiginalBaseFolder := filepath.Clean(originalOption.BaseFolder)

	for k, v := range newOption.IniMaps {
		trimmed := strings.TrimPrefix(*newOption.IniMaps[k].Base, strings.TrimPrefix(cleanORiginalBaseFolder, "../"))

		newIniBase := filepath.Join(cleanBasePath, trimmed)
		v.Base = &newIniBase
		newOption.IniMaps[k] = v
	}

	outFileName := filepath.Join(filepath.Dir(cleanConfigLocation), newOptionName) + ".json"

	newOption.WriteConfig(outFileName)
	//
	// log.Println(filepath.Base(cleanBasePath))
	// log.Println(filepath.Dir(cleanConfigLocation))

}

// func findMatchingTopLevelDir(originalBase, newBase, configLocation string) string {
// 	originalBase, newBase = filepath.Clean(originalBase), filepath.Clean(newBase)
//
// 	configSplit := strings.Split(configLocation, "/")
// 	newBaseSplit := strings.Split(newBase, "/")
//
// 	log.Println(new)
// 	log.Println(configLocation)
// 	return ""
// }
