package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"log"
	"os"
)

func (o *Option) localName() string {
	return fmt.Sprintf("./%s", o.safeName())
}

func (o *Option) doGitRepo() error {

	auth := &http.BasicAuth{
		Username: o.Git.Username,
		Password: o.Git.AuthToken,
	}

	_, err := os.Stat(o.localName())
	// create the repository object
	var r *git.Repository

	// if the stat has an error that means the folder does not exist
	if err != nil {
		r, err = git.PlainClone(o.localName(), false, &git.CloneOptions{
			URL:      o.Git.URL,
			Progress: os.Stdout,
			Auth:     auth,
		})
		if err != nil {
			log.Println("here", err)
		}
	} else {
		r, err = git.PlainOpen(o.localName())
		if err != nil {
			log.Println("there", err)
		}
	}

	if r == nil {
		return fmt.Errorf("problem loading repository")
	}

	// setup the worktree
	w, err := r.Worktree()
	if err != nil {
		log.Panicln(err)
	}

	// Pull the repository

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
		Force:      true,
	})
	if err != nil {
		log.Println(err)
	}

	// fetch the branches
	err = r.Fetch(&git.FetchOptions{
		Auth:     auth,
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println(err)
	}

	branch := "master"
	if o.Git.Branch != "" {
		branch = o.Git.Branch
	}

	ref := plumbing.NewBranchReferenceName(branch)

	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref,
		Force:  true,
	})
	if err != nil {
		fmt.Println(err)
	}

	o.BaseFolder = o.safeName()

	return nil
}
