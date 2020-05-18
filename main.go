package main

import (
	//	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

func downloadAll() {
	fmt.Println("Getting ALL issues till date")
}

func getIssue(num int) error {
	var url string
	if num == 0 { // just the latest issue
		url = "https://xkcd.com/info.0.json"
	} else {
		url = "https://xkcd.com/" + strconv.Itoa(num) + "/info.0.json"
	}

	fmt.Println("Fetching url: " + url)

	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	out, err := os.Create(path.Base(r.Request.URL.String()))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r.Body)
	return err
}

func main() {

	var (
		all      bool
		specific int
		xrange   int
		help     bool
		version  bool
	)

	flag.BoolVar(&all, "a", false, "Download all")
	flag.IntVar(&specific, "n", 0, "Download specific number")
	flag.IntVar(&xrange, "r", 0, "Download a range From To")
	flag.BoolVar(&help, "h", false, "Display help")
	flag.BoolVar(&version, "v", false, "Display version")
	flag.Parse()

	if len(os.Args) == 1 {
		getIssue(0)
	} else if help == true {
		flag.PrintDefaults()
		os.Exit(1)
	} else if version == true {
		fmt.Println("XKCD Downloader v0.1 ")
	} else if all == true {
		downloadAll()
	} else if specific != 0 {
		getIssue(specific)
	}

}
