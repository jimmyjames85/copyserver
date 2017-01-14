package copyserver

import (
	"encoding/json"
	"fmt"
	"time"
)

func logf(format string, a ... interface{}) {
	fmt.Printf("%s: %s\n", time.Now().String(), fmt.Sprintf(format, a...))
}

func ToJSON(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("failed to marshal data: %s", err)
		return fmt.Sprintf("%v", obj)
	}
	return string(b)
}
