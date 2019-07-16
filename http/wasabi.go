package http

import (
	"net/http"

	"github.com/ko1eda/backupmanager/wasabi"
)

// WasabiHandler handles http requests to S3
type WasabiHandler struct {
	s3service *wasabi.S3Service
	// hasher *Hasher
}

// NewWasabiHandler returns a new handler
func NewWasabiHandler(s *wasabi.S3Service) *WasabiHandler {
	return &WasabiHandler{s}
}

func (h *WasabiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	hostname := params.Get("host")

	// if there is no hosthname parameter
	if hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.s3service.CreateBucket(hostname)

	w.WriteHeader(http.StatusCreated)
}
