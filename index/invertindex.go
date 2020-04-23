/*

Package index builds the inverse index for the entered document, save/load index to entered file, to search over the build index with thread-safe functions.

Usage

To create new empty index instance use NewInvertIndex function:

	i := index.NewInvertIndex()

which would create instance in thread-safe way starting internal channel listener to add new tokens.

To fill index call MakeBuild with path to directory:

	i.MakeBuild("path/to/directory")

MakeBuild parse all files in entered directory, clear tokens, add them in the index .

To save index use WriteResult with filename:

	i.WriteResult("output/filename")

To search in index file use MakeSearch:

	i.MakeSearch("in/tokens", []string{"1.txt", "2.txt", "3.txt"})

which builds Rate struct.

*/
package index

import (
	"github.com/go-pg/pg/v9"
	"github.com/polisgo2020/search-KudinovKV/database"
	"log"
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
	wg     *sync.WaitGroup
}

// GetRateCount return countMatch struct Rate field
func (r Rate) GetRateCount() int {
	return r.countMatch
}

// GetRateName return fileName struct Rate field
func (r Rate) GetRateName() string {
	return r.fileName
}

// GetWg return WaitGroup struct InvertIndex field
func (i InvertIndex) GetWg() *sync.WaitGroup {
	return i.wg
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
		wg:     &sync.WaitGroup{},
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

// MakeSearchDB search tokens in database
func (i InvertIndex) MakeSearchDB(tokens []string, pg *pg.DB) []Rate {
	out := []Rate{}
	for _, token := range tokens {
		i, err := database.GetFiles(token, pg)
		if err != nil {
			zl.Debug().
				Msgf("Can't select in database , %v", err)
			return nil
		}
		for _, element := range i {
			flag := 0
			for i, value := range out {
				if strings.EqualFold(value.fileName, element.FileName) {
					out[i].countMatch += 1
					flag = 1
				}
			}
			if flag == 0 {
				out = append(out, Rate{
					fileName:   element.FileName,
					countMatch: 1,
				})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].countMatch > out[j].countMatch
	})
	return out
}

// MakeBuild read files and added token in the channel
func (i InvertIndex) MakeBuild(path string) {
	defer i.wg.Done()
	data, err := file.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	tokens := PrepareTokens(data)
	for _, token := range tokens {
		i.dataCh <- []string{token, path}
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

// WriteResultDB write maps in database
func (i InvertIndex) WriteResultDB(pg *pg.DB) {
	defer i.mutex.Unlock()
	close(i.dataCh)
	i.mutex.Lock()
	var newIndexDB []database.Index
	for token, fileNames := range i.index {
		for _, fileName := range fileNames {
			newIndexDB = append(newIndexDB, database.Index{
				FileName: fileName,
				Token:    token,
			})
		}
	}
	if err := database.AddIndex(newIndexDB, pg); err != nil {
		zl.Fatal().Err(err).
			Msg("Cant write in database")
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

// PrepareTokens remove space literal and stop-words from data string , split and translates to lower
func PrepareTokens(data string) []string {
	cleanSting := stopwords.CleanString(data, "en", true)
	tokens := strings.Fields(cleanSting)
	for i := range tokens {
		tokens[i] = strings.ToLower(tokens[i])
	}
	return tokens
}
