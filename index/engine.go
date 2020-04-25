package index

import (
	"github.com/polisgo2020/search-KudinovKV/file"
	zl "github.com/rs/zerolog/log"
	"sync"
)

// IndexEngine is the interface for the data storage object
type InvertIndexEngine interface {
	Add(token, fileName string)
	Get(tokens []string) ([]Rate, error)
	Close()
}

// Index use Engine to storage documents and find tokens
type Index struct {
	Engine InvertIndexEngine
	DataCh chan []string
	Wg     *sync.WaitGroup
	Mutex  *sync.Mutex
}

// listener listen chanel and use Engine to write data
func (i Index) listener() {
	defer i.Mutex.Unlock()
	i.Mutex.Lock()
	for input := range i.DataCh {
		i.Engine.Add(input[0], input[1])
	}
}

// NewIndex create new instance of engine and start listener
func NewIndex(engine InvertIndexEngine) *Index {
	i := Index{
		Engine: engine,
		DataCh: make(chan []string),
		Wg:     &sync.WaitGroup{},
		Mutex:  &sync.Mutex{},
	}
	go i.listener()
	return &i
}

// MakeBuild read files and added token in the channel
func (i Index) MakeBuild(path string) {
	defer i.Wg.Done()
	data, err := file.ReadFile(path)
	if err != nil {
		zl.Fatal().Err(err).
			Msg("Can't read file")
	}
	tokens := PrepareTokens(data)
	for _, token := range tokens {
		i.DataCh <- []string{token, path}
	}
}
