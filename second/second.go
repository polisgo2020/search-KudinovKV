package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	file "../file"
)

type InvertIndex map[string][]int

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

// parseFile return slice of files and map
func parseFile(data string) ([]int, InvertIndex) {
	var listOfFiles []int
	var index InvertIndex
	index = map[string][]int{}

	datastrings := strings.Split(data, "\n")
	for _, correctstring := range datastrings {
		if correctstring == "" {
			break
		}
		keys := strings.Split(correctstring, ":")
		values := strings.Split(keys[1], ",")
		for _, value := range values {
			number, _ := strconv.Atoi(value)
			index[keys[0]] = append(index[keys[0]], number)
			if ok := contains(listOfFiles, number); !ok {
				listOfFiles = append(listOfFiles, number)
			}
		}
	}
	return listOfFiles, index
}

// contains check element in int array
func contains(arr []int, element int) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}

// makeSearch find In string tokens in the index map
func makeSearch(In []string, listOfFiles []int, index InvertIndex) {
	var out []int
	maxpoints := 0

	for i := range listOfFiles {
		count := 0
		for j := range In {
			if ok := contains(index[In[j]], listOfFiles[i]); ok {
				count++
			}
		}
		if count > maxpoints {
			maxpoints = count
		}
		out = append(out, count)
	}
	i := maxpoints
	for i != -1 {
		for j := range out {
			if out[j] == i {
				fmt.Println(j+1, " file got ", out[j], " points !")
			}
		}
		i--
	}
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

	listOfFiles, maps := parseFile(data)
	makeSearch(in, listOfFiles, maps)
}
