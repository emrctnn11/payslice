package middleware

import "net/http"

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow connection from anywhere
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// allow specific http methods
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTION, PUT, DELETE")
		// Allow specific headers (Content-Type is crucial for JSON)
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		//handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the real handler.
		next(w, r)
	}
}
