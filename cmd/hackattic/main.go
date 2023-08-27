package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pdmatrix/hackattic/pkg/challenge"
)

var challengeName string
var runUntillPassed bool
var playground bool

func init() {
	flag.StringVar(&challengeName, "challenge", "", "Name of the challenge")
	flag.BoolVar(&runUntillPassed, "run-untill-passed", false, "Run untill passed")
	flag.BoolVar(&playground, "playground", false, "Use playground mode")
	flag.Parse()
}

func main() {
	if challengeName == "" {
		log.Fatal("You need to provide challenge name")
	}

	for {
		res, err := challenge.GetSolution(challengeName, playground)
		if err != nil {
			fmt.Printf("Error while solving %s challenge: %v", challengeName, err)
			os.Exit(1)
		}

		if !runUntillPassed {
			break
		}

		if strings.Contains(res, "passed") {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
