package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shreyassanthu77/cisp/interpreter"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: crap <input>")
		return
	}

	for _, arg := range args {
		file, err := os.ReadFile(arg)
		if err != nil {
			fmt.Printf("Error loading file: %s\n", err)
			continue
		}
		input := string(file)
		fmt.Println(">> Executing:", arg)
		fmt.Println("-------------------------")

		interpreter, err := interpreter.New(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		t := time.Now()
		res := interpreter.Run()
		done := time.Since(t)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("-------------------------")
		fmt.Printf("Main Returned %v in: %v\n\n", res, done)
	}
}
