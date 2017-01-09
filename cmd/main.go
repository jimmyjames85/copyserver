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

func todoWebAdd(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, fmt.Sprintf(`<html>
  <form action="http://%s:%d/todo/add">
    <input type="text" name="item"><br>
    <input type="text" name="item"><br>
    <input type="text" name="item"><br>
    <input type="text" name="item"><br>
    <input type="text" name="item"><br>
    <input type="submit" value="Submit"><br>
  </form>
</html>
`, host, port))
}
//
//func todoSaveToDisk(w http.ResponseWriter, r *http.Request) {
//	s := toJSON(todolist)
//	fmt.Printf("saving: %s\n", s)
//	writeToDiskAndIgnoreErr(s, "/tmp/copyserver.data")
//}
//func todoLoadToDisk(w http.ResponseWriter, r *http.Request) {
//	s := loadFromDisk("/tmp/copyserver.data")
//	io.WriteString(w, s)
//}
//
//func loadFromDisk(fileloc string) string{
//	b, err := ioutil.ReadFile(fileloc)
//	if err!=nil {
//		return ""
//	}
//	return string(b)
//}
//
//func writeToDiskAndIgnoreErr(data string, fileloc string) {
//	d := []byte(data)
//	err := ioutil.WriteFile(fileloc, d, 0644)
//	if err != nil {
//		fmt.Printf("err writint to file %s: %s\n", fileloc, err)
//	}
//}
//func handleLoadSavedData(w http.ResponseWriter, r *http.Request) {
//	io.WriteString(w, loadFromDisk("/tmp/copyserver.data"))
//}

func main() {
	host = os.Getenv("HOST")
	if host == "" {
		log.Fatal("could not read environment variable HOST")
	}
	var err error
	port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(fmt.Sprintf("could not read environment variable PORT: %v", err))
	}

	http.HandleFunc("/copy", copy)
	http.HandleFunc("/pastejson", pastejson)
	http.HandleFunc("/paste", paste)
	http.HandleFunc("/aliases", aliases)

	http.HandleFunc("/todo/web/add", todoWebAdd)

	http.HandleFunc("/todo/add", copyserver.HandleListAdd)
	http.HandleFunc("/todo/get", copyserver.HandleListGet)
	http.HandleFunc("/todo/getjson", copyserver.HandleListGetJSON)
	http.HandleFunc("/todo/remove", copyserver.HandleListRemove)
	http.HandleFunc("/todo/setindex", copyserver.HandleListSetIndex)

	//http.HandleFunc("/todo/save", todoSaveToDisk)
	//http.HandleFunc("/todo/load", todoLoadToDisk)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}
