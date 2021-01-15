package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/status/{node}", statusHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Node: %v\n", vars["node"])
}
