package main

import (
	"fmt"
	"os"

	"github.com/browserutils/ese/parser"
)

var (
	catalog_command = "catalog"
)

func doCatalog() error {
	long_values := false
	args := command_args

	if len(args) > 0 && args[0] == "--long_values" {
		long_values = true
		args = args[1:]
	}

	if len(args) != 1 {
		return fmt.Errorf("usage: %s catalog [--long_values] <file>", os.Args[0])
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

	catalog, err := parser.ReadCatalog(ese_ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s", catalog.Dump(parser.DumpOptions{
		Indexes:         true,
		Tables:          true,
		LongValueTables: long_values,
	}))
	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case catalog_command:
			err := doCatalog()
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
