package create

import (
	"fmt"

	"github.com/kjbreil/ziploc/option"
)

type Single struct {
	JustWalk         bool
	KeepTemp         bool
	WriteConfig      bool
	PullFromInstance bool
	Option           *option.Option
}

func (s *Single) Run() error {
	var err error

	if s.Option.Git != nil {
		err = s.Option.DoGitRepo()
		if err != nil {
			return fmt.Errorf("git repo read error: %w", err)
		}
	}

	// if instanceDir is set walk it
	if s.Option.InstanceDir != nil {
		err = s.Option.WalkInstance(*s.Option.InstanceDir)
		if err != nil {
			return err
		}

		if s.JustWalk {
			// copy from the instance dir if just walking
			err = copyInstanceDir(s.Option)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// get the files
	err = s.Option.GetBaseFiles()
	if err != nil {
		return err
	}

	// make dss object for base files
	err = s.Option.GetBaseDSS()
	if err != nil {
		return err
	}

	// find any INI's and add them to the map
	err = s.Option.FindINI()
	if err != nil {
		return err
	}

	err = copyInstanceDir(s.Option)
	if err != nil {
		return err
	}

	err = s.Option.FromBase("", s.KeepTemp)
	if err != nil {
		return err
	}

	if !s.WriteConfig {
		return nil
	}

	return s.Option.WriteConfig()
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
