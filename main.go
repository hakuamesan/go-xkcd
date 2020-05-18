package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
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

type Latest struct {
	Num int `json:"num"`
}

func getLatestCount() int {
	r, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		log.Fatal("Error getting latest count")
	}
	defer r.Body.Close()

	latest, _ := ioutil.ReadAll(r.Body)
	var l Latest
	json.Unmarshal(latest, &l)
	if debug == true {
		log.Println("Total count is : " + strconv.Itoa(l.Num))
	}
	totalCount = l.Num
	return l.Num
}

var debug = true
var totalCount = 0

func getComicRange(from int, to int) {
	log.Println("Getting ALL issues till date")
	log.Println("Total counts is ", to)

	for i := from; i <= to; i++ {
		getComic(i)
	}
}

func getComic(num int) error {
	var url string
	n := strconv.Itoa(num)

	if num == 0 { // just the latest issue
		url = "https://xkcd.com/info.0.json"
	} else {
		url = "https://xkcd.com/" + n + "/info.0.json"
	}

	log.Println("Fetching url: " + url)

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

	/*
		if debug == true {
			log.Println("JSON data is:")
			log.Println(string(data))
		}
	*/

	var xkcdImg Xkcd
	json.Unmarshal(data, &xkcdImg)
	if n == "0" {
		n = strconv.Itoa(xkcdImg.Num)
	}

	imgName := "xkcd-" + n + "-" + path.Base(xkcdImg.Img)
	if debug == true {
		log.Println("Img Path: " + imgName)
		//log.Println(xkcdImg.Alt)
	}

	fname := "xkcd-" + n + ".txt"
	log.Println("filename: " + fname)

	if _, err := os.Stat(fname); err == nil {
		log.Println("Txt description already saved. ")
	} else if os.IsNotExist(err) {

		f, err := os.Create(fname)
		if err != nil {
			return err
		}
		txt, err := f.WriteString(xkcdImg.Alt)
		if err != nil {
			return err
		}
		if txt == 0 {
			log.Fatal("Filesize is 0")
		}
		f.Sync()
		defer f.Close()
	}

	if _, err := os.Stat(imgName); err == nil {
		log.Println("File " + imgName + " exists. Skipping...")

	} else if os.IsNotExist(err) {

		img, err := http.Get(xkcdImg.Img)
		if err != nil {
			return err
		}
		defer img.Body.Close()
		imgfile, err := os.Create(imgName)
		_, err = io.Copy(imgfile, img.Body)

		defer imgfile.Close()

	}

	return err
}

func main() {
	startTime := time.Now()

	var (
		all      bool
		specific int
		xrange   int
		help     bool
		version  bool
		from     int
		to       int
	)

	flag.BoolVar(&all, "a", false, "Download all")
	flag.IntVar(&specific, "n", 0, "Download specific number")
	flag.IntVar(&xrange, "r", 0, "Download a range From To")
	flag.IntVar(&from, "f", 0, "Download all images starting From")
	flag.IntVar(&to, "t", 0, "Download all images up To")
	flag.BoolVar(&help, "h", false, "Display help")
	flag.BoolVar(&version, "v", false, "Display version")
	flag.Parse()

	if len(os.Args) == 1 {
		getComic(0)
	} else if help == true {
		fmt.Println("XKCD Downloader v0.1 ")
		fmt.Println("-------------------- ")
		flag.PrintDefaults()
		os.Exit(1)
	} else if version == true {
		fmt.Println("XKCD Downloader v0.1 ")
	} else if all == true {
		//downloadAll()
		from = 1
		to = getLatestCount()
		getComicRange(from, to)

	} else if specific != 0 {
		getComic(specific)
	}

	if from != 0 && to != 0 && from > to {
		from, to = to, from
	}

	if to != 0 && from == 0 {
		from = 1
	}

	if from != 0 && to == 0 {
		to = getLatestCount()
	}

	if from != 0 {
		log.Println("Starting from " + strconv.Itoa(from))
	}
	if to != 0 {
		log.Println("Ending up to " + strconv.Itoa(to))
	}

	getComicRange(from, to)

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("Total time taken: ", diff.Seconds(), " seconds")
}
