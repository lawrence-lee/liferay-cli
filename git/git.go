package git

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/viper"
	"liferay.com/lcectl/constants"
)

func init() {
	dirname, err := os.UserHomeDir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault(constants.Const.RepoDir, filepath.ToSlash(path.Join(dirname, ".lcectl", "sources", "localdev")))
	viper.SetDefault(constants.Const.RepoRemote, "https://github.com/gamerson/lxc-localdev")
	viper.SetDefault(constants.Const.RepoBranch, "master")
	viper.SetDefault(constants.Const.RepoSync, true)
}

func SyncGit(verbose bool) {
	if repoSync := viper.GetBool(constants.Const.RepoSync); !repoSync {
		return
	}

	var s *spinner.Spinner

	if !verbose {
		s = spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Color("green")
		s.Suffix = " Synchronizing 'localdev' sources..."
		s.FinalMSG = fmt.Sprintf("\u2705 Synced 'localdev' sources.\n")
		s.Start()
		defer s.Stop()
	}

	repoDir := viper.GetString(constants.Const.RepoDir)
	repo, err := git.PlainOpen(repoDir)

	cloned := false

	if err != nil {
		os.MkdirAll(repoDir, os.ModePerm)

		repo, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			Depth:        1,
			RemoteName:   "origin",
			SingleBranch: true,
			URL:          viper.GetString(constants.Const.RepoRemote),
		})

		if err != nil {
			log.Fatal("Clone error: ", err)
		} else {
			fmt.Printf("Cloned 'localdev' sources to %s\n", repoDir)
		}

		cloned = true
	}

	if !cloned {
		worktree, err := repo.Worktree()

		if err != nil {
			log.Fatal("worktree error: ", err)
		}

		if err = worktree.Pull(&git.PullOptions{
			Force: true,
			ReferenceName: plumbing.ReferenceName(
				"refs/heads/" + viper.GetString(constants.Const.RepoBranch)),
			RemoteName:   "origin",
			SingleBranch: true,
		}); err != nil {

			if err == git.NoErrAlreadyUpToDate || err == transport.ErrEmptyUploadPackRequest {
				if s != nil {
					s.FinalMSG = fmt.Sprintf("\u2705 'localdev' sources %s.\n", git.NoErrAlreadyUpToDate)
				} else {
					fmt.Printf("'localdev' sources %s.\n", git.NoErrAlreadyUpToDate)
				}

				return
			}

			if s != nil {
				s.FinalMSG = fmt.Sprintf("\u2718 'localdev' sources error %s.\n", err)
				os.Exit(1)
			} else {
				log.Fatalf("'localdev' sources error %s.\n", err)
			}
		}

		if s != nil {
			s.FinalMSG = fmt.Sprintf("\u2705 'localdev' sources updated.\n")
		} else {
			fmt.Printf("'localdev' sources updated.\n")
		}
	}
}
