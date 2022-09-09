package main

import (
	"fmt"
	"net/http"
)

func redirectHandler(data string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/twitback" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != "GET" {
			http.Error(w, "Method not supported", http.StatusNotFound)
			return
		}

		stat := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")

		if stat != data {
			http.Error(w, "Not allowed", http.StatusForbidden)
			return
		}

		// store in data base or pass through go channels
		// fmt.Printf("code: %s", code)

		getAccesToken(code)
	}

}

func runServer(data string) {
	http.HandleFunc("/twitback", redirectHandler(data))
	fmt.Println("listing on port 5000...")
	http.ListenAndServe(":5000", nil)

}
