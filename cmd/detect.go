/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kjbreil/ziploc/option"

	"github.com/spf13/cobra"
)

// singleCmd represents the single command
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "detect custom files",
	Long:  `.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var o *option.Option
		var err error

		configLoc := filepath.Clean(args[0])

		o, err = option.ReadConfig(configLoc, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("error reading config file %s", configLoc), "err", err)

			os.Exit(1)
		}

		err = o.LookupCustomFiles()
		if err != nil {
			logger.Error("error looking up custom files", "err", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(detectCmd)

}
