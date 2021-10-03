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

	webping.ProcessSubmittedUrls(submittedUrls)
}
