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
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_url "net/url"
	"os"
	"strings"

	"github.com/dyatlov/go-opengraph/opengraph"
)

const appVersion = "1.0.1"

func main() {
	version := flag.Bool("v", false, "prints current opengraph version")
	url := flag.String("url", "", "fetch url and extract OpenGraph info from there")
	flag.Parse()

	if *version {
		fmt.Println(appVersion)
		return
	}

	// allow url to be provided without flag too, by default
	if *url == "" && flag.NArg() == 1 {
		*url = flag.Arg(0)
	}

	if *url != "" {
		u, err := _url.ParseRequestURI(*url)
		if err != nil {
			log.Fatalf("Error parsing url: %s\n", err)
		} else if !strings.HasPrefix(u.Scheme, "http") {
			log.Fatal(u.Scheme)
			log.Fatalf("URL should have http(s) protocol: %s\n", *url)
		}
	}

	var reader io.Reader

	if *url != "" {
		resp, err := http.Get(*url)
		if err != nil {
			log.Fatalf("Error while fetching url %s: %s", *url, err)
		}

		reader = resp.Body

		defer resp.Body.Close()
	} else {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			// pipe
			reader = bufio.NewReader(os.Stdin)
		} else {
			flag.Usage()
			return
		}
	}

	og := opengraph.NewOpenGraph()
	if err := og.ProcessHTML(reader); err != nil {
		log.Fatalf("Error processing html: %s", err)
	}

	output, _ := json.MarshalIndent(og, "", "  ")
	fmt.Println(string(output))
}
