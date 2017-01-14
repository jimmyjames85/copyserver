package copyserver

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const defaultlist = "main"
//var saveFileloc = "/tmp/todolists.json"
//var theapikey = "something"
//var host="localhost"
//var port=12345

//func SetHost(hostname string){
//	host=hostname
//}
//
//func SetPort(p int){
//	port=p
//}
//
//func SetSaveFileloc(fileloc string){
//	saveFileloc = fileloc
//}
//
//func SetApikey(apikey string){
//	theapikey = apikey
//}

// e.g.
//
// curl localhost:1234/todo/add -d list=grocery -d item=milk -d item=bread
//
func (c *copyserver)handleListAdd(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}

	items := r.Form["item"]
	if len(items) == 0 {
		io.WriteString(w, outcomeMessage(false, "no items to add")) //todo display available endpoints
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	for _, listname := range listnames {
		c.todolists.AddItems(listname, items...)
	}

	io.WriteString(w, outcomeMessage(true, ""))
}

// e.g.
//
// curl localhost:1234/todo/remove -d list=grocery -d index=0 -d index=3
//
func (c *copyserver) handleListRemove(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}

	items := r.Form["index"]
	if len(items) == 0 {
		io.WriteString(w, outcomeMessage(false, "no items to remove")) //todo display available endpoints
		return
	}

	var indicies []int
	for _, indexStr := range items {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			io.WriteString(w, outcomeMessage(false, fmt.Sprintf("Unable to parse index string: %s\n", index)))
			continue
		}
		indicies = append(indicies, index)
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	for _, listname := range listnames {
		c.todolists.RemoveItems(listname, indicies...)
	}

	io.WriteString(w, outcomeMessage(true, ""))
}

// e.g.
//
// curl localhost:1234/todo/get -d list=grocery
//
func (c *copyserver) handleListGet(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	io.WriteString(w, ToJSON(c.todolists.GetLists(listnames...)))
}

// e.g.
//
// curl localhost:1234/todo/getall
//
func (c *copyserver) handleListGetAll(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}
	io.WriteString(w, ToJSON(c.todolists.GetAllLists()))
}

// e.g.
//
// curl localhost:1234/todo/save
//
func (c *copyserver) handleSaveListsToDisk(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}
	password := r.Form["password"]
	if password == nil || len(password) != 1 || password[0] != c.apikey {
		io.WriteString(w, outcomeMessage(false, "incorrect credentials"))
		return
	}
	err := c.todolists.SavetoDisk(c.savefile)
	if err != nil {
		io.WriteString(w, outcomeMessage(false, fmt.Sprintf("%s", err)))
		return
	}

	io.WriteString(w, outcomeMessage(true, ""))
}

// e.g.
//
// curl localhost:1234/todo/load
//
func (c *copyserver) handleLoadListsFromDisk(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}

	err := c.todolists.LoadFromDisk(c.savefile)
	if err != nil {
		io.WriteString(w, outcomeMessage(false, fmt.Sprintf("%s", err)))
		return
	}

	io.WriteString(w, outcomeMessage(true, ""))
}

// e.g.
//
// curl localhost:1234/todo/getindexed -d list=grocery
//
func (c *copyserver) handleListGetIndexed(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	}

	lists := c.todolists.GetLists(listnames...)
	for listname, list := range lists {
		io.WriteString(w, fmt.Sprintf("%s\n", listname))
		for i, item := range list {
			io.WriteString(w, fmt.Sprintf("  %d: %s\n", i, item))
		}
	}
}

// e.g.
//
// curl localhost:1234/todo/setindex -d list=grocery -d index=2 -d newindex=4
//
func (c *copyserver) handleListSetIndex(w http.ResponseWriter, r *http.Request) {
	if !handleParseFormData(w, r) {
		return
	}

	listnames := r.Form["list"]
	if listnames == nil {
		listnames = append(listnames, defaultlist)
	} else if len(listnames) > 1 {
		io.WriteString(w, outcomeMessage(false, "shamefully refusing request for multiple lists"))
		return
	}
	listname := listnames[0]
	indexstr := r.Form["index"]
	newIndexstr := r.Form["newindex"]
	if len(indexstr) != 1 || len(newIndexstr) != 1 {
		io.WriteString(w, outcomeMessage(false, "please specify exactly one index and exactly one newindex"))
		return
	}

	index, err := strconv.Atoi(indexstr[0])
	if err != nil {
		io.WriteString(w, outcomeMessage(false, fmt.Sprintf("unable to parse index: %s", err)))
		return
	}
	newIndex, err := strconv.Atoi(newIndexstr[0])
	if err != nil {
		io.WriteString(w, outcomeMessage(false, fmt.Sprintf("unable to parse newindex: %s\n", err)))
		return
	}
	c.todolists.SetPriority(listname, index, newIndex)
	io.WriteString(w, outcomeMessage(true, ""))
}

func (c *copyserver) handleWebAdd(w http.ResponseWriter, r *http.Request) {
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
`, c.host, c.port))
}

func (c *copyserver) handleWebGetAll(w http.ResponseWriter, r *http.Request) {
	html := "<html>"
	for listname, list := range c.todolists {
		html += fmt.Sprintf("%s<hr><table>",listname )
		for i, item := range list {
			rmBtn := fmt.Sprintf(`<form action="http://%s:%d/todo/remove">
			<input type="hidden" name="list" value="%s">
			<input type="hidden" name="index" value="%d">
			<input type="submit" value="rm"></form>`, c.host, c.port, listname, i)
			html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", i, item, rmBtn)
		}
		html += "</table><br>"
	}
	html += "</html>"
	io.WriteString(w, html)
}
func handleParseFormData(w http.ResponseWriter, r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		io.WriteString(w, outcomeMessage(false, fmt.Sprintf("failed to parse form data: %s", err)))
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}
	return true
}

//todo what is a better way to do this????
func outcomeMessage(outcome bool, msg string) string {
	m := make(map[string]string)
	if len(msg) != 0 {
		m["message"] = msg
	}
	m["result"] = "failed"
	if outcome {
		m["result"] = "success"
	}
	return ToJSON(m)
}
