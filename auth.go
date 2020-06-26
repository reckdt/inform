package main

import (
	"fmt"
	"net/http"
)

func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth")

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		username, password, authOK := r.BasicAuth()

		fmt.Println(username, password, authOK)

		h.ServeHTTP(w, r)
	}
}
