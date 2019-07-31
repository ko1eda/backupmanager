package http

import (
	"log"
	"net/http"
)

// SecretKeyValidation validates that the incoming request has a valid secret_key paramter
func (s *Server) secretKeyValidation(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := r.URL.Query().Get("secret_key")

		if s.Validator.Validate(secret) != true {
			log.Println("InvalidSecretKeyRequest")

			// err := s.Mailer.Send(
			// 	"",
			// 	"support@creatingdigital.com",
			// 	"InvalidSecretKeyRequest",
			// 	"Someone has tried to access go backup generator api without a valid secret key",
			// )

			// if err != nil {
			// 	log.Println("MailSendError: ", err)
			// }

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.ServeHTTP(w, r)
	}
}
