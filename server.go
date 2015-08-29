package secondly

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/GeertJohan/go.rice"
)

func startServer(addr string) {
	staticHandler := http.FileServer(rice.MustFindBox("static").HTTPBox())

	mux := http.NewServeMux()
	mux.HandleFunc("/fields.json", fieldsHandler)
	mux.HandleFunc("/save", saveHandler)

	// Static
	mux.Handle("/app.js", staticHandler)
	mux.Handle("/app.css", staticHandler)
	mux.Handle("/config.html", staticHandler)
	// Redirect from root to a static file. Ugly yet effective.
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/" {
			http.Redirect(rw, req, "/config.html", http.StatusMovedPermanently)
		}
	})

	log.Println("Starting configuration server on", addr)
	go http.ListenAndServe(addr, mux)
}

func fieldsHandler(rw http.ResponseWriter, req *http.Request) {
	fields := extractFields(config, "")
	body, err := json.Marshal(fields)
	if err != nil {
		panic(err)
	}

	rw.Write(body)
}

func saveHandler(rw http.ResponseWriter, req *http.Request) {
	cbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	updateConfig(cbody)
	writeConfig()

	resp := struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
	}{
		Success: true,
		Msg:     "Config successfully updated",
	}
	body, _ := json.Marshal(resp)
	rw.Write(body)
}
