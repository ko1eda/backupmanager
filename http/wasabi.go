package http

import (
	"encoding/json"
	"net/http"

	"github.com/ko1eda/backupmanager/log"
	"github.com/ko1eda/backupmanager/wasabi"
)

// WasabiHandler handles http requests to S3
type wasabiHandler struct {
	s3service  *wasabi.S3Service
	iamservice *wasabi.IAMService
	logger     *log.Logger
	// hasher *Hasher
}

// NewWasabiHandler returns a new handler
func newWasabiHandler(s3s *wasabi.S3Service, iam *wasabi.IAMService, logger *log.Logger) *wasabiHandler {
	return &wasabiHandler{s3s, iam, logger}
}

func (h *wasabiHandler) handleCreateBackupInfrastructure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		hostname := params.Get("host")

		// if there is no hosthname parameter
		if hostname == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := h.iamservice.CreateUser(hostname)
		key := h.iamservice.CreateAccessKeyForUser(user)

		h.s3service.CreateBucket(hostname)
		policy := h.iamservice.CreateLimitedAccessBucketPolicy(hostname)

		h.iamservice.AttachPolicyToUser(policy, user)

		// convert key.AccessKeyId and key.SecretAccessKey to json
		// return response json
		json, err := json.Marshal(&struct {
			AccessKey       string
			SecretAccessKey string
		}{
			*key.AccessKeyId,
			*key.SecretAccessKey,
		})

		if err != nil {
			h.logger.Log("AccessKeyMarshalError: ", err)
			return
		}

		// str := fmt.Sprintf("%s:%s", *key.AccessKeyId, *key.SecretAccessKey)

		// // by writing anything to the reponse body we do not
		// // need a writeheader as go automatically adds 200
		// w.Write([]byte(str))

		// by writing anything to the reponse body we do not
		// need a writeheader as go automatically adds 200
		w.Write(json)
	}
}

// func (h *wasabiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	params := r.URL.Query()

// 	hostname := params.Get("host")

// 	// if there is no hosthname parameter
// 	if hostname == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// user := h.iamservice.CreateUser(hostname)
// 	// _ = h.iamservice.CreateAccessKeyForUser(user)

// 	// h.s3service.CreateBucket(hostname)
// 	// policy := h.iamservice.CreateLimitedAccessBucketPolicy(hostname)

// 	// h.iamservice.AttachPolicyToUser(policy, user)

// 	// convert key.AccessKeyId and key.SecretAccessKey to json
// 	// return response json

// 	w.WriteHeader(http.StatusCreated)
// }
