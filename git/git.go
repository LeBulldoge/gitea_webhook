package git

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var (
	repoDir    = flag.String("repo", "", "Target repo directory")
	pemKey     = flag.String("pem", "", "Path to pem key for ssh auth")
	passphrase = flag.String("pass", "", "Passphrase for private key")
)

func Pull() error {
	slog.Info("performing git pull")

	repo, err := git.PlainOpen(*repoDir)
	if err != nil {
		return fmt.Errorf("error opening git repo: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting git worktree: %w", err)
	}

	opt := git.PullOptions{
		RemoteName: "origin",
	}

	if *pemKey != "" {
		keys, err := ssh.NewPublicKeysFromFile("git", *pemKey, *passphrase)
		if err != nil {
			return fmt.Errorf("error creating public keys: %w", err)
		}
		opt.Auth = keys
	}

	err = wt.Pull(&opt)
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		slog.Info("target repository already up to date")
		return nil
	} else if err != nil {
		return fmt.Errorf("error doing git pull: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return fmt.Errorf("error getting git head ref: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return fmt.Errorf("error getting commit from ref: %w", err)
	}

	slog.Info("pull finished", "latest commit", commit)

	return nil
}
