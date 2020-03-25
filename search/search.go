package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	file "github.com/polisgo2020/search-KudinovKV/file"
	index "github.com/polisgo2020/search-KudinovKV/index"
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
	in = index.PrepareTokens(strings.Join(in, " "))

	data, err := file.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}

	maps := index.NewInvertIndex()
	listOfFiles := maps.ParseIndexFile(data)
	searchResult := maps.MakeSearch(in, listOfFiles)

	for i, elem := range searchResult {
		fmt.Println(i+1, " file got ", elem, " points !")
	}
}
