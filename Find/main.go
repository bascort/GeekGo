package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	srt := "yandex"
	sites := []string{
		"https://geekbrains.ru",
		"https://yandex.ru",
		"https://google.com",
		"https://habr.com",
		"",
	}

	arr, err := search(srt, sites)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range arr {
		fmt.Printf("Site - %s\n", v)
	}
}

func search(str string, sites []string) ([]string, error) {
	arr := make([]string, 0, len(sites))

	for _, site := range sites {
		if site != "" {
			resp, err := http.Get(site)
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
