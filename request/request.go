package request

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type RequestResult struct {
	Url            string
	ElapsedSeconds float64
	StatusCode     string
}

func sendRequest(c chan RequestResult, wg *sync.WaitGroup, urlString string) {
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

	c <- RequestResult{
		Url:            urlString,
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

func ProcessSubmittedUrls(submittedUrls []string, onResult func(requestResult RequestResult)) {
	// this function is the producer

	// the channel buffer will need to be, at least, the total number of submitted url parameters
	c := make(chan RequestResult, len(submittedUrls))
	wg := sync.WaitGroup{}

	defer wg.Wait()

	// total requests are the number of actual sent requests
	totalRequests := 0

	for _, urlString := range submittedUrls {
		parsedUrl, err := parseUrl(urlString)

		if err != nil {
			invalidUrlResult := RequestResult{
				Url:        urlString,
				StatusCode: "INVALID",
			}
			onResult(invalidUrlResult)
			continue
		}

		go sendRequest(c, &wg, parsedUrl)
		totalRequests++
		wg.Add(1)
	}

	totalResponses := 0

	for result := range c {
		totalResponses++
		onResult(result)

		if totalResponses == totalRequests {
			// closing the channel causes this loop to end
			close(c)
		}
	}
}
