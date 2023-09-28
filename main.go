package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/LeBulldoge/gitea_webhook/webhook"
)

func main() {
	flag.Parse()

	slog.Info("starting webhook server")

	http.HandleFunc("/webhook", webhook.Handle)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else if err != nil {
		slog.Error("error starting server", "err", err)
		os.Exit(1)
	}
}
