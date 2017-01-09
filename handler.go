package copyserver

import (
	"net/http"
	"io"
	"fmt"
	"strconv"
)

const defaultlist = ""

//func _boilerPlateStuffs(w http.ResponseWriter, r *http.Request) (map[string][]string, error){
//
//	err := r.ParseForm()
//	if err != nil {
//		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
//		w.WriteHeader(http.StatusInternalServerError)
//		return nil, err
//	}
//
//	listnames := r.Form["list"]
//
//	if listnames == nil {
//		listnames = append(listnames, defaultlist)
//	}
//
//}

func HandleListAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := r.Form["item"]
	if len(items) == 0 {
		io.WriteString(w, "no items to add\n")
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	for _, listname := range listnames {
		AddItems(listname, items...)
	}

	io.WriteString(w, "success\n")
}

func HandleListRemove(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	items := r.Form["index"]
	if len(items) == 0 {
		io.WriteString(w, "no items to remove\n")
		return
	}

	var indicies []int
	for _, indexStr := range items {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			io.WriteString(w, fmt.Sprintf("Unable to parse index string: %s\n", index))
			continue
		}
		indicies = append(indicies, index)
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	for _, listname := range listnames {
		RemoveItems(listname, indicies...)
	}

	io.WriteString(w, "success\n")
}

func HandleListGetJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	io.WriteString(w, toJSON(GetLists(listnames...)))
}

func HandleListGet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	lists := GetLists(listnames...)
	for listname, list := range lists {
		io.WriteString(w, fmt.Sprintf("%s\n", listname))
		for i, item := range list {
			io.WriteString(w, fmt.Sprintf("  %d: %s\n", i, item))
		}
	}
}

func HandleListSetIndex(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, fmt.Sprintf("failed to parse form data: %s\n", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	} else if len(listnames) > 1 {
		io.WriteString(w, fmt.Sprintf("shamefully refusing request for multiple lists\n"))
		return
	}
	listname := listnames[0]
	indexstr := r.Form["index"]
	newIndexstr := r.Form["newindex"]
	if len(indexstr) != 1 || len(newIndexstr) != 1 {
		io.WriteString(w, fmt.Sprintf("please specify exactly one index and exactly one newindex\n"))
		return
	}

	index, err := strconv.Atoi(indexstr[0])
	if err != nil {
		io.WriteString(w, fmt.Sprintf("unable to parse index: %s\n", err))
		return
	}
	newIndex, err := strconv.Atoi(newIndexstr[0])
	if err != nil {
		io.WriteString(w, fmt.Sprintf("unable to parse newindex: %s\n", err))
		return
	}
	SetPriority(listname, index, newIndex)
	io.WriteString(w, "success\n")
}