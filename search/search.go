package main

import (
	"log"
	"os"

	file "../file"
	index "../index"
)

// parseArgs return slice of string with args
func parseArgs() []string {
	var in []string

	for i := range os.Args {
		if i == 0 || i == 1 {
			continue
		}
		in = append(in, os.Args[i])
	}
	return in
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Invalid number of arguments. Example of call: /path/to/index/file tokens-to-search")
	}

	in := parseArgs()

	data, err := file.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}

	listOfFiles, maps := index.ParseFile(data)
	index.MakeSearch(in, listOfFiles, maps)
}
