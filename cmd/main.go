package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"strings"
	"os"
	"log"
	"strconv"
	"github.com/jimmyjames85/copyserver"
)

var host string
var port int

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

func aliases(w http.ResponseWriter, r *http.Request) {
	aliases := make([]string, 0)
	aliases = append(aliases, fmt.Sprintf("alias ccp='curl %s:%d/copy -d '", host, port))
	aliases = append(aliases, fmt.Sprintf("alias cpaste='curl %s:%d/paste'", host, port))
	aliases = append(aliases, fmt.Sprintf("alias todoadd='curl %s:%d/todo/add -d '", host, port))
	aliases = append(aliases, fmt.Sprintf("alias todorm='curl %s:%d/todo/remove -d'", host, port))
	aliases = append(aliases, fmt.Sprintf("alias todols='curl %s:%d/todo/get'", host, port))

	io.WriteString(w, strings.Join(aliases, "\n"))
}

func todoWebGetAll(w http.ResponseWriter, r *http.Request) {

}

func todoWebGet(w http.ResponseWriter, r *http.Request) {

}


func main() {
	host = os.Getenv("HOST")
	if len(host) == 0 {
		log.Fatalf("could not read environment variable HOST\n")
	}

	var err error
	port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("could not read environment variable PORT: %v\n", err)
	}

	apikey := os.Getenv("APIKEY")
	if len(apikey) == 0 {
		log.Fatalf("could not read environment variable APIKEY\n")
	}

	savefile := os.Getenv("SAVEFILE")
	if len(savefile) == 0 {
		log.Fatalf("could not read environment variable SAVEFILE\n")
	}

	copyserver := copyserver.NewCopyserver(port,host,apikey,savefile)
	copyserver.Serve()

	//http.HandleFunc("/copy", copy)
	//http.HandleFunc("/pastejson", pastejson)
	//http.HandleFunc("/paste", paste)
	//http.HandleFunc("/aliases", aliases)

}
