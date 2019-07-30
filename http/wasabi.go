package http

import (
	"net/http"
)

// handleCreateBackupInfrastructure creates an IAM user, an S3 bucket
// and then maps a restrictive policy to that user to allow limited access to the bucket.
// If there is an eror a log will be written and an email will be sent.
func (s *Server) handleCreateBackupInfrastructure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hostname := r.URL.Query().Get("host")
		if hostname == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// err := s.Mailer.Send("koleda.christopher@gmail.com", "cd-backup-generator@securedatatransit.com", "Test Message")

		// if err != nil {
		// 	log.Println("MailSendError: ", err)
		// }

		// user, err := s.IAMService.CreateUser(hostname)
		// if err != nil {
		// 	log.Println("CreateUserError: ", err)

		// 	// send email

		// 	w.WriteHeader(http.StatusFailedDependency)
		// 	return
		// }

		// key, err := s.IAMService.CreateAccessKeyForUser(user)
		// if err != nil {
		// 	log.Println("CreateKeyError: ", err)

		// 	// send email

		// 	w.WriteHeader(http.StatusFailedDependency)
		// 	return
		// }

		// bucket, err := s.S3Service.CreateBucket(hostname)
		// if err != nil {
		// 	log.Println("CreateBucketError: ", err)

		// 	//send email

		// 	w.WriteHeader(http.StatusFailedDependency)
		// 	return
		// }

		// policy, err := s.IAMService.CreateLimitedAccessBucketPolicy(bucket)
		// if err != nil {
		// 	log.Println("CreatePolicyError: ", err)

		// 	// send email

		// 	w.WriteHeader(http.StatusFailedDependency)
		// 	return
		// }

		// err = s.IAMService.AttachPolicyToUser(policy, user)
		// if err != nil {
		// 	log.Println("AttachPolicyError: ", err)

		// 	// send email

		// 	w.WriteHeader(http.StatusFailedDependency)
		// 	return
		// }

		// // by writing anything to the reponse body we do not
		// // need a writeheader as go automatically adds 200
		// str := fmt.Sprintf("%s,%s", *key.AccessKeyId, *key.SecretAccessKey)
		w.Write([]byte("hello"))
	}
}
