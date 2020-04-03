package index

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bbalet/stopwords"
	"github.com/polisgo2020/search-KudinovKV/file"
)

type InvertIndex struct {
	index  map[string][]int
	dataCh chan []string
	mutex  *sync.Mutex
}

type Rate struct {
	fileID     int
	countMatch int
}

// GetRateFields return fields struct Rate
func (r Rate) GetRateFields() (fileID, countMatch int) {
	return r.fileID, r.countMatch
}

// Listener got tokens from cannel and added to maps
func (i InvertIndex) Listener() {
	defer i.mutex.Unlock()
	i.mutex.Lock()
	for input := range i.dataCh {
		token := input[0]
		k, _ := strconv.Atoi(input[1])
		i.addToken(token, k)
	}
}

// NewInvertIndex return empty InvertIndex and start listener gorutine
func NewInvertIndex() *InvertIndex {
	i := InvertIndex{
		index:  map[string][]int{},
		dataCh: make(chan []string),
		mutex:  &sync.Mutex{},
	}
	go i.Listener()
	return &i
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
func (i InvertIndex) ParseIndexFile(data string) []int {
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
			i.index[keys[0]] = append(i.index[keys[0]], number)
			if ok := Contains(listOfFiles, number); !ok {
				listOfFiles = append(listOfFiles, number)
			}
		}
	}
	sort.Ints(listOfFiles)
	return listOfFiles
}

// MakeSearch find in string tokens in the index map
func (i InvertIndex) MakeSearch(in []string, listOfFiles []int) []Rate {
	out := []Rate{}

	for k := range listOfFiles {
		count := 0
		for j := range in {
			if ok := Contains(i.index[in[j]], listOfFiles[k]); ok {
				count++
			}
		}
		out = append(out, Rate{fileID: k, countMatch: count})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].countMatch > out[j].countMatch
	})
	return out
}

// MakeBuild read files and added token in the cannel
func (i InvertIndex) MakeBuild(dirname string, f os.FileInfo, fileID int, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := file.ReadFile(filepath.Join(dirname, f.Name()))
	if err != nil {
		log.Fatalln(err)
	}
	tokens := PrepareTokens(data)
	for _, token := range tokens {
		i.dataCh <- []string{token, strconv.Itoa(fileID)}
	}
}

// WriteResult write maps in file
func (i InvertIndex) WriteResult(outputFilename string) {
	defer i.mutex.Unlock()
	close(i.dataCh)
	i.mutex.Lock()
	var resultString string
	for key, value := range i.index {
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

// addToken add new token in index map
func (i InvertIndex) addToken(token string, fileID int) {
	_, ok := i.index[token]
	b := Contains(i.index[token], fileID)
	if !ok || !b {
		i.index[token] = append(i.index[token], fileID)
		log.Println("Token : ", token)
		log.Println("Value: ", i.index[token])
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
