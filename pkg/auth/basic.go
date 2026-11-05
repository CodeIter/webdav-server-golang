// Package auth provides HTTP authentication middleware.
package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// BasicAuth returns middleware that requires the given HTTP Basic credentials.
// Both values are hashed before comparison so the check runs in constant time
// and does not leak the length of the expected credentials.
func BasicAuth(realm, username, password string) func(http.Handler) http.Handler {
	wantUser := sha256.Sum256([]byte(username))
	wantPass := sha256.Sum256([]byte(password))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || !matches(user, wantUser) || !matches(pass, wantPass) {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`", charset="UTF-8"`)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func matches(got string, want [sha256.Size]byte) bool {
	sum := sha256.Sum256([]byte(got))
	return subtle.ConstantTimeCompare(sum[:], want[:]) == 1
}
