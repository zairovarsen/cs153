package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arsenzairov/cs153/parser"
	"github.com/arsenzairov/cs153/scanner"
)

func main() {
	if len(os.Args) >= 2 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			os.Exit(1)
		}
		err = runProgram(file, file.Name())
		if err != nil {
			os.Exit(1)
		}
	} else {
		runRepl(os.Stdin)
	}

}

func runProgram(file *os.File, name string) error {
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("Error getting the file info: %v\n", err)
	}
	estimatedSize := int(fileInfo.Size())

	scanner := scanner.NewScanner(file, name, estimatedSize)
	p := parser.NewParser(scanner)
	return p.Parse()
}

func runRepl(stdin io.Reader) {

	input := bufio.NewScanner(stdin)
	fmt.Println("Enter statements (type 'exit' to quit):")

	for {
		fmt.Print("> ")
		if !input.Scan() {
			break
		}

		line := input.Text()
		if line == "exit" {
			break
		}

		reader := strings.NewReader(line)
		scanner := scanner.NewScanner(reader, "main.go")

		p := parser.NewParser(scanner)
		err := p.Parse()
		if err != nil {
			os.Exit(1)
		}
	}

}
