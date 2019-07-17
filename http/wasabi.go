package http

import (
	"net/http"

	"github.com/ko1eda/backupmanager/wasabi"
)

// WasabiHandler handles http requests to S3
type WasabiHandler struct {
	s3service  *wasabi.S3Service
	iamservice *wasabi.IAMService
	// hasher *Hasher
}

// NewWasabiHandler returns a new handler
func NewWasabiHandler(s3s *wasabi.S3Service, iam *wasabi.IAMService) *WasabiHandler {
	return &WasabiHandler{s3s, iam}
}

func (h *WasabiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	hostname := params.Get("host")

	// if there is no hosthname parameter
	if hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := h.iamservice.CreateUser(hostname)
	_ = h.iamservice.CreateAccessKeyForUser(user)

	h.s3service.CreateBucket(hostname)
	policy := h.iamservice.CreateLimitedAccessBucketPolicy(hostname)

	h.iamservice.AttachPolicyToUser(policy, user)

	// convert key.AccessKeyId and key.SecretAccessKey to json
	// return response json

	w.WriteHeader(http.StatusCreated)
}
