package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type RequestFormat struct {
	Search string `json: "search"`
	Sites  []Site `json: "sites"`
}

type Site string

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", firstHandle)
	router.HandleFunc("/search", searchHandle)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func firstHandle(w http.ResponseWriter, r *http.Request) {
	sructJSON := `{"search":"yandex",
	"sites":[
		"https://geekbrains.ru",
		"https://yandex.ru",
		"https://google.com",
		"https://habr.com",
		""
		]
	}`

	resp, err := http.Post("http://"+r.Host+"/search", "application/json", strings.NewReader(sructJSON))

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer resp.Body.Close()

}

func searchHandle(w http.ResponseWriter, r *http.Request) {
	resp, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	jsonStruct := new(RequestFormat)

	err = json.Unmarshal(resp, jsonStruct)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	site, err := search(jsonStruct.Search, jsonStruct.Sites)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	b, err := json.Marshal(site)
	fmt.Fprintf(w, "%s", b)
}

func search(str string, sites []Site) ([]Site, error) {
	arr := make([]Site, 0, len(sites))

	for _, site := range sites {
		if site != "" {
			resp, err := http.Get(fmt.Sprintf("%s", site))
			if err != nil {
				return nil, err
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()

			i := strings.Contains(string(body), str)
			if i {
				arr = append(arr, site)
			}
		}
	}

	return arr, nil
}
