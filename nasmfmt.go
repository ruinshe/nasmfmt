package main

import (
	"bufio"
	"bytes"
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
	Version: "0.2.1",
	Args:    cobra.ExactArgs(1),
	Run:     process,
}

func process(command *cobra.Command, args []string) {
	f, err := os.Open(args[0])
	if err != nil {
		log.Fatal(fmt.Sprintf("Error occurs when opening file: %s", args[0]), err)
		os.Exit(1)
	}

	// Directly read the string from stdin.
	reader := bufio.NewReader(f)
	indent := 0
	buffer := bytes.Buffer{}

	for {
		// Actually here second returned value isPrefix may be false,
		// but I ignore it for keep the code simple.
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("Reading error from stdin: ", err)
			os.Exit(1)
		}

		var tokens []string
		trimed := strings.TrimSpace(string(line))
		if strings.HasPrefix(trimed, ";;") {
			tokens = []string{trimed}
		} else {
			tokens = strings.Fields(trimed)
		}
		formatted := strings.Join(tokens, " ")
		if len(tokens) == 1 && tokens[0][len(tokens[0])-1] == ':' {
			buffer.WriteString(formatted)
			indent = 4
		} else {
			if len(formatted) > 0 {
				for i := 0; i < indent; i++ {
					buffer.WriteString(" ")
				}
			}
			buffer.WriteString(formatted)

			if len(tokens) >= 2 && tokens[0] == ";;" && tokens[1] == "section_end" {
				indent = 0
			}
		}
		buffer.WriteString("\n")
	}
	f.Close()

	f, err = os.Create(args[0])
	if err != nil {
		log.Fatal(fmt.Sprintf("Error occurs when creating file: %s", args[0]), err)
		os.Exit(1)
	}
	f.WriteString(buffer.String())
	f.Close()
}

func main() {
	rootCommand.PersistentFlags().BoolP("write", "w", true,
		"Write the file (excluded for bypass Emacs).")

	err := rootCommand.Execute()
	if err != nil {
		log.Fatal("Error occurs when executing conversion.", err)
		os.Exit(1)
	}
}
