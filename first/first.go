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
	"github.com/bbalet/stopwords"
)

type InvertIndex map[string][]int

// contains check element in int array
func contains(arr []int, element int) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}

// addToken add new token in map index
func addToken(index InvertIndex, token string, fileID int) {
	_, ok := index[token]
	b := contains((index)[token], fileID)
	if !ok || !b {
		index[token] = append(index[token], fileID)
		fmt.Println("Token : ", token)
		fmt.Println("Value: ", index[token])
		fmt.Println()
	}
}

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

	var maps InvertIndex
	maps = map[string][]int{}

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
			addToken(maps, token, i)
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
