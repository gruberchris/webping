package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
)

func sendRequest(c chan string, wg *sync.WaitGroup, urlString string) {
	defer wg.Done()

	res, err := http.Get(urlString)
	resultMessage := ""

	if err != nil {
		resultMessage = fmt.Sprintf("[unknown host] %s", urlString)
	} else {
		resultMessage = fmt.Sprintf("[%d] %s", res.StatusCode, urlString)
	}

	c <- resultMessage
}

func parseUrl(urlString string) (string, error) {
	u, err := url.Parse(urlString)

	if err != nil {
		return "", err
	}

	if u.Scheme == "" {
		return "https://" + urlString, nil
	}

	return urlString, nil
}

func processSubmittedUrls(submittedUrls []string) {
	// the channel buffer will need to be, at least, the total number of submitted url parameters
	channelBufferLength := len(submittedUrls)

	c := make(chan string, channelBufferLength)
	wg := sync.WaitGroup{}

	// total requests are the number of actual sent requests
	totalRequests := 0

	for _, urlString := range submittedUrls {
		parsedUrl, err := parseUrl(urlString)

		if err != nil {
			fmt.Println(fmt.Sprintf("[invalid url] %s", urlString))
			continue
		}

		go sendRequest(c, &wg, parsedUrl)
		totalRequests++
		wg.Add(1)
	}

	i := 1

	for resultMessage := range c {
		fmt.Println(resultMessage)

		// must close the channel to exit this loop
		if i == totalRequests {
			close(c)
		}

		i++
	}

	wg.Wait()
}

func main() {
	flag.Parse()

	submittedUrls := flag.Args()

	if len(submittedUrls) == 0 {
		fmt.Println("Usage: webping.exe <url1> <urls2> ... <urln>")
		os.Exit(1)
	}

	processSubmittedUrls(submittedUrls)
}
