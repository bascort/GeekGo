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
	var err error

	for _, v := range sites {
		if v != "" {
			resp, err := http.Get(v)
			if err != nil {
				return arr, err
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return arr, err
			}

			defer resp.Body.Close()

			stringBody := string(body)
			i := strings.Count(stringBody, str)
			if i > 0 {
				arr = append(arr, v)
			}
		}
	}

	return arr, err
}
