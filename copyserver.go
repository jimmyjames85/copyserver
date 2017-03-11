package copyserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type copyserver struct {
	host             string
	port             int
	apikey           string
	savefile         string
	todolists        Todos
}

func NewCopyserver(port int, host, apikey, savefile string) *copyserver {
	c := &copyserver{
		host:      host,
		port:      port,
		apikey:    apikey,
		savefile:  savefile,
		todolists: NewTodos(),
	}
	return c
}

// this function blocks
func (c *copyserver) Serve() error {

	http.HandleFunc("/todo/web/add", c.handleWebAdd)
	http.HandleFunc("/todo/web/add_redirect", c.handleWebAddWithRedirect)
	http.HandleFunc("/todo/web/remove_redirect", c.handleWebRemoveWithRedirect)
	http.HandleFunc("/todo/web/getall", c.handleWebGetAll)
	http.HandleFunc("/test", c.handleTest)
	http.HandleFunc("/todo/add", c.handleListAdd)
	http.HandleFunc("/todo/get", c.handleListGet)
	http.HandleFunc("/todo/getall", c.handleListGetAll)
	http.HandleFunc("/todo/getindexed", c.handleListGetIndexed)
	http.HandleFunc("/todo/remove", c.handleListRemove)
	http.HandleFunc("/todo/setindex", c.handleListSetIndex)
	http.HandleFunc("/todo/save", c.handleSaveListsToDisk) //todo save on every modification
	http.HandleFunc("/todo/load", c.handleLoadListsFromDisk)
	http.HandleFunc("/healthcheck", c.handleHealthcheck)

	if _, err := os.Stat(c.savefile); err == nil {
		err := c.todolists.LoadFromDisk(c.savefile)
		if err != nil {
			log.Fatalf("unable to load from previous file: %s\n", c.savefile)
		}
	}

	//todo save on every modification
	// save on a cron
	go func() {
		saveTimer := time.Tick(60 * time.Second) //todo set an env var
		for _ = range saveTimer {
			err := c.todolists.SavetoDisk(c.savefile)
			if err != nil {
				fmt.Printf(outcomeMessage(false, fmt.Sprintf("%s", err))) //todo notify
				return
			}
		}
	}()
	return http.ListenAndServe(fmt.Sprintf(":%d", c.port), nil)
}
