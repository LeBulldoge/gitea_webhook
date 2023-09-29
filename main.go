package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/LeBulldoge/gitea_webhook/git"
	"github.com/LeBulldoge/gitea_webhook/webhook"
)

var (
	repoDir    = flag.String("repo", "", "Target repo directory")
	pemKey     = flag.String("pem", "", "Path to pem key for ssh auth")
	passphrase = flag.String("pass", "", "Passphrase for private key")
)

func main() {
	flag.Parse()

	slog.Info("starting webhook server")

	gt, err := git.New(*repoDir, *pemKey, *passphrase)
	if err != nil {
		slog.Error("failure initializing git repo", "err", err)
		return
	}

	handler := webhook.NewHandler(gt)
	http.Handle("/webhook", &handler)

	err = http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else if err != nil {
		slog.Error("error starting server", "err", err)
		os.Exit(1)
	}
}
