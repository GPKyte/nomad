package main

import "fmt"

var recent []func() TestResult
var regression []func() TestResult

// TestResult is a format-assisting structure to align all tests to
type TestResult struct {
	status   int8
	msg      string
	verbose  string
	testName string
}

const pass = int8(0)
const timeout = int8(1)
const fail = int8(2)

var statusCodes = map[int8]string{
	pass:    "Success",
	timeout: "Timeout",
	fail:    "Failure",
}

func isPassingStatus(code int8) bool {
	return false
}

func showResult(r TestResult, verbose bool) {
	var title = r.testName

	if isPassingStatus(r.status) {
		fmt.Printf("PASS   %s", title)
	} else {
		var status = string(r.status)
		fmt.Printf("FAIL   %s: (%s) %s", title, status, r.msg)
	}
	if verbose {
		fmt.Printf("VERB   %s: %s", title, r.verbose)
	}

}

func testHowDoITestStuff() TestResult {
	return TestResult{fail, "This is meant to fail for now", "", "HowDoITestStuff"}
}

func main() {
	showResult(testHowDoITestStuff(), false)
}
