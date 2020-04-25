/*

Package index builds the inverse index for the entered document, save/load index to entered file, to search over the build index with thread-safe functions.

Usage

To create new empty index instance use NewInvertIndex function:

	i := index.NewInvertIndex()

which would create instance in thread-safe way starting internal channel listener to add new tokens.

To fill index call i.Add :

	i.Add("think" , "1.txt")

To save index use WriteResult with filename:

	i.WriteResult("output/filename")

To search in index file use Get:

	i.Get([]string{"hello", "world", "!!!"})

which builds Rate struct.

*/
package index

import (
	"github.com/bbalet/stopwords"
	"github.com/polisgo2020/search-KudinovKV/file"
	zl "github.com/rs/zerolog/log"
	"sort"
	"strings"
)

type Rate struct {
	FileName   string
	CountMatch int
}

type InvertIndex struct {
	index  map[string][]string
	dataCh chan []string
}

// Listener got tokens from channel and added to maps
func (i InvertIndex) Listener() {
	for input := range i.dataCh {
		i.Add(input[0], input[1])
	}
}

// NewInvertIndex return empty InvertIndex and start listener gorutine
func NewInvertIndex() *InvertIndex {
	i := InvertIndex{
		index:  map[string][]string{},
		dataCh: make(chan []string),
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

// ParseIndexFile added index in map
func (i InvertIndex) ParseIndexFile(data string) {
	datastrings := strings.Split(data, "\n")
	for _, correctstring := range datastrings {
		if correctstring == "" {
			break
		}
		keys := strings.Split(correctstring, ":")
		values := strings.Split(keys[1], ",")
		for _, value := range values {
			i.index[keys[0]] = append(i.index[keys[0]], value)
		}
	}
}

// MakeSearch find tokens in the map
func (i InvertIndex) Get(tokens []string) ([]Rate, error) {
	out := []Rate{}

	for _, token := range tokens {
		if _, ok := i.index[token]; ok {
			for _, fileName := range i.index[token] {
				flag := false
				for i, element := range out {
					if strings.EqualFold(element.FileName, fileName) {
						out[i].CountMatch += 1
						flag = true
						break
					}
				}
				if !flag {
					out = append(out, Rate{
						FileName:   fileName,
						CountMatch: 1,
					})
				}
			}
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CountMatch > out[j].CountMatch
	})
	return out, nil
}

// WriteResult write maps in file
func (i InvertIndex) WriteResult(outputFilename string) {
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

// Add add new token in index map
func (i InvertIndex) Add(token, fileName string) {
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

// Close all
func (i *InvertIndex) Close() {}

// PrepareTokens remove space literal and stop-words from data string , split and translates to lower
func PrepareTokens(data string) []string {
	cleanSting := stopwords.CleanString(data, "en", true)
	tokens := strings.Fields(cleanSting)
	for i := range tokens {
		tokens[i] = strings.ToLower(tokens[i])
	}
	return tokens
}
