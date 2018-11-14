package campaigner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

func logFormattedJson(message string, data interface{}) {
	tmp, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if len(message) > 0 {
		log.Printf("\n\n" + message + "\n")
	}
	log.Printf("\n%s", string(tmp))
}

func writeIndentedJson(path string, data []byte) {
	var o bytes.Buffer

	json.Indent(&o, data, "", "\t")

	err := ioutil.WriteFile(path, o.Bytes(), 0644)
	if err != nil {
		log.Printf("Could not write indented json file %s: %s", path, err)
	}
}
