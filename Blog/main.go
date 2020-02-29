package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var dir string

    flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/views/")))
	//router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/posts/{id:[0-9]+}", postsHandler).Methods("GET")
	http.Handle("/", router)

	fmt.Println("Server is listing...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.vars(r)
	id := vars["id"]

	if id
}
