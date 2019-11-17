package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"
)

var rootCommand = &cobra.Command{
	Use:   "nasmfmt",
	Short: "The tool to format NASM file according to myself format.",
	Long: `The NASM file formatter to format the .nasm files for my own daily work.

This simple program will read the nasm file as stdin and output the formatted script
as stdout.`,
	Version: "0.0.1",
	Run:     process,
}

func process(command *cobra.Command, args []string) {
	// Directly read the string from stdin.
	reader := bufio.NewReader(os.Stdin)
	indent := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Reading error from stdin: ", err)
			os.Exit(1)
		}

		tokens := strings.Fields(line)
		formatted := strings.Join(tokens, " ")
		if len(tokens) == 1 && tokens[0][len(tokens[0])-1] == ':' {
			fmt.Println(formatted)
			indent = 4
		} else {
			for i := 0; i < indent; i++ {
				fmt.Print(" ")
			}
			fmt.Println(formatted)
		}
	}
}

func main() {
	err := rootCommand.Execute()
	if err != nil {
		log.Fatal("Error occurs when executing conversion.", err)
		os.Exit(1)
	}
}
