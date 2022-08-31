package bot

import (
	"encoding/json"
	"log"
)

func PrettyPrint(msg string, obj interface{}) {
	log.Println(msg)
	s, _ := json.MarshalIndent(obj, "", "\t")
	log.Println(string(s))
}
