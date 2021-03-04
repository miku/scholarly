package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sethgrid/pester"
)

var (
	start    = flag.Int("s", 0, "start")
	end      = flag.Int("e", 1000, "end")
	user     = flag.String("u", "", "username")
	password = flag.String("p", "", "password")
	sleep    = flag.Duration("d", 1*time.Second, "delay")
)

type UserResponse struct {
	Stat  string `json:"stat"`
	Users struct {
		End   int64  `json:"end"`
		Next  string `json:"next"`
		Start int64  `json:"start"`
		User  []struct {
			Groups struct {
				End   int64 `json:"end"`
				Start int64 `json:"start"`
			} `json:"groups"`
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"user"`
	} `json:"users"`
}

func main() {
	flag.Parse()
	client := pester.New()
	for {
		// curl -XGET --user username:apikey "https://www.bibsonomy.org/api/users?end=2&format=json"
		link := fmt.Sprintf("https://www.bibsonomy.org/api/users?start=%d&end=%d&format=json", *start, *end)
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.SetBasicAuth(*user, *password)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%d %s", resp.StatusCode, link)
		if resp.StatusCode == 404 {
			break
		}
		if resp.StatusCode >= 400 {
			log.Fatalf("got http %d", resp.StatusCode)
		}
		tee := io.TeeReader(resp.Body, os.Stdout)
		dec := json.NewDecoder(tee)
		var uresp UserResponse
		if err := dec.Decode(&uresp); err != nil {
			log.Fatal(err)
		}
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		if uresp.Users.Next == "" {
			break
		} else {
			time.Sleep(*sleep)
		}
		*start = *start + 1000
		*end = *end + 1000
	}

}
