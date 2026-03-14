package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/browserutils/ese/ordereddict"
	"github.com/browserutils/ese/parser"
)

var (
	dump_command = "dump"

	STOP_ERROR = errors.New("Stop")
)

func marshalOrderedValue(v interface{}) ([]byte, error) {
	switch t := v.(type) {
	case *ordereddict.Dict:
		return marshalOrderedDict(t)
	default:
		return json.Marshal(v)
	}
}

func marshalOrderedDict(row *ordereddict.Dict) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')

	for idx, key := range row.Keys() {
		if idx > 0 {
			buf.WriteByte(',')
		}

		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		value, _ := row.Get(key)
		valueBytes, err := marshalOrderedValue(value)
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func doDump() error {
	args := command_args
	limit := 0
	file := ""
	tables := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == "--limit" && i+1 < len(args) {
			var err error
			limit, err = strconv.Atoi(args[i+1])
			if err != nil {
				return err
			}
			i++
			continue
		}

		if file == "" {
			file = args[i]
		} else {
			tables = append(tables, args[i])
		}
	}

	if file == "" {
		return fmt.Errorf("usage: %s dump [--limit N] <file> [table ...]", os.Args[0])
	}

	fd, err := os.Open(file)
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

	if len(tables) == 0 {
		tables = catalog.Tables.Keys()
	}

	for _, t := range tables {
		count := 0

		err = catalog.DumpTable(t, func(row *ordereddict.Dict) error {
			serialized, err := marshalOrderedDict(row)
			if err != nil {
				return err
			}

			count++
			fmt.Printf("%v\n", string(serialized))
			if limit > 0 &&
				count >= limit {
				return STOP_ERROR
			}

			return nil
		})
	}

	if err == STOP_ERROR {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func init() {
	command_handlers = append(command_handlers, func(command string) bool {
		switch command {
		case dump_command:
			err := doDump()
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
