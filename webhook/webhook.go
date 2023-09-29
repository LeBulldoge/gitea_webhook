package webhook

import (
	"log/slog"
	"net/http"

	"github.com/LeBulldoge/gitea_webhook/git"
)

type WebhookHandler struct {
	git *git.Git
}

func NewHandler(gt *git.Git) WebhookHandler {
	return WebhookHandler{
		git: gt,
	}
}

func (m *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		slog.Error("request type should be post")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error("Content-Type should be 'application/json'")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := m.git.Pull()
	if err != nil {
		slog.Error("failed running git", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
