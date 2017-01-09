package copyserver

import (
	"fmt"
	"sort"
)

var copydata string
var lists map[string][]string

func populateListsVariable() {
	// TODO find a better way todo this
	// probably make a func NewCopyServer
	if lists == nil {
		lists = make(map[string][]string)
	}
}

func QuickCopy(s string) string {
	ret := copydata
	copydata = s
	return ret
}

func QuickPaste() string {
	return copydata
}
func GetLists(listnames...string) map[string][]string {
	ret := make(map[string][]string,0)

	for _, listname := range listnames {
		if list, ok := lists[listname]; ok {

			ret[listname] = append(ret[listname], list...)
		}
	}
	return ret
}

func AddItems(listname string, items ...string) {
	populateListsVariable()
	lists[listname] = append(lists[listname], items...)
}

func RemoveItems(listname string, indicies ...int) ([]string, error) {
	populateListsVariable()
	if _, ok := lists[listname]; !ok {
		return nil, fmt.Errorf("no such list: %s", listname)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(indicies))) // TODO why does this work; i think sort.Reverse defines diff sort methods maybe

	removed := make([]string, 0)

	for _, index := range indicies {
		if index < 0 || index >= len(lists[listname]) {
			continue
		}
		removed = append(removed, lists[listname][index])
		lists[listname] = append(lists[listname][:index], lists[listname][index + 1:]...)
	}
	return removed, nil
}

func SetPriority(listname string, index int, newIndex int) {
	populateListsVariable()

	list, ok := lists[listname]
	if !ok {
		return
	} else if index < 0 || index >= len(list) || newIndex < 0 || newIndex >= len(list) {
		return
	}

	// by default we move [index] to the right
	direction := 1
	if newIndex < index {
		direction = -1
	}

	for index != newIndex{
		lists[listname][index], lists[listname][index+direction] = lists[listname][index+direction], lists[listname][index]
		index+= direction

	}
}
