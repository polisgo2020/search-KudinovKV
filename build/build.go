package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	file "github.com/polisgo2020/search-KudinovKV/file"
	index "github.com/polisgo2020/search-KudinovKV/index"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatalln("Invalid number of arguments. Example of call: /path/to/files /path/to/output")
	}

	maps := index.NewInvertIndex()

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}

	maps.MakeBuild(os.Args[1], files)

	var resultString string
	for key, value := range maps {
		var IDs []string

		for _, i := range value {
			IDs = append(IDs, strconv.Itoa(i))
		}
		resultString += key + ":" + strings.Join(IDs, ",") + "\n"
	}

	err = file.WriteFile(resultString, os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
}
