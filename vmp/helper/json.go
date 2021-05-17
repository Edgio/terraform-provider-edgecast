package helper

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
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
func Log(msg string) {
	f, err := os.OpenFile("zone.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf(">>>>>> %s", msg)
}
func LogComarison(a string, b string) {
	f, err := os.OpenFile("zone.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("%s vs %s", a, b)
}
func LogIntComarison(a int, b int) {
	f, err := os.OpenFile("zone.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("%d vs %d", a, b)
}
func LogInstanceToPrettyJson(title string, instance interface{}) {
	f, err := os.OpenFile("zone.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	e, _ := json.MarshalIndent(instance, "", "    ")

	logger.Print("=====================================================================")
	logger.Printf("[[[%s]]]:", title)
	logger.Printf("[[[Parsed Json]]]:%s", e)
	logger.Print("=====================================================================")
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
