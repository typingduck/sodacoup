package main

// Read a sudoku problem from stdin and solve it.
// If a problem is not given on stdin then just do a sample problem found
// in 'sample_problem'.
//
// Running:
//    Run just the sample file:
//        ./sodacoup -s
//    Give a problem using stdin:
//         cat problem | ./sodacoup
//    Print steps done:
//         cat problem | ./sodacoup -v
//

import (
	"fmt"
	"github.com/typingduck/sodacoup/sodacouplib"
	"io/ioutil"
	"log"
	"os"
)

const sampleFile = "sample_problem"

func main() {
	verbose, getProblemFromStdin := parseCommandline()

	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	problem := loadProblem(getProblemFromStdin)

	sud, err := sodacouplib.NewSudokuSquare(problem)
	if err != nil {
		printFatal("Problem text doesn't look like a valid sudoku problem")
	}
	fmt.Println("PROBLEM:")
	fmt.Println(sud)

	err = sud.Solve()
	if err != nil {
		printFatal("error: %s", err)
	}
	fmt.Println("SOLUTION:")
	fmt.Println(sud)
}

func parseCommandline() (bool, bool) {
	verbose := len(os.Args) > 1 && os.Args[1] == "-v"
	doSample := (!verbose && len(os.Args) > 1) || (verbose && len(os.Args) > 2)
	return verbose, !doSample
}

func loadProblem(getProblemFromStdin bool) string {
	if getProblemFromStdin {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			printFatal("Failed reading from stdin. %s", err)
		}
		return string(bytes)
	}
	bytes, err := ioutil.ReadFile(sampleFile)
	if err != nil {
		printFatal("Failed reading from sample file. %s", err)
	}
	return string(bytes)
}

func printFatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
