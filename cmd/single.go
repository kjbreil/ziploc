/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kjbreil/ziploc/create"
	"github.com/kjbreil/ziploc/option"

	"github.com/spf13/cobra"
)

// singleCmd represents the single command
var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "create a option/sample from a single config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var o *option.Option
		var err error

		configLoc := filepath.Clean(args[0])

		doSmsx, _ := cmd.Flags().GetBool("smsx")
		if doSmsx {
			o, err = option.ReadSmsxConfig(configLoc)
		} else {
			o, err = option.ReadConfig(configLoc, logger)
		}
		if err != nil {
			logger.Error(fmt.Sprintf("error reading config file %s", configLoc), "err", err)

			os.Exit(1)
		}

		justWalk, _ := cmd.Flags().GetBool("just-walk")
		keepTemp, _ := cmd.Flags().GetBool("keep-temp")
		writeConfig, _ := cmd.Flags().GetBool("write-config")
		outDir, _ := cmd.Flags().GetString("out-dir")
		pullFromInstance, _ := cmd.Flags().GetBool("pull-instance")
		ignoreInstance, _ := cmd.Flags().GetBool("ignore-instance")

		if ignoreInstance {
			o.InstanceDir = nil
		}

		if outDir != "" {
			o.ZipOut = outDir
		}

		s := create.Single{
			JustWalk:         justWalk,
			KeepTemp:         keepTemp,
			WriteConfig:      writeConfig,
			PullFromInstance: pullFromInstance,
			Option:           o,
			DetectChanges:    true,
		}

		err = s.Run()
		if err != nil {
			logger.Error(fmt.Sprintf("error running %s", o.Name), "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)

	singleCmd.PersistentFlags().StringP("out-dir", "o", "", "override zip output directory defined in the json")
	singleCmd.PersistentFlags().StringP("only-custom", "c", "", "only copy files that are detected as custom")
	singleCmd.PersistentFlags().StringP("pull-instance", "i", "", "pull from instance")
	singleCmd.PersistentFlags().BoolP("just-walk", "j", false, "just walk the instance, copying files over without creating a zip")
	singleCmd.PersistentFlags().BoolP("ignore-instance", "u", false, "skip using anything to do with instance")
	singleCmd.PersistentFlags().BoolP("keep-temp", "k", false, "keep temp files")
	singleCmd.PersistentFlags().BoolP("smsx", "s", false, "use SMSX config format")
	singleCmd.PersistentFlags().BoolP("write-config", "w", false, "output config file after generation")
}
