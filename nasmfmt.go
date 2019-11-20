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
	Version: "0.3.0",
	Args:    cobra.ExactArgs(1),
	Run:     format,
}

func process(reader *bufio.Reader, buffer *bytes.Buffer) {
	indent := 0

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
		trimed := strings.TrimSpace(string(line))

		// Process the token splitting.
		var tokens []string
		inString := false
		current := ""
		for i, c := range trimed {
			if !inString && c == ';' {
				// In this case, we directly using suffix.
				current += string(trimed[i:])
				break
			} else if c == ' ' {
				if !inString && len(current) > 0 {
					tokens = append(tokens, current)
					current = ""
				} else if inString {
					current += string(c)
				}
			} else if c == '"' {
				if inString && trimed[i-1] == '\\' {
					// "...\" case.
				} else {
					inString = !inString
				}
				current += string(c)
			} else {
				current += string(c)
			}
		}
		if len(current) > 0 {
			tokens = append(tokens, current)
		}

		formatted := strings.Join(tokens, " ")
		if len(tokens) > 0 {
			current = tokens[len(tokens)-1]
		} else {
			current = ""
		}
		if len(formatted) == 0 {
			// Skip for empty block
		} else if !strings.HasPrefix(current, ";") && formatted[len(current)-1] == ':' {
			buffer.WriteString(formatted)
			indent = 4
		} else {
			for i := 0; i < indent; i++ {
				buffer.WriteString(" ")
			}
			buffer.WriteString(formatted)

			if len(tokens) > 0 && strings.HasPrefix(tokens[0], ";; section_end") {
				indent = 0
			}
		}
		buffer.WriteString("\n")
	}
}

func format(command *cobra.Command, args []string) {
	f, err := os.Open(args[0])
	if err != nil {
		log.Fatal(fmt.Sprintf("Error occurs when opening file: %s", args[0]), err)
		os.Exit(1)
	}

	reader := bufio.NewReader(f)
	buffer := bytes.Buffer{}
	process(reader, &buffer)
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
