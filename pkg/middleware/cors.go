package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			headers.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,PATCH")
			headers.Set("Access-Control-Allow-Headers", "authorization,content-type,content-length")
			headers.Set("Access-Control-Max-Age", "86400")
			return
		}
		next.ServeHTTP(w, r)
	})
}
