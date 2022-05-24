// This is a simple program utilising dyatlov/go-opengraph library.
// It outputs Open Graph data either by downloading and parsing html
// or by reading html directly from a pipe
//
// Examples:
//
// Download and parse html page:
// ./opengraph https://www.youtube.com/watch?v=yhoI42bdwU4
//
// Parse piped html
// curl https://www.youtube.com/watch?v=yhoI42bdwU4 | ./opengraph
//

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/dyatlov/go-opengraph/opengraph"
)

func printHelp() {
	fmt.Printf("Usage: %s <url>\n", os.Args[0])
	os.Exit(0)
}

func main() {
	var reader io.Reader

	if len(os.Args) == 2 {
		url := os.Args[1]
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error while fetching url %s: %s", url, err)
		}

		reader = resp.Body

		defer resp.Body.Close()
	} else if len(os.Args) == 1 {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			// pipe
			reader = bufio.NewReader(os.Stdin)
		} else {
			printHelp()
		}
	} else {
		printHelp()
	}

	og := opengraph.NewOpenGraph()
	if err := og.ProcessHTML(reader); err != nil {
		log.Fatalf("Error processing html: %s", err)
	}

	output, _ := json.MarshalIndent(og, "", "  ")
	fmt.Println(string(output))
}
