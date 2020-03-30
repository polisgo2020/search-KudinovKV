package index

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/bbalet/stopwords"
	"github.com/polisgo2020/search-KudinovKV/file"
)

type InvertIndex map[string][]int

// Listener got tokens from cannel and added to maps
func (index InvertIndex) Listener(dataCh <-chan []string, mutex *sync.Mutex) {
	mutex.Lock()
	for input := range dataCh {
		token := input[0]
		i, _ := strconv.Atoi(input[1])
		index.AddToken(token, i)
	}
	mutex.Unlock()
}

// NewInvertIndex return empty InvertIndex
func NewInvertIndex() InvertIndex {
	index := InvertIndex{}
	return index
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

// ParseIndexFile added index in map and return slice of files
func (index InvertIndex) ParseIndexFile(data string) []int {
	var listOfFiles []int

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
	return listOfFiles
}

// MakeSearch find in string tokens in the index map
func (index InvertIndex) MakeSearch(in []string, listOfFiles []int) []int {
	var out []int
	var searchResult []int
	maxpoints := 0

	for i := range listOfFiles {
		count := 0
		for j := range in {
			if ok := Contains(index[in[j]], listOfFiles[i]); ok {
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
				searchResult = append(searchResult, out[j])
			}
		}
		i--
	}
	return searchResult
}

// MakeBuild read files and added token in the cannel
func MakeBuild(dirname string, f os.FileInfo, i int, out chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := file.ReadFile(filepath.Join(dirname, f.Name()))
	if err != nil {
		log.Fatalln(err)
	}
	tokens := PrepareTokens(data)
	for _, token := range tokens {
		info := []string{token, strconv.Itoa(i)}
		out <- info
	}
}

// WriteResult write maps in file
func (index InvertIndex) WriteResult(outputFilename string) {
	var resultString string
	for key, value := range index {
		var IDs []string

		for _, i := range value {
			IDs = append(IDs, strconv.Itoa(i))
		}
		resultString += key + ":" + strings.Join(IDs, ",") + "\n"
	}
	err := file.WriteFile(resultString, outputFilename)
	if err != nil {
		log.Fatalln(err)
	}
}

// AddToken add new token in map index
func (index InvertIndex) AddToken(token string, fileID int) {
	_, ok := index[token]
	b := Contains((index)[token], fileID)
	if !ok || !b {
		index[token] = append(index[token], fileID)
		log.Println("Token : ", token)
		log.Println("Value: ", index[token])
		log.Println()
	}
}

// PrepareTokens remove space literaral and stopwords from data string , splited and translates to lower
func PrepareTokens(data string) []string {
	cleanSting := stopwords.CleanString(data, "en", true)
	tokens := strings.Fields(cleanSting)
	for i := range tokens {
		tokens[i] = strings.ToLower(tokens[i])
	}
	return tokens
}
