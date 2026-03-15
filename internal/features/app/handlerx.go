package app

import (
	"io"
	"log"
	"net/http"
)

func (h Handler) receive(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	log.Printf("GET /webhook query=%s", r.URL.Query())

	log.Printf("GET /webhook mode=%q token=%q challenge=%q", mode, token, challenge)

	if mode == "subscribe" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(challenge))
		return
	}

	http.Error(w, "forbidden", http.StatusForbidden)
}

func (h Handler) verify(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("POST /webhook headers=%v", r.Header)
	log.Printf("POST /webhook body=%s", string(body))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("EVENT_RECEIVED"))
}
