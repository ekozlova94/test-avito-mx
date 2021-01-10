package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/load-file", loadFile)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

func loadFile(w http.ResponseWriter, _ *http.Request) {
	file, err := ioutil.ReadFile("./example-goods.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	if _, errWrite := w.Write(file); errWrite != nil {
		log.Fatal(errWrite)
	}
}
