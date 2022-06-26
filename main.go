package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]

	if operation == "" {
		return ErrorNoOperation
	}

	fileName := args["fileName"]

	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	userList, err := loadUserList(fileName)
	if err != nil {
		return err
	}

	switch operation {
	case "list":
		data, err := userList.Dump()
		if err != nil {
			return err
		}

		if len(data) > 0 {
			fmt.Fprint(writer, string(data))
		}

	case "add":
		item := args["item"]

		if item == "" {
			return errors.New("-item flag has to be specified")
		}

		user, err := parseUser(item)
		if err != nil {
			return err
		}

		err = userList.Add(user)
		if err != nil {
			fmt.Fprint(writer, err)
			return nil
		}

		userList.saveUserList(fileName)

	case "remove":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		}

		if err = userList.Remove(args["id"]); err != nil {
			fmt.Fprint(writer, err)
			return nil
		}

		userList.saveUserList(fileName)

	case "findById":
		if args["id"] == "" {
			return errors.New("-id flag has to be specified")
		}

		if !userList.Has(args["id"]) {
			fmt.Fprint(writer, "")
		}

		var data []byte

		data, err = userList.DumpItem(args["id"])
		if err != nil {
			data = []byte("")
		}

		fmt.Fprint(writer, string(data))

	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
