package main

import (
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/polisgo2020/search-KudinovKV/index"
)

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

	dataCh := make(chan []string)

	maps := index.NewInvertIndex(dataCh)

	for i, f := range files {
		wg.Add(1)
		go index.MakeBuild(os.Args[1], f, i, dataCh, wg)
	}

	wg.Wait()
	close(dataCh)

	maps.WriteResult(os.Args[2])
}
