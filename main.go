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
	return l.Num
}

var debug = false

func getComicRange(from int, to int) {
	fmt.Println("Getting ALL issues till date")
	fmt.Println("Total counts is ", to)

	// changed to the latest getting downloaded first to the oldest
	// so this works for update that only downloads the newer ones
	for i := to; i <= from; i-- {
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
	fname := "xkcd-" + n + ".txt"
	if debug == true {
		log.Println("[DEBUG] Img Path: " + imgName)
		log.Println("[DEBUG] filename: " + fname)
		//log.Println(xkcdImg.Alt)
	}

	if _, err := os.Stat(fname); err == nil {
		fmt.Println("[WARN] Txt description is already saved. ")
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
			log.Fatal("[ERROR] Filesize is 0")
		}
		f.Sync()
		defer f.Close()
	}

	if _, err := os.Stat(imgName); err == nil {
		log.Fatal("Image File " + imgName + " exists. Skipping...")

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
	flag.BoolVar(&debug, "d", false, "Verbose debug info")

	flag.Parse()

	if len(os.Args) == 1 {
		getComic(0)
		//		os.Exit(1)
	} else if help == true {
		fmt.Println("XKCD Downloader v0.1 ")
		fmt.Println("-------------------- ")
		flag.PrintDefaults()
		os.Exit(1)
	} else if version == true {
		fmt.Println("XKCD Downloader v0.1 ")
		os.Exit(1)
	} else if all == true {
		//downloadAll()
		from = 1
		to = getLatestCount()
		getComicRange(from, to)

	} else if specific != 0 {
		getComic(specific)
	} else {

		latest := getLatestCount()

		if from != 0 && to != 0 && from > to {
			from, to = to, from
		}

		if to != 0 && from == 0 {
			from = 1
		}

		if from != 0 && to == 0 {
			to = latest
		}

		if to > latest {
			to = latest
		}

		if from >= latest && to >= latest {
			getComic(0)
		} else {

			if from != 0 {
				log.Println("Starting from " + strconv.Itoa(from))
			}
			if to != 0 {
				log.Println("Ending up to " + strconv.Itoa(to))
			}

			getComicRange(from, to)
		}

	}

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("Total time taken: ", diff.Seconds(), " seconds")
}
