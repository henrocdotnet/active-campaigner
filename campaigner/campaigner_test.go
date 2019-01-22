package campaigner

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type configSetup struct {
	APIToken      string `envconfig:"api_token"       required:"true"`
	BaseURL       string `envconfig:"base_url"        required:"true"`
	UnitTestEmail string `envconfig:"unit_test_email" required:"true"`
	UnitTestPhone string `envconfig:"unit_test_phone" default:"2125551212"`
}

type unitTest func(*testing.T)

var (
	config configSetup
	C      Campaigner
	NOW    = time.Now().Format("20060102_150405")
)

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
	l := flag.Lookup("test.run")
	if len(l.Value.String()) > 0 {
		// log.Printf("test run? %s\n", l)
		os.Exit(m.Run())

	}

	// Set filter to group runners.  This allows tests to be run both piecemeal or ordered / "suite".
	err := flag.Set("test.run", "TestContactSuite|TestTagSuite|TestContactTaggingSuite|TestOrganizationSuite")
	if err != nil {
		log.Fatal(err)
	}
	l2 := flag.Lookup("test.run")
	log.Printf("test run? %s\n", l2)

	os.Exit(m.Run())
}

func TestPrint(t *testing.T) {
	log.Println("Testing print.")
}

func TestRunTestWithName(t *testing.T) {
	runTestWithPackagePath(t, TestPrint)
}

// Runs a unit test and prints it's path.
func runTestWithPackagePath(t *testing.T, u unitTest) {
	rv := reflect.ValueOf(u)
	name := runtime.FuncForPC(rv.Pointer()).Name()
	t.Run(name, u)
}
