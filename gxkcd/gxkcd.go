package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type result byte

type xkcd struct {
	Num        int    `json:"num"`
	Title      string `json:"title"`
	Safe_title string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

const (
	OK result = iota
	ERR_NEXT
	ERR_RETRY
)

var logger = log.Default()

func _download_one(url string, fname string) result {
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("Error %s", err)
		return ERR_RETRY
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		outfile, err := os.Create(fname)
		if err != nil {
			logger.Printf("Error creating file: %s", err)
			return ERR_RETRY
		}
		defer outfile.Close()
		_, err = io.Copy(outfile, resp.Body)
		if err != nil {
			logger.Printf("Error copying body: %s", err)
			return ERR_RETRY
		}
		logger.Println("OK")
		return OK
	case 404:
		logger.Println("Not found, skipping")
		return ERR_NEXT
	default:
		logger.Printf("Unexpected code %s", resp.Status)
		return ERR_RETRY
	}
}

func download(dir string, maxfail int) {
	logger.Printf("Download: dir %s, maxfail %d", dir, maxfail)
	fails := 0
	for i := 1; ; {
		logger.Printf("Trying XKCD #%d", i)
		switch res := _download_one(
			fmt.Sprintf("https://xkcd.com/%d/info.0.json", i),
			fmt.Sprintf("%s/%04d.json", dir, i),
		); res {
		case OK:
			fails = 0
			i++
		case ERR_RETRY:
			fails++
			logger.Printf("Error %d of %d, retry", fails, maxfail)
		case ERR_NEXT:
			fails++
			i++
			logger.Printf("Error %d of %d, skip", fails, maxfail)
		}
		if fails >= maxfail {
			logger.Println("Max error reached, stopping")
			return
		}

	}
}

func allterms(s string, terms []string) bool {
	for _, t := range terms {
		if !strings.Contains(s, strings.ToLower(t)) {
			return false
		}
	}
	return true
}

func search(dir string, terms []string) {
	if len(terms) == 0 {
		return
	}
	fn := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		blob, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		x := xkcd{}
		err = json.Unmarshal(blob, &x)
		ss := strings.ToLower(strings.Join([]string{x.Title, x.Safe_title, x.Transcript, x.Alt}, "\n"))
		if allterms(ss, terms) {
			logger.Printf("Found match: %d", x.Num)
		}
		return nil
	}

	err := filepath.WalkDir(dir, fn)
	if err != nil {
		logger.Printf("Error %s", err)
	}
}

func main() {
	var dir string
	var maxfail int

	flag.StringVar(&dir, "dir", ".", "Data directory")
	flag.IntVar(&maxfail, "maxfail", 5, "Max number of errors in sequence")
	flag.Parse()

	switch a := flag.Arg(0); a {
	case "":
		logger.Fatal("Need action: download or search")
	case "download":
		logger.Println("Downloading XKCD")
		download(dir, maxfail)
	case "search":
		logger.Println("Search")
		search(dir, flag.Args()[1:])
	default:
		logger.Fatalf("Unknown action %#v", a)
	}

}
