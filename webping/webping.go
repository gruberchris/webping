package webping

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type WebpingResult struct {
	RequestUrl     string
	ElapsedSeconds float64
	StatusCode     string
}

func sendRequest(c chan WebpingResult, wg *sync.WaitGroup, urlString string) {
	// this function is a consumer
	defer wg.Done()

	startTime := time.Now()
	res, err := http.Get(urlString)
	elapsedTime := time.Since(startTime)
	statusCode := ""

	if err != nil {
		if match, _ := regexp.MatchString(".*lookup.*", err.Error()); match {
			statusCode = "UNKNOWN HOST"
		} else if match, _ := regexp.MatchString(".*connection.*.refused.*", err.Error()); match {
			statusCode = "CONN REFUSED"
		} else {
			statusCode = "NET ERROR"
		}
	} else {
		statusCode = strconv.Itoa(res.StatusCode)
	}

	c <- WebpingResult{
		RequestUrl:     urlString,
		ElapsedSeconds: elapsedTime.Seconds(),
		StatusCode:     statusCode,
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

func ProcessSubmittedUrls(submittedUrls []string, outMessage func(webpingResult WebpingResult)) {
	// this function is the producer

	// the channel buffer will need to be, at least, the total number of submitted url parameters
	c := make(chan WebpingResult, len(submittedUrls))
	wg := sync.WaitGroup{}

	defer wg.Wait()

	// total requests are the number of actual sent requests
	totalRequests := 0

	for _, urlString := range submittedUrls {
		parsedUrl, err := parseUrl(urlString)

		if err != nil {
			invalidUrlResult := WebpingResult{
				RequestUrl: urlString,
				StatusCode: "INVALID",
			}
			outMessage(invalidUrlResult)
			continue
		}

		go sendRequest(c, &wg, parsedUrl)
		totalRequests++
		wg.Add(1)
	}

	totalResponses := 0

	for webpingResult := range c {
		totalResponses++
		outMessage(webpingResult)

		if totalResponses == totalRequests {
			// closing the channel causes this loop to end
			close(c)
		}
	}
}
