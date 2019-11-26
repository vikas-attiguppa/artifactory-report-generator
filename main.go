package main

import (
	"fmt"
	"github.com/vikas-attiguppa/artifactory-report-generator/artifactory"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	response, err := artifactory.DefaultClient().GetTopArchivesForRepo(r.URL.Path[1:])
	if err != nil {
		fmt.Print("Something went wrong" + err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
