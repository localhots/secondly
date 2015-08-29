package secondly

import (
	"encoding/json"
	"log"
	"net/http"
)

func startServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/fields.json", fieldsHandler)

	log.Println("Starting configuration server on", addr)
	go http.ListenAndServe(addr, mux)
}

func fieldsHandler(rw http.ResponseWriter, req *http.Request) {
	fields := extractFields(config, "")
	body, _ := json.Marshal(fields)
	rw.Write(body)
}
