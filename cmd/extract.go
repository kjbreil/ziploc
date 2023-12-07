/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/kjbreil/ziploc/extract"
	"github.com/kjbreil/ziploc/option"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract a zip file into a folder and setup config. Pass the zip",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !(len(args) == 1 || len(args) == 2) {
			return fmt.Errorf("accepts 1 or 2 args, received %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		zipPath := filepath.Clean(args[0])
		basePath := filepath.Dir(zipPath)

		if len(args) == 2 {
			basePath = filepath.Clean(args[1])
		}

		doSmsx, _ := cmd.Flags().GetBool("smsx")

		multi, _ := cmd.Flags().GetBool("multi")

		if multi {
			folder, err := os.Stat(filepath.Dir(zipPath))
			if err != nil {
				panic(err)
			}
			if !folder.IsDir() {
				panic("something went wrong, path is not a folder")
			}
			err = filepath.Walk(zipPath, func(path string, info fs.FileInfo, err error) error {
				log.Println(path)

				if filepath.Ext(path) != ".zip" {
					return nil
				}
				if doSmsx {
					createSmsxFolder(path, basePath)
				} else {
					createZiplocFolder(path, basePath)
				}

				return nil
			})
			if err != nil {
				log.Panicln(err)
			}

		} else {
			if filepath.Ext(zipPath) != ".zip" {
				fmt.Printf("file provided was not a zip, did you mean to do a multi (-m)")
				return
			}
			if doSmsx {
				createSmsxFolder(zipPath, basePath)
			} else {
				createZiplocFolder(zipPath, basePath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
	extractCmd.PersistentFlags().BoolP("smsx", "s", false, "create SMSX compatible output directory")
	extractCmd.PersistentFlags().BoolP("multi", "m", false, "given a directory create folder for all zips")
	// extractCmd.PersistentFlags().StringP("out-dir", "o", "", "output to a certain directory, defaults to zip directory")
}

func createSmsxFolder(path, basePath string) {
	var err error
	var newOption *option.Option

	newOption, err = extract.ReadZip(path, basePath, true, nil, true)
	if err != nil {
		log.Panicln(err)
	}

	configName := filepath.Join(newOption.BaseFolder, newOption.Name+".smsx")

	err = newOption.WriteSmsx(configName)
	if err != nil {
		log.Panicln(err)
	}
	// make gitignore

	os.WriteFile(filepath.Join(newOption.BaseFolder, ".gitignore"), []byte("[Rr]elease/%\n[Tt]emp/%\n"), 0666)
}

func createZiplocFolder(path, basePath string) {
	log.Println("doing single zip")
	var err error
	var fromOption *option.Option
	var newOption *option.Option

	// config exists so make o from that
	// if *configLocation != "" {
	// 	fromOption, err = option.ReadConfig(*configLocation)
	// 	if err != nil {
	// 		log.Panicln(err)
	// 	}
	// }

	newOption, err = extract.ReadZip(path, basePath, true, fromOption, false)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(newOption)
	// get the folder that the fromOption is in
	// configLocation := filepath.Dir(*configLocation)
	newOption.SetConfigLocation(filepath.Join(basePath, newOption.Name) + ".json")
	newOption.WriteConfig()
}
