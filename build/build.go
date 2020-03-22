package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	file "../file"
	index "../index"
	"github.com/bbalet/stopwords"
)

// prepareTokens remove space literaral and stopwords from data string , splited and translates to lower
func prepareTokens(data string) []string {
	cleanSting := stopwords.CleanString(data, "en", true)
	tokens := strings.Fields(cleanSting)
	for i := range tokens {
		tokens[i] = strings.ToLower(tokens[i])
	}
	return tokens
}

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

	for i, f := range files {
		fmt.Println(f.Name())
		data, err := file.ReadFile(filepath.Join(os.Args[1], f.Name()))
		if err != nil {
			log.Fatalln(err)
			continue
		}
		tokens := prepareTokens(data)
		for _, token := range tokens {
			index.AddToken(maps, token, i)
		}
	}

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
