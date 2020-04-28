package main

// Generates sudoku puzzles.

import (
	"fmt"
	"github.com/typingduck/sodacoup/sodacouplib"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

func main() {
	verbose := len(os.Args) > 1 && os.Args[1] == "-v"
	if !verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	sodacouplib.RandomSeed()
	var i int64
	m := 82
	for {
		i++
		rand.Seed(i)
		s, e := sodacouplib.GenerateProblem()
		if e != nil {
			printFatal("error:%s", e)
		}

		c := s.SetCount()
		if c < m {
			fmt.Printf("%d seed puzzle (%d filled):\n", i, c)
			fmt.Println(s)
			m = c
		}
	}
}

func printFatal(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
