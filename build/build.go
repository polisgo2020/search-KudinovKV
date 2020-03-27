package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/polisgo2020/search-KudinovKV/index"
)

// listener got tokens from cannel and added to maps
func listener(dataCh <-chan []string, outputFilename string, maps index.InvertIndex, bufferMutex *sync.Mutex) {
	for input := range dataCh {
		token := input[0]
		i, _ := strconv.Atoi(input[1])
		bufferMutex.Lock()
		maps.AddToken(token, i)
		bufferMutex.Unlock()
	}
}

func main() {

	if len(os.Args) < 3 {
		log.Fatalln("Invalid number of arguments. Example of call: /path/to/files /path/to/output")
	}

	files, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatalln(err)
		return
	}
	wg := &sync.WaitGroup{}
	channelMutex := &sync.Mutex{}
	bufferMutex := &sync.Mutex{}
	maps := index.NewInvertIndex()

	dataCh := make(chan []string)
	defer close(dataCh)

	go listener(dataCh, os.Args[2], maps, bufferMutex)

	for i, f := range files {
		wg.Add(1)
		go index.MakeBuild(os.Args[1], f, i, dataCh, wg, channelMutex)
	}

	wg.Wait()

	bufferMutex.Lock()
	maps.WriteResult(os.Args[2])
	bufferMutex.Unlock()
}
