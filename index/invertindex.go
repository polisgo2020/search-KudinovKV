package index

import (
	"fmt"
	"strconv"
	"strings"
)

type InvertIndex map[string][]int

// NewInvertIndex return empty InvertIndex
func NewInvertIndex() InvertIndex {
	return map[string][]int{}
}

// Contains check element in int array
func Contains(arr []int, element int) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}

// ParseFile return slice of files and map
func ParseFile(data string) ([]int, InvertIndex) {
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
			if ok := Contains(listOfFiles, number); !ok {
				listOfFiles = append(listOfFiles, number)
			}
		}
	}
	return listOfFiles, index
}

// MakeSearch find In string tokens in the index map
func MakeSearch(In []string, listOfFiles []int, index InvertIndex) {
	var out []int
	maxpoints := 0

	for i := range listOfFiles {
		count := 0
		for j := range In {
			if ok := Contains(index[In[j]], listOfFiles[i]); ok {
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

// AddToken add new token in map index
func AddToken(index InvertIndex, token string, fileID int) {
	_, ok := index[token]
	b := Contains((index)[token], fileID)
	if !ok || !b {
		index[token] = append(index[token], fileID)
		fmt.Println("Token : ", token)
		fmt.Println("Value: ", index[token])
		fmt.Println()
	}
}
