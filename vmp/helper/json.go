package helper

import (
	"bytes"
	"encoding/json"
	"log"
)

func LogRequestBody(method string, url string, body interface{}) {
	fb, _ := json.MarshalIndent(body, "", "    ")
	// Read the Body content
	log.Print("=====================================================================")
	log.Printf("[[[REQUEST-URI]]]:[%s] %s", method, url)
	log.Printf("[[[REQUEST-BODY]]]:%s", fb)
	log.Print("=====================================================================")
}

func LogPrettyJson(title string, jsonString string) {

	log.Print("=====================================================================")
	log.Printf("[[[%s]]]:", title)
	log.Printf("[[[REQUEST-BODY]]]:%s", jsonPrettyPrint(jsonString))
	log.Print("=====================================================================")
}

func LogInstanceToPrettyJson(title string, instance interface{}) {
	e, _ := json.MarshalIndent(instance, "", "    ")

	log.Print("=====================================================================")
	log.Printf("[[[%s]]]:", title)
	log.Printf("[[[REQUEST-BODY]]]:%s", e)
	log.Print("=====================================================================")
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
