package index

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/bbalet/stopwords"
	"github.com/polisgo2020/search-KudinovKV/file"
	zl "github.com/rs/zerolog/log"
)

type Rate struct {
	fileName   string
	countMatch int
}

type InvertIndex struct {
	index  map[string][]string
	dataCh chan []string
	mutex  *sync.Mutex
}

// GetRateCount return count field struct Rate
func (r Rate) GetRateCount() int {
	return r.countMatch
}

// GetRateName return name field struct Rate
func (r Rate) GetRateName() string {
	return r.fileName
}

// Listener got tokens from channel and added to maps
func (i InvertIndex) Listener() {
	defer i.mutex.Unlock()
	i.mutex.Lock()
	for input := range i.dataCh {
		i.addToken(input[0], input[1])
	}
}

// NewInvertIndex return empty InvertIndex and start listener gorutine
func NewInvertIndex() *InvertIndex {
	i := InvertIndex{
		index:  map[string][]string{},
		dataCh: make(chan []string),
		mutex:  &sync.Mutex{},
	}
	go i.Listener()
	return &i
}

// Contains check element in int array
func Contains(arr []string, element string) bool {
	for _, a := range arr {
		if a == element {
			return true
		}
	}
	return false
}

// ParseIndexFile added index in map and return slice of files
func (i InvertIndex) ParseIndexFile(data string) []string {
	listOfFiles := []string{}
	datastrings := strings.Split(data, "\n")
	for _, correctstring := range datastrings {
		if correctstring == "" {
			break
		}
		keys := strings.Split(correctstring, ":")
		values := strings.Split(keys[1], ",")
		for _, value := range values {
			i.index[keys[0]] = append(i.index[keys[0]], value)
			if ok := Contains(listOfFiles, value); !ok {
				listOfFiles = append(listOfFiles, value)
			}
		}
	}
	return listOfFiles
}

// MakeSearch find in string tokens in the index map
func (i InvertIndex) MakeSearch(in, listOfFiles []string) []Rate {
	out := []Rate{}
	for k, filename := range listOfFiles {
		count := 0
		for j := range in {
			if ok := Contains(i.index[in[j]], listOfFiles[k]); ok {
				count++
			}
		}
		out = append(out, Rate{
			fileName:   filename,
			countMatch: count,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].countMatch > out[j].countMatch
	})
	return out
}

// MakeBuild read files and added token in the channel
func (i InvertIndex) MakeBuild(dirname string, f os.FileInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := file.ReadFile(filepath.Join(dirname, f.Name()))
	if err != nil {
		log.Fatalln(err)
	}
	tokens := PrepareTokens(data)
	for _, token := range tokens {
		i.dataCh <- []string{token, f.Name()}
	}
}

// WriteResult write maps in file
func (i InvertIndex) WriteResult(outputFilename string) {
	defer i.mutex.Unlock()
	close(i.dataCh)
	i.mutex.Lock()
	var resultString string
	for key, value := range i.index {
		var fileNames []string
		for _, filename := range value {
			fileNames = append(fileNames, filename)
		}
		resultString += key + ":" + strings.Join(fileNames, ",") + "\n"
	}
	err := file.WriteFile(resultString, outputFilename)
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Cant write in index file")
	}
}

// addToken add new token in index map
func (i InvertIndex) addToken(token, fileName string) {
	_, ok := i.index[token]
	b := Contains(i.index[token], fileName)
	if !ok || !b {
		i.index[token] = append(i.index[token], fileName)
		zl.Debug().
			Msgf("Token: %s", token)
		zl.Debug().
			Msgf("Value: %v", i.index[token])
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
