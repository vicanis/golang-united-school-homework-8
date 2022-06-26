package main

import (
	"errors"
	"flag"
)

type Arguments map[string]string

func parseArgs() Arguments {
	operation := flag.String("operation", "", "operation type")
	id := flag.String("id", "", "use id to delete")
	item := flag.String("item", "", "user item json")
	fileName := flag.String("fileName", "", "storage file name")

	flag.Parse()

	args := Arguments{
		"operation": "",
		"id":        "",
		"item":      "",
		"fileName":  "",
	}

	if operation != nil && *operation != "" {
		args["operation"] = *operation
	}

	if id != nil && *id != "" {
		args["id"] = *id
	}

	if item != nil && *item != "" {
		args["item"] = *item
	}

	if fileName != nil && *fileName != "" {
		args["fileName"] = *fileName
	}

	return args
}

var ErrorNoOperation = errors.New("-operation flag has to be specified")
