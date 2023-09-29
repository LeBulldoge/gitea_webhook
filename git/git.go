package git

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Git struct {
	repoDir    string
	pemKey     string
	passphrase string

	repo *git.Repository
}

func New(repoDir string, pemKey string, passphrase string) (*Git, error) {
	res := &Git{
		repoDir:    repoDir,
		pemKey:     pemKey,
		passphrase: passphrase,
	}

	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return nil, fmt.Errorf("error opening git repo: %w", err)
	}

	res.repo = repo

	return res, nil
}

func (m *Git) Pull() error {
	slog.Info("performing git pull")

	repo, err := git.PlainOpen(m.repoDir)
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

	if m.pemKey != "" {
		keys, err := ssh.NewPublicKeysFromFile("git", m.pemKey, m.passphrase)
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
