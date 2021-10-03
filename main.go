package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type WebpingResult struct {
	Message        string
	ElapsedSeconds float64
}

func sendRequest(c chan WebpingResult, wg *sync.WaitGroup, urlString string) {
	// this function is a consumer
	defer wg.Done()

	startTime := time.Now()

	res, err := http.Get(urlString)

	elapsedTime := time.Since(startTime)
	resultMessage := ""

	if err != nil {
		resultMessage = fmt.Sprintf("[unknown host] %s", urlString)
	} else {
		resultMessage = fmt.Sprintf("[%d] %s", res.StatusCode, urlString)
	}

	c <- WebpingResult{
		Message:        resultMessage,
		ElapsedSeconds: elapsedTime.Seconds(),
	}
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
	// this function is the producer

	// the channel buffer will need to be, at least, the total number of submitted url parameters
	c := make(chan WebpingResult, len(submittedUrls))
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

	for i := 1; i < totalRequests; i++ {
		webpingResult := <-c

		formattedMessage := fmt.Sprintf("%s in %v seconds", webpingResult.Message, webpingResult.ElapsedSeconds)
		fmt.Println(formattedMessage)

		// must close the channel when finished
		if i == totalRequests {
			close(c)
			break
		}
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
