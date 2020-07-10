package main

import (
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
	// "github.com/gorilla/mux"
)

func main() {
	http.HandleFunc("/ivr/", ivr)

	// r := mux.NewRouter()
	// r.HandleFunc("/products/{key}", ProductHandler)
	// r.HandleFunc("/articles/{category}/", ArticlesCategoryHandler)
	// r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)

	port := ":8080"
	log.Printf("Listening port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func ivr(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("request: %#v\n", r)
	fmt.Printf("url: %#v\n", r.URL)
	// bs, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if err != nil {
	// 	log.Fatal("read body error")
	// }
	// fmt.Printf("body: %s\n", bs)

	fmt.Printf("a: %s, b: %s\n", r.FormValue("a"), r.FormValue("b"))
	w.Write([]byte(r.URL.Path))
}
