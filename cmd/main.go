package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var copydata string

func copy(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(r.Form)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to marshal form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	copydata = string(b)
}

func pastejson(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, copydata)
}

func paste(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	bytes := []byte(copydata)

	err := json.Unmarshal(bytes, &data)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to unmarshal copydata: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// only print the data
	if len(data) == 1 {
		for key, val := range data {
			if fmt.Sprintf("%v", val) == "[]" {
				io.WriteString(w, key)
			} else {
				io.WriteString(w, copydata)
			}
		}
	} else {
		io.WriteString(w, copydata)
	}
}

func main() {
	http.HandleFunc("/copy", copy)
	http.HandleFunc("/pastejson", pastejson)
	http.HandleFunc("/paste", paste)
	http.ListenAndServe(":27182", nil)
}
