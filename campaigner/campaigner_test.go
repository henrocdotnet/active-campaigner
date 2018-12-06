package campaigner

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"
)

type configSetup struct {
	APIToken string `envconfig:"api_token"`
	BaseURL  string `envconfig:"base_url"`
}

type unitTest func(*testing.T)

var (
	config configSetup
	C Campaigner
)

// TODO: Move flag parsing into TestMain at some point.
func init() {
	err := envconfig.Process("ac", &config)
	if err != nil {
		log.Fatal(err)
	}

	C = Campaigner{APIToken: config.APIToken, BaseURL: config.BaseURL}
}


func TestMain(m *testing.M) {
	if flag.Parsed() == false {
		flag.Parse()
	}

	// Allow tests to be run based on test.run filter.
	log.Println("HERE?")
	l := flag.Lookup("test.run")
	if len(l.Value.String()) > 0 {
		log.Printf("test run? %s\n", l)
		os.Exit(m.Run())

	}

	// Set filter to group runners.  This allows tests to be run both piecemeal or ordered / "suite".
	flag.Set("test.run", "TestOrganizationList_Success")
	l2 := flag.Lookup("test.run")
	log.Printf("test run? %s\n", l2)

	os.Exit(m.Run())
}

func TestPrint(t *testing.T) {
	log.Println("Testing print.")
}

func TestRunTestWithName(t *testing.T) {
	runTestWithName(t, TestPrint)
}

// Runs a unit test and prints it's name.
func runTestWithName(t *testing.T, u unitTest) {
	rv := reflect.ValueOf(u)
	name := runtime.FuncForPC(rv.Pointer()).Name()
	t.Run(name, u)
}
