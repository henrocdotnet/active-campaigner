package campaigner

import (
	"flag"
	"log"
	"os"
	"reflect"
	"runtime"
	"testing"
)


type unitTest func(*testing.T)


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


func runTestWithName(t *testing.T, u unitTest) {
	rv := reflect.ValueOf(u)
	name := runtime.FuncForPC(rv.Pointer()).Name()
	t.Run(name, u)
}
