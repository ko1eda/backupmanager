package http

import "net/http"

// WithSecretKeyValidation validates that the incoming request has a valid secret_key paramter
func (s *Server) WithSecretKeyValidation(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := r.URL.Query().Get("secret_key")

		if s.Validator.Validate(secret) != true {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h.ServeHTTP(w, r)
	}
}
