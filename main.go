package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var wg sync.WaitGroup
var mut sync.Mutex

func sendRequest(urlString string) {
	defer wg.Done()

	res, err := http.Get(urlString)
	resultMessage := ""

	if err != nil {
		resultMessage = fmt.Sprintf("[unknown host] %s", urlString)
	} else {
		resultMessage = fmt.Sprintf("[%d] %s", res.StatusCode, urlString)
	}

	mut.Lock()
	defer mut.Unlock()

	fmt.Println(resultMessage)
}

func parseUrl(urlString string) (string, error) {
	u, err := url.Parse(urlString)

	if err != nil { return "", err }

	if u.Scheme == "" {
		return "https://" + urlString, nil
	}

	return urlString, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: go run main.go <url1> <urls2> ... <urln>")
	}

	for _, urlString := range os.Args[1:] {
		parsedUrl, err := parseUrl(urlString)

		if err != nil {
			fmt.Println(fmt.Sprintf("[invalid url] %s", urlString))
			continue
		}

		go sendRequest(parsedUrl)

		wg.Add(1)
	}

	wg.Wait()
}
