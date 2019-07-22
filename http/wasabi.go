package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ko1eda/backupmanager/wasabi"
)

// WasabiHandler handles http requests to S3
type wasabiHandler struct {
	s3service  *wasabi.S3Service
	iamservice *wasabi.IAMService
	//mailer
	// hasher *Hasher
}

// NewWasabiHandler returns a new handler
func newWasabiHandler(s3s *wasabi.S3Service, iam *wasabi.IAMService) *wasabiHandler {
	return &wasabiHandler{s3s, iam}
}

func (h *wasabiHandler) handleCreateBackupInfrastructure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Query().Get("host")
		if hostname == "" {
			// if there is no hostname parameter the req is malformed
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := h.iamservice.CreateUser(hostname)
		if err != nil {
			log.Println("CreateUserError: ", err)

			// send email

			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		key, err := h.iamservice.CreateAccessKeyForUser(user)
		if err != nil {
			log.Println("CreateKeyError: ", err)

			// send email

			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		bucket, err := h.s3service.CreateBucket(hostname)
		if err != nil {
			log.Println("CreateBucketError: ", err)

			//send email

			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		policy, err := h.iamservice.CreateLimitedAccessBucketPolicy(bucket)
		if err != nil {
			log.Println("CreatePolicyError: ", err)

			// send email

			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		err = h.iamservice.AttachPolicyToUser(policy, user)
		if err != nil {
			log.Println("AttachPolicyError: ", err)

			// send email

			w.WriteHeader(http.StatusFailedDependency)
			return
		}

		// by writing anything to the reponse body we do not
		// need a writeheader as go automatically adds 200
		str := fmt.Sprintf("%s,%s", *key.AccessKeyId, *key.SecretAccessKey)
		w.Write([]byte(str))
	}
}
