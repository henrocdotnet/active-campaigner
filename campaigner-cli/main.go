package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/henrocdotnet/active-campaigner/campaigner"
	"github.com/kelseyhightower/envconfig"
)

var (
	config configSetup
)

type configSetup struct {
	ApiToken string `envconfig:"api_token"`
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

	c := campaigner.Campaigner{ApiToken: config.ApiToken, BaseURL: config.BaseURL}

	if len(args) < 2 {
		printUsage()
		os.Exit(-1)
	}

	switch args[1] {
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
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Printf("%#v\n", os.Args)
}
