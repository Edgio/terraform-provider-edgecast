package helper

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

// Log API request body, method, url in json pretty format.
func LogRequestBody(method string, url string, body interface{}) {
	fb, _ := json.MarshalIndent(body, "", "    ")
	// Read the Body content
	log.Print("=====================================================================")
	log.Printf("[[[REQUEST-URI]]]:[%s] %s", method, url)
	log.Printf("[[[REQUEST-BODY]]]:%s", fb)
	log.Print("=====================================================================")
}

// Log json string with pretty format with a message
func LogPrettyJson(message string, jsonString string) {

	log.Print("=====================================================================")
	log.Printf("[[[%s]]]:", message)
	log.Printf("[[[Json]]]:%s", jsonPrettyPrint(jsonString))
	log.Print("=====================================================================")
}

// Log message in the file
// message: message
// instance: any data structure, like map, slice, instance of struct
// file: file name. file is created in the folder that tf.exe exeduted
func Log(msg string, file string) {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf(">>>>>> %s", msg)
}

// Log key value pair or strings to compare in the log
// message: message
// instance: any data structure, like map, slice, instance of struct
// file: file name. file is created in the folder that tf.exe exeduted
func LogComarison(a string, b string, file string) {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("%s vs %s", a, b)
}

// Log key value pair or two int instances in the log
// message: message
// instance: any data structure, like map, slice, instance of struct
// file: file name. file is created in the folder that tf.exe exeduted
func LogIntComarison(a int, b int, file string) {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Printf("%d vs %d", a, b)
}

// Log jsonfied instance with message
// message: message
// instance: any data structure, like map, slice, instance of struct
// file: file name. file is created in the folder that tf.exe exeduted
func LogInstanceToPrettyJson(message string, instance interface{}, file string) {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	e, _ := json.MarshalIndent(instance, "", "    ")

	logger.Print("=====================================================================")
	logger.Printf("[[[%s]]]:", message)
	logger.Printf("[[[Parsed Json]]]:%s", e)
	logger.Print("=====================================================================")
}

// Make json string formatted in terraform.log
func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
