package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"os"
	"log"
)

var host string
var port int

var copydata string
var todolist []string

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

func aliases(w http.ResponseWriter, r *http.Request) {
	aliases := make([]string, 0)
	aliases = append(aliases, fmt.Sprintf("alias ccp='curl %s:%d/copy -d '", host, port))
	aliases = append(aliases, fmt.Sprintf("alias cpaste='curl %s:%d/paste'", host, port))
	aliases = append(aliases, fmt.Sprintf("alias todoadd='curl %s:%d/todo/add -d '", host, port))
	aliases = append(aliases, fmt.Sprintf("alias todorm='curl %s:%d/todo/remove -d'", host, port))
	aliases = append(aliases, fmt.Sprintf("alias todols='curl %s:%d/todo/get'", host, port))

	io.WriteString(w, strings.Join(aliases, "\n"))
}

func todoAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for key, value := range r.Form {
		if key == "item" {
			for _, item := range value {
				todolist = append(todolist, item)
			}
		} else {
			todolist = append(todolist, key)
		}
	}
}

func todoRemove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for key, value := range r.Form {

		var items []string

		if key == "item" {
			items = value
		} else {
			items = append(items, key)
		}

		if len(items) > 1 {
			io.WriteString(w, "I can only remove one item at a time\n")
		}

		if len(items) >= 1 {
			inum, err := strconv.Atoi(items[0])
			if err != nil {
				io.WriteString(w, fmt.Sprintf("Unable to parse item number: %s\n", items[0]))
			} else if inum < 0 || inum >= len(todolist) {
				io.WriteString(w, fmt.Sprintf("item number out of range: %d\n", inum))
			} else {
				removedItem := todolist[inum]
				todolist = append(todolist[:inum], todolist[inum+1:]...)
				io.WriteString(w, fmt.Sprintf("removed item: %s\n", removedItem))
			}
		}

	}

}
func todoGet(w http.ResponseWriter, r *http.Request) {
	if todolist == nil {
		return
	}

	for i, item := range todolist {
		io.WriteString(w, fmt.Sprintf("%d: %s\n", i, item))
	}
}

func main() {
	host = os.Getenv("HOST")
	if  host == "" {
		log.Fatal("could not read environment variable HOST")
	}
	var err error
	port, err = strconv.Atoi(os.Getenv("PORT"))
	if err !=nil {
		log.Fatal(fmt.Sprintf("could not read environment variable PORT: %v",err))
	}

	http.HandleFunc("/copy", copy)
	http.HandleFunc("/pastejson", pastejson)
	http.HandleFunc("/paste", paste)
	http.HandleFunc("/aliases", aliases)
	http.HandleFunc("/todo/add", todoAdd)
	http.HandleFunc("/todo/get", todoGet)
	http.HandleFunc("/todo/remove", todoRemove)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
