package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("../src/index.html")
	w.Write(file)
}
func main() {
	http.HandleFunc("/", mainHandler)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir("../src/static")))
	http.Handle("/static/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
