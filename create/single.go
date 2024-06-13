package create

import (
	"fmt"

	"github.com/kjbreil/ziploc/option"
)

func Single(o *option.Option, justWalk, keepTemp, writeConfig bool) error {
	var err error

	if o.Git != nil {
		err = o.DoGitRepo()
		if err != nil {
			return fmt.Errorf("git repo read error: %w", err)
		}
	}

	// if instanceDir is set walk it
	if o.InstanceDir != nil {
		err = o.WalkInstance(*o.InstanceDir)
		if err != nil {
			return err
		}

		if justWalk {
			// copy from the instance dir if just walking
			err = copyInstanceDir(o)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// get the files
	err = o.GetBaseFiles()
	if err != nil {
		return err
	}

	// make dss object for base files
	err = o.GetBaseDSS()
	if err != nil {
		return err
	}

	// find any INI's and add them to the map
	err = o.FindINI()
	if err != nil {
		return err
	}

	err = copyInstanceDir(o)
	if err != nil {
		return err
	}

	err = o.FromBase("", keepTemp)
	if err != nil {
		return err
	}

	if !writeConfig {
		return nil
	}

	return o.WriteConfig()
}

func copyInstanceDir(o *option.Option) error {
	if o.InstanceDir != nil {
		err := o.FromInstance()
		if err != nil {
			return err
		}
	}
	return nil
}
