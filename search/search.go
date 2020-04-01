package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/polisgo2020/search-KudinovKV/file"
	"github.com/polisgo2020/search-KudinovKV/index"
)

var (
	maps        index.InvertIndex
	listOfFiles []int
	topPage     = string(`
<!DOCTYPE html>
<html lang="en">
<head>
	<title>Contact V4</title>
	<meta charset="UTF-8">
	</head>
<style>
* {
	margin: 0px; 
	padding: 0px; 
	box-sizing: border-box;
}

body, html {
	height: 100%;
	font-family: Poppins-Regular, sans-serif;
}
.container-contact100 {
  width: 100%;  
  min-height: 100vh;
  display: -webkit-box;
  display: -webkit-flex;
  display: -moz-box;
  display: -ms-flexbox;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  padding: 15px;
  background: #a64bf4;
  background: -webkit-linear-gradient(45deg, #00dbde, #fc00ff);
  background: -o-linear-gradient(45deg, #00dbde, #fc00ff);
  background: -moz-linear-gradient(45deg, #00dbde, #fc00ff);
  background: linear-gradient(45deg, #00dbde, #fc00ff);
  
}

.wrap-contact100 {
  width: 50%;
  background: #c0cff7;
  border-radius: 10px;
  overflow: hidden;
  padding: 42px 55px 45px 55px;
}
.contact100-form-title {
  display: block;
  font-size: 39px;
  color: #333333;
  line-height: 1.2;
  text-align: center;
}
.rounded {
	counter-reset: li; 
	list-style: none; 
	font: 14px "Trebuchet MS", "Lucida Sans";
	padding: 0;
	text-shadow: 0 1px 0 rgba(255,255,255,.5);
}
.rounded a {
	position: relative;
	display: block;
	padding: .4em .4em .4em 2em;
	margin: .5em 0;
	background: #5b8bea;
	color: #444;
	text-decoration: none;
	border-radius: .3em;
	transition: .3s ease-out;
}
.rounded a:hover {background: #E9E4E0;}
.rounded a:hover:before {transform: rotate(360deg);}
.rounded a:before {
	content: counter(li);
	counter-increment: li;
	position: absolute;
	left: -1.3em;
	top: 50%;
	margin-top: -1.3em;
	background: #8FD4C1;
	height: 2em;
	width: 2em;
	line-height: 2em;
	border: .3em solid white;
	text-align: center;
	font-weight: bold;
	border-radius: 2em;
	transition: all .3s ease-out;
}
</style>
<body>
	<div class="container-contact100">
		<div class="wrap-contact100">
			<span class="contact100-form-title">
				<ol class="rounded">`)
	bottomPage = string(`
				</ol>
			</span>
		</div>
	</div>
</body>
</html>
`)
)

func handler(w http.ResponseWriter, r *http.Request) {
	tokens := r.URL.Query().Get("tokens")
	if tokens == "" {
		fmt.Fprintln(w, topPage+"Incorrect request!<br>Example: <i>ip:port/?param=tokens to search with space</i>"+bottomPage)
		return
	}

	in := index.PrepareTokens(tokens)
	searchResult := maps.MakeSearch(in, listOfFiles)
	resultString := string("")
	for _, elem := range searchResult {
		resultString += "<li><a href=\"#\"> file got " + strconv.Itoa(int(elem)) + " points !</a></li>"
	}
	fmt.Fprintln(w, topPage+resultString+bottomPage)
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
