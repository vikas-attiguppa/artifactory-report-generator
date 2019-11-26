package main

import (
	"artifactory-report-generator/artifactory"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	response, err := artifactory.DefaultClient().GetArchivesForRepo(r.URL.Path[1:])
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
