package campaigner

import (
	"bytes"
	"encoding/json"
	"github.com/kr/pretty"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

// Fixes issues caused by some ID numbers being returned as both strings and numbers in the JSON (from the same API calls).
type Int64json int64

func (i Int64json) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *Int64json) UnmarshalJSON(data []byte) error {
	re := regexp.MustCompile("[^0-9]")
	s := re.ReplaceAllString(string(data), "")

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*i = Int64json(n)
	return nil
}

func dump(i interface{}) {
	log.Printf("%# v", pretty.Formatter(i))
}

func dumpWithMessage(i interface{}, m string) {
	log.Printf("%s\n%# v", m, pretty.Formatter(i))
}


func logFormattedJSON(message string, data interface{}) {
	tmp, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if len(message) > 0 {
		log.Printf("\n\n" + message + "\n")
	}
	log.Printf("\n%s", string(tmp))
}

func writeIndentedJSON(path string, data []byte) {
	var o bytes.Buffer

	err := json.Indent(&o, data, "", "\t")
	if err != nil {
		log.Printf("Could not indent JSON: %s", err)
	}

	err = ioutil.WriteFile(path, o.Bytes(), 0644)
	if err != nil {
		log.Printf("Could not write indented json file %s: %s", path, err)
	}
}
