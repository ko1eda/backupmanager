package http

import (
	"log"
	"net/http"
	"strings"
)

// SecretKeyValidation validates that the incoming request has a valid secret_key paramter
func (s *Server) secretKeyValidation(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := r.URL.Query().Get("secret_key")

		// try and get any nafarious ip addresses from the request
		ff := r.URL.Query().Get("X-Forwarded-For")
		ips := strings.Split(ff, ",")

		if s.Validator.Validate(secret) != true {
			log.Println("InvalidSecretKeyRequest: IPS - ", ips)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.ServeHTTP(w, r)
	}
}
