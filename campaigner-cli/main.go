package main

import (
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"log"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/henrocdotnet/active-campaigner/campaigner"
	"github.com/kelseyhightower/envconfig"
)

var (
	config envConfig
)

type envConfig struct {
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
// TODO(improvements): Errors should probably go to STDERR.
func main() {
	args := os.Args

	c := campaigner.Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}

	if len(args) < 3 {
		printUsage()
		os.Exit(-1)
	}

	switch args[1] {
	case "contact":
		switch args[2] {
		case "delete":
			if len(args) != 4 {
				printUsage()
				os.Exit(-1)
			}

			id, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			err = c.ContactDelete(id)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

		case "list":
			var (
				limit = campaigner.DEFAULT_LIST_LIMIT
				offset = campaigner.DEFAULT_LIST_OFFSET
				err error
			)

			if len(args) == 5 {
				limit, err = strconv.Atoi(args[3])
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}

				offset, err = strconv.Atoi(args[4])
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
			}

			if len(args) != 3 && len(args) != 5 {
				printUsage()
				os.Exit(-1)
			}

			r, err := c.ContactList(limit, offset)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			fmt.Printf("Listing Contacts:\n\n")
			for _, y := range r.Contacts {
				fmt.Printf("\t%d: %s\n", y.ID, y.EmailAddress)
			}
			fmt.Println("")
			fmt.Printf("Limit: %d, Offset: %d, Total: %d\n", limit, offset, r.Meta.Total)
			fmt.Println("")

		case "read":
			id, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			r, err := c.ContactRead(id)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Printf("% #v\n", pretty.Formatter(r))

			fields, err := c.FieldList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			for _, fieldValue := range r.FieldValues {
				var n string
				for _, b := range fields.Fields {
					//log.Printf("%d %d %s\n", b.ID, fieldValue.FieldID.Int64(), b.Title)
					if b.ID == fieldValue.FieldID.Int64() {
						n = b.Title
					}

				}
				fmt.Printf("Field: %d %s %s\n", fieldValue.FieldID, n, fieldValue.Value)

			}

		case "tags":
			id, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			r, err := c.ContactTagReadByContactID(id)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Printf("% #v\n", pretty.Formatter(r))
		}
	case "field":
		switch args[2] {
		case "list":
			r, err := c.FieldList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			for _, y := range r.Fields {
				fmt.Printf("\tID: %d, Title: %s\n", y.ID, y.Title)
			}
		}
	case "list":
		switch args[2] {
		case "list":
			r, err := c.ListList()
			if err != nil {
				handleError(err)

			}

			for _, y := range r.Lists {
				fmt.Printf("\tID: %d: Name: %s\n", y.ID, y.Name)
			}
		case "read":
			id, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				handleError(err)
			}

			r, err := c.ListRead(id)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("% #v\n", pretty.Formatter(r))
		}
	case "org":
		switch args[2] {
		case "read":
			id, err := strconv.ParseInt(args[3], 10, 64)
			handleError(err)

			r, err := c.OrganizationRead(id)
			handleError(err)

			fmt.Printf("% #v\n", pretty.Formatter(r))
		case "delete":
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
		case "list":
			var (
				limit = campaigner.DEFAULT_LIST_LIMIT
				offset = campaigner.DEFAULT_LIST_OFFSET
				err error
			)

			if len(args) == 5 {
				limit, err = strconv.Atoi(args[3])
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}

				offset, err = strconv.Atoi(args[4])
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
			}

			if len(args) != 3 && len(args) != 5 {
				printUsage()
				os.Exit(-1)
			}

			r, err := c.OrganizationList(limit, offset)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			fmt.Printf("Listing Organizations:\n\n")
			for _, y := range r.Organizations {
				fmt.Printf("\t%d: %s (%s)\n", y.ID, y.Name, y.ContactCount)
			}
			fmt.Println("")
			fmt.Printf("Limit: %d, Offset: %d, Total: %d\n\n", limit, offset, r.Meta.Total)
		default:
			printUsage()
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
			fMap := template.FuncMap{
				"cleanTagName": func(s string) string {
					re, err := regexp.Compile("[^a-zA-Z0-9]+")
					if err != nil {
						fmt.Println(err)
						os.Exit(-1)
					}
					s = re.ReplaceAllString(s, "")

					return s
				},
			}

			r, err := c.TagList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			t1 := template.Must(template.New("template").
				Funcs(fMap).
				Parse(classTemplate))
			err = t1.Execute(os.Stdout, r)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			fields, err := c.FieldList()
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			t2 := template.Must(template.New("template").
				Funcs(fMap).
				Parse(fieldsTemplate))
			err = t2.Execute(os.Stdout, fields)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

			lists, err := c.ListList()
			if err != nil {
				handleError(err)
			}
			t3 := template.Must(template.New("template").Funcs(fMap).Parse(listsTemplate))
			err = t3.Execute(os.Stdout, lists)
			if err != nil {
				handleError(err)
			}

		case "read":
			id, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				handleError(err)
			}
			r, err := c.TagRead(id)
			if err != nil {
				handleError(err)
			}

			fmt.Printf("Tag: %# v\n", pretty.Formatter(r))
		}

	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Printf("%#v\n", os.Args)

	tmpl := `
Usage:
	cli <contact|tag|org> 

	contact <list|read>
		     list: List contacts.
		read <id>: Read contact.
	
	org <delete|list>
		delete <id>: Delete organization.
		       list: List organizations.
	tag <delete|generate|list>
		list: List tags.

`

	fmt.Printf(tmpl)
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}
}
