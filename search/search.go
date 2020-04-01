package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/polisgo2020/search-KudinovKV/file"
	"github.com/polisgo2020/search-KudinovKV/index"
)

var (
	maps        index.InvertIndex
	listOfFiles []int
)

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "URL:", r.URL.String())

	tokens := r.FormValue("tokens")
	if tokens == "" {
		fmt.Fprintln(w, "Incorrect request! Example : ip:port/?param=tokens to search with space.")
		return
	}

	in := index.PrepareTokens(tokens)
	searchResult := maps.MakeSearch(in, listOfFiles)

	for i, elem := range searchResult {
		fmt.Fprintln(w, i+1, " file got ", elem, " points !")
	}
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalln("Invalid number of arguments. Example of call: /path/to/index/file ip-address port")
	}

	ip := os.Args[2]
	port := os.Args[3]
	mux := http.NewServeMux()

	data, err := file.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	maps = index.NewInvertIndex()
	listOfFiles = maps.ParseIndexFile(data)

	mux.HandleFunc("/", handler)

	server := http.Server{
		Addr:         ip + ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("starting server at ", ip, ":", port)
	server.ListenAndServe()
}
