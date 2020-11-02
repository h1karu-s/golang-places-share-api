package middlewares

import (
	"net/http"
)

//AddAccessHeader .
func AddAccessHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			next.ServeHTTP(w, r)
	})
}