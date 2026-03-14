package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/browserutils/ese/parser"
)

var (
	page_command = "page"
)

func doPage() error {
	args := command_args
	if len(args) != 2 {
		return fmt.Errorf("usage: %s page <file> <page_number>", os.Args[0])
	}

	fd, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fd.Close()

	ese_ctx, err := parser.NewESEContext(fd)
	if err != nil {
		return err
	}

	page_number, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}

	parser.DumpPage(ese_ctx, page_number)
	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case page_command:
			err := doPage()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		default:
			return false
		}
		return true
	})
}
