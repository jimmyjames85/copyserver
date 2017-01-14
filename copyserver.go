package copyserver

import (
	"fmt"
	"net/http"
)

type copyserver struct {
	host      string
	port      int
	apikey    string
	savefile  string
	todolists Todos
}

func NewCopyserver(port int, host, apikey, savefile string) *copyserver {
	c := &copyserver{
		host:   host,
		port:   port,
		apikey: apikey,
		savefile: savefile,
		todolists: NewTodos(),
	}
	return c
}

// this function blocks
func (c *copyserver) Serve() error {

	http.HandleFunc("/todo/web/add", c.handleWebAdd)
	//http.HandleFunc("/todo/web/get", todoWebGet)
	http.HandleFunc("/todo/web/getall", c.handleWebGetAll)

	http.HandleFunc("/todo/add", c.handleListAdd)
	http.HandleFunc("/todo/get", c.handleListGet)
	http.HandleFunc("/todo/getall", c.handleListGetAll)
	http.HandleFunc("/todo/getindexed", c.handleListGetIndexed)
	http.HandleFunc("/todo/remove", c.handleListRemove)
	http.HandleFunc("/todo/setindex", c.handleListSetIndex)

	http.HandleFunc("/todo/save", c.handleSaveListsToDisk) //todo put this on a cron
	http.HandleFunc("/todo/load", c.handleLoadListsFromDisk)

	return http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil)
}
