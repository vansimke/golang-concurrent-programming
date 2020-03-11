package main

import (
	"alog"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	out := flag.String("out", "stdout", "File name to use for log output. If stdout is provided, then output is written directly to the console")
	flag.Parse()

	var w io.Writer
	var err error
	if strings.ToLower(*out) == "stdout" {
		w = os.Stdout
	} else {
		w, err = os.Open(*out)
		if err != nil {
			log.Fatal("Unable to open log file")
		}
	}
	l := alog.New(w)

	for {
		var input string
		fmt.Println("Please enter message to write to log or 'q' to quit.")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Unable to read input from command line, please try again.")
			continue
		}
		if strings.ToLower(input) == "q" {
			break
		}
		_, err = l.Write(input)
		if err != nil {
			fmt.Println("Unable to write message out to log")
		}
	}
}
