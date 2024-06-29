package main

import "net/http"

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Any code here will execute on the way down the chain.
		w.Header().Set("Content-Security-Policy","default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options","nosniff")
		w.Header().Set("X-Frame-Options","deny")
		w.Header().Set("X-XSS-Protection","0")

		next.ServeHTTP(w, r)
		// Any code here will execute on the way back up the chain.
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}