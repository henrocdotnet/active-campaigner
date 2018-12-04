package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/henrocdotnet/active-campaigner/campaigner"
	"github.com/kelseyhightower/envconfig"
)

var (
	config configSetup
)

type configSetup struct {
	APIToken string `envconfig:"api_token"`
	BaseURL  string `envconfig:"base_url"`
}

func init() {
	err := envconfig.Process("ac", &config)
	if err != nil {
		log.Fatal(err)
	}
}

// TODO(improvements): Arguments hacked in here, spike solution, clean up.
func main() {
	args := os.Args

	c := campaigner.Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	if len(args) < 2 {
		printUsage()
		os.Exit(-1)
	}

	switch args[1] {
	case "contact":
		switch args[2] {
		case "list":
			err := c.ContactList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

		}
	case "org":
		switch args[2] {
		case "delete":
			fmt.Printf("FOUND DELETE ORG\n")

			id, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			err = c.OrganizationDelete(int64(id))
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Printf("Organization %d deleted successfully.\n", id)

			break
		case "list":
			r, err := c.OrganizationList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Printf("%#v\n", r)
		}
	case "tag":
		switch args[2] {
		case "list":
			r, err := c.TagList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			f, err := json.MarshalIndent(r, "", "\t")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", string(f))
		case "delete":
			id, err := strconv.Atoi(args[3])
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			err = c.TagDelete(int64(id))
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Printf("Tag %d deleted successfully.\n", id)
		case "generate":
			r, err := c.TagList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			t := template.Must(template.New("template").
				Funcs(template.FuncMap{
					"cleanTagName": func(s string) string {
						re, err := regexp.Compile("[^a-zA-Z0-9]+")
						if err != nil {
							fmt.Println(err)
							os.Exit(-1)
						}
						s = re.ReplaceAllString(s, "")

						// s = strings.Replace()

						return s
					},
				}).
				Parse(classTemplate))
			err = t.Execute(os.Stdout, r)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
		}

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Printf("%#v\n", os.Args)
}
