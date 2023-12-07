/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kjbreil/ziploc/create"
	"github.com/kjbreil/ziploc/option"
	"github.com/spf13/cobra"
)

// multiCmd represents the multi command
var multiCmd = &cobra.Command{
	Use:   "multi",
	Short: "Generate multiple options/samples",
	Long: `multi allows you to generate multiple options/samples 
using a single run. By default it looks in a given directory for
config files however it can be set to recursively look in 
subdirectories for valid config files.
For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topDir := filepath.Dir(filepath.Clean(args[0]))

		folder, err := os.Stat(topDir)
		if err != nil {
			fmt.Printf("directory read error: %v", err)
			return
		}
		if !folder.IsDir() {
			fmt.Printf("something went wrong, not a folder")
			return
		}
		recursive, _ := cmd.Flags().GetBool("recursive")
		doSmsx, _ := cmd.Flags().GetBool("smsx")

		options := findConfigsInDirectory(topDir, doSmsx, recursive)

		justWalk, _ := cmd.Flags().GetBool("just-walk")
		keepTemp, _ := cmd.Flags().GetBool("keep-temp")
		writeConfig, _ := cmd.Flags().GetBool("write-config")
		version, _ := cmd.Flags().GetString("set-version")

		for _, o := range options {
			if version != "" {
				o.Version = version
			}
			err = create.Single(o, justWalk, keepTemp, writeConfig)
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(multiCmd)

	multiCmd.PersistentFlags().BoolP("recursive", "r", false, "recursively look for valid config files")

	multiCmd.PersistentFlags().StringP("out-dir", "o", "", "override zip output directory defined in the json")
	multiCmd.PersistentFlags().BoolP("just-walk", "j", false, "just walk the instance, copying files over without creating a zip")
	multiCmd.PersistentFlags().BoolP("keep-temp", "k", false, "keep temp files")
	multiCmd.PersistentFlags().BoolP("smsx", "s", false, "use SMSX config format")
	multiCmd.PersistentFlags().BoolP("write-config", "w", false, "output config file after generation")
	multiCmd.PersistentFlags().StringP("set-version", "v", "", "force a specific version tag")

}

func findConfigsInDirectory(dir string, smsx, recursive bool) map[string]*option.Option {
	options := make(map[string]*option.Option)

	// TODO: Capture errors and return a list of them
	_ = filepath.Walk(filepath.Dir(dir), func(path string, info fs.FileInfo, err error) error {

		// if info.IsDir() && recursive {
		// 	add := findConfigsInDirectory(path, smsx, recursive)
		// 	for k, v := range add {
		// 		options[k] = v
		// 	}
		// }

		baseDir := filepath.Dir(path)
		if dir != baseDir && !recursive {
			return nil
		}

		if !smsx && filepath.Ext(path) == ".json" {
			o, _ := option.ReadConfig(path)
			options[path] = o
		}

		if smsx && filepath.Ext(path) == ".smsx" {
			o, _ := option.ReadSmsxConfig(path)
			options[path] = o
		}

		return nil
	})

	return options
}
