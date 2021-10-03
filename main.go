package main

import (
	"flag"
	"fmt"
	"os"
	"webping/webping"
)

func main() {
	flag.Parse()

	submittedUrls := flag.Args()

	if len(submittedUrls) == 0 {
		fmt.Println("Usage: webping.exe <url1> <urls2> ... <urln>")
		os.Exit(1)
	}

	printMessage := func(webpingResult webping.WebpingResult) {
		formattedMessage := fmt.Sprintf("[%s] %s in %v seconds", webpingResult.StatusCode, webpingResult.RequestUrl, webpingResult.ElapsedSeconds)
		fmt.Println(formattedMessage)
	}

	webping.ProcessSubmittedUrls(submittedUrls, printMessage)
}
