package main

import (
	"bytes"
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
	router.HandleFunc("/cookieInit", cookieInitHandler)
	router.HandleFunc("/cookieGet", cookieGetHandler)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func firstHandle(w http.ResponseWriter, r *http.Request) {
	structJSON, err := json.Marshal(RequestFormat{Search: "yandex",
		Sites: []Site{"https://geekbrains.ru", "https://yandex.ru", "https://google.com", "https://habr.com", ""}})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp, err := http.Post("http://"+r.Host+"/search", "application/json", bytes.NewReader(structJSON))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "%s", body)
}

func searchHandle(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Error: body = nil", 500)
		return
	}
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

	sites, errsString := search(jsonStruct.Search, jsonStruct.Sites)
	if errsString != "" {
		http.Error(w, errsString, 500)
	}

	b, err := json.Marshal(sites)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func search(str string, sites []Site) ([]Site, string) {
	arr := make([]Site, 0, len(sites))
	errs := ""

	for _, site := range sites {
		if site != "" {
			resp, err := http.Get(string(site))
			if err != nil {
				errs += "Error GET: " + err.Error() + "\n"
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errs += "Error GET: " + err.Error() + "\n"
				continue
			}
			defer resp.Body.Close()

			if strings.Contains(string(body), str) {
				arr = append(arr, site)
			}
		}
	}

	return arr, errs
}

func cookieInitHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "GB",
		Value: "task",
	})
}

func cookieGetHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("GB")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "Name = %s, Value = %s", cookie.Name, cookie.Value)
}
