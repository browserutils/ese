package main

import (
	"fmt"
	"os"

	"github.com/browserutils/ese/parser"
)

type CommandHandler func(command string) bool

var (
	command_handlers []CommandHandler
	command_args     []string
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [--debug] <command> [args]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  catalog [--long_values] <file>\n")
	fmt.Fprintf(os.Stderr, "  dump [--limit N] <file> [table ...]\n")
	fmt.Fprintf(os.Stderr, "  page <file> <page_number>\n")
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 && args[0] == "--debug" {
		parser.Debug = true
		parser.DebugWalk = true
		args = args[1:]
	}

	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	command := args[0]
	command_args = args[1:]

	for _, command_handler := range command_handlers {
		if command_handler(command) {
			return
		}
	}

	usage()
	os.Exit(1)
}
