package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/zbioe/imdb/genrer"
	"github.com/zbioe/imdb/title"
)

const resultPATH = "./results"
const itemsPerReq = 50

var (
	SearchURL = "https://www.imdb.com/search/title"
	limit     = flag.Int("limit", 500, "limit per genrer")
	adult     = flag.Bool("adult", true, "incluse adult results")
	debug     = flag.Bool("debug", false, "verbose debug mode")
	sort      = flag.String("sort", "user_rating,desc", "sorted by")
)

func main() {
	flag.Parse()
	setup()

	resp, err := http.Get(SearchURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if _, err := os.Stat(resultPATH); os.IsNotExist(err) {
		err := os.Mkdir(resultPATH, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	var query = url.Values{}
	for _, s := range strings.SplitN(*sort, ",", -1) {
		query.Add("sort", s)
	}
	if *adult {
		query.Add("adult", "include")
	}

	var wg = new(sync.WaitGroup)
	var rawquery = query.Encode()
	for g := range genrer.Parse(resp.Body) {
		var rq = fmt.Sprintf("%s&genres=%s", rawquery, g)
		if *debug {
			log.Printf("start collect %s", g)
		}
		wg.Add(1)
		go collectTitles(wg, g, *debug, rq, *limit)
	}
	wg.Wait()
}

func collectTitles(
	wg *sync.WaitGroup,
	g genrer.Genrer,
	debug bool,
	rawquery string,
	limit int,
) {
	defer wg.Done()
	var sum int
	var npage = calculatePages(limit, itemsPerReq)
	var filepath = fmt.Sprintf("%s/%s.jsonl", resultPATH, g)

	f, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	for p := 0; p < npage; p++ {
		var rq = fmt.Sprintf("%s&start=%d", rawquery, (p*itemsPerReq)+1)
		var rawurl = fmt.Sprintf("%s?%s", SearchURL, rq)

		if debug {
			log.Printf("send request to %s", rawurl)
		}
		resp, err := http.Get(rawurl)
		if err != nil {
			log.Fatal(err)
		}

		if debug {
			log.Printf("start process titles of %s", rawurl)
		}
		result := title.Parse(resp.Body)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
		for t := range result.Titles {
			sum++
			if debug {
				log.Printf("%s: process %dº title of %dº page", g, sum, p)
			}
			if err := encoder.Encode(t); err != nil {
				log.Fatal(err)
			}
			if sum == limit {
				break
			}
		}

		if debug {
			log.Printf("finish process titles of %s", rawurl)
		}
		resp.Body.Close()
	}
	if debug {
		log.Printf("finish collect %s", g)
	}
}

func calculatePages(limit, itemsPerReq int) int {
	var npage int
	divided := float64(limit) / float64(itemsPerReq)
	if truncated := math.Trunc(divided); truncated == divided {
		npage = int(truncated)
	} else {
		npage = int(truncated) + 1
	}
	return npage
}

func setup() {
	log.SetFlags(0)
	log.SetPrefix("imdb: ")
}
