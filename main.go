package main

import (
	fmt "github.com/starkandwayne/goutils/ansi"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(200)
		fmt.Fprintf(w, "@B{cloudvet} is up and @G{ready for testing!}\n")
	})
	http.ListenAndServe(":"+port, nil)
}
