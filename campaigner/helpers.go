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

type int64json int64;

func (i int64json) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(i))
}

func (i *int64json) UnmarshalJSON(data []byte) error {
	re := regexp.MustCompile("[^0-9]")
	s := re.ReplaceAllString(string(data), "")

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*i = int64json(n)
	return nil
}

func dump(i interface{}) {
	log.Printf("%# v", pretty.Formatter(i))
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

	json.Indent(&o, data, "", "\t")

	err := ioutil.WriteFile(path, o.Bytes(), 0644)
	if err != nil {
		log.Printf("Could not write indented json file %s: %s", path, err)
	}
}
