package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Xkcd struct {
	Month      string `json:"month"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
	Num        int    `json:"num"`
}

var debug = true

func downloadAll() {
	fmt.Println("Getting ALL issues till date")
}

func getIssue(num int) error {
	var url string
	n := strconv.Itoa(num)

	if num == 0 { // just the latest issue
		url = "https://xkcd.com/info.0.json"
	} else {
		url = "https://xkcd.com/" + n + "/info.0.json"
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

	data, _ := ioutil.ReadAll(r.Body)

	if debug == true {
		fmt.Println("JSON data is:")
		fmt.Println(string(data))
	}
	var xkcdImg Xkcd
	json.Unmarshal(data, &xkcdImg)
	if n == "0" {
		n = strconv.Itoa(xkcdImg.Num)
	}

	imgName := "xkcd-" + n + "-" + path.Base(xkcdImg.Img)
	if debug == true {
		fmt.Println("Img Path: " + imgName)
		fmt.Println(xkcdImg.Alt)
	}

	fname := "xkcd-" + n + ".txt"
	fmt.Println("filename: " + fname)
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	txt, err := f.WriteString(xkcdImg.Alt)
	if err != nil {
		return err
	}
	if txt == 0 {
		fmt.Println("Filesize is 0")
	}
	f.Sync()
	defer f.Close()

	img, err := http.Get(xkcdImg.Img)
	if err != nil {
		return err
	}
	defer img.Body.Close()
	imgfile, err := os.Create(imgName)
	_, err = io.Copy(imgfile, img.Body)

	defer imgfile.Close()
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
