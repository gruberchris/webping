package webping

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"
)

type webpingResult struct {
	Message        string
	ElapsedSeconds float64
}

func sendRequest(c chan webpingResult, wg *sync.WaitGroup, urlString string) {
	// this function is a consumer
	defer wg.Done()

	startTime := time.Now()
	res, err := http.Get(urlString)
	elapsedTime := time.Since(startTime)
	resultMessage := ""

	if err != nil {
		if match, _ := regexp.MatchString(".*lookup.*", err.Error()); match {
			resultMessage = fmt.Sprintf("[UNKNOWN HOST] %s", urlString)
		} else if match, _ := regexp.MatchString(".*connection.*.refused.*", err.Error()); match {
			resultMessage = fmt.Sprintf("[CONN REFUSED] %s", urlString)
		} else {
			fmt.Println(err)
			resultMessage = fmt.Sprintf("[NET ERROR] %s", urlString)
		}
	} else {
		resultMessage = fmt.Sprintf("[%d] %s", res.StatusCode, urlString)
	}

	c <- webpingResult{
		Message:        resultMessage,
		ElapsedSeconds: elapsedTime.Seconds(),
	}
}

func parseUrl(urlString string) (string, error) {
	u, err := url.ParseRequestURI(urlString)

	if err != nil || u.Host == "" {
		u, err := url.ParseRequestURI("https://" + urlString)

		if err != nil {
			return "", err
		}

		return u.Scheme + "://" + u.Host, nil
	}

	return u.Scheme + "://" + u.Host, nil
}

func ProcessSubmittedUrls(submittedUrls []string) {
	// this function is the producer

	// the channel buffer will need to be, at least, the total number of submitted url parameters
	c := make(chan webpingResult, len(submittedUrls))
	wg := sync.WaitGroup{}

	defer wg.Wait()

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

	totalResponses := 0

	for webpingResult := range c {
		totalResponses++
		formattedMessage := fmt.Sprintf("%s in %v seconds", webpingResult.Message, webpingResult.ElapsedSeconds)
		fmt.Println(formattedMessage)

		if totalResponses == totalRequests {
			// closing the channel causes this loop to end
			close(c)
		}
	}

	// fmt.Println("Processed " + strconv.Itoa(totalResponses) + " of " + strconv.Itoa(totalRequests) + " responses.")
}
