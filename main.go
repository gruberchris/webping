package main

import (
	"flag"
	"fmt"
	"github.com/gruberchris/webping/request"
	"os"
)

func main() {
	flag.Parse()

	submittedUrls := flag.Args()

	if len(submittedUrls) == 0 {
		fmt.Println("Usage: webping.exe <url1> <urls2> ... <urln>")
		os.Exit(1)
	}

	printMessage := func(webpingResult request.RequestResult) {
		formattedMessage := fmt.Sprintf("[%s] %s in %v seconds", webpingResult.StatusCode, webpingResult.Url, webpingResult.ElapsedSeconds)
		fmt.Println(formattedMessage)
	}

	request.ProcessSubmittedUrls(submittedUrls, printMessage)
}
