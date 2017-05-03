package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type InterestingResp struct {
	Photos Photos `json:"photos"`
}

type Photos struct {
	Page    int     `json:"page"`
	Pages   int     `json:"pages"`
	Perpage int     `json:"perpage"`
	Total   int     `json:"total"`
	Photo   []Photo `json:"photo"`
}

type Photo struct {
	ID     string `json:"id"`
	Owner  string `json:"owner"`
	Secret string `json:"secret"`
	Server string `json:"server"`
	Farm   int    `json:"farm"`
}

func (p *Photo) URL() string {
	fmt.Printf("https://farm%v.staticflickr.com/%v/%v_%v_q.jpg\n", p.Farm, p.Server, p.ID, p.Secret)
	return fmt.Sprintf("https://farm%v.staticflickr.com/%v/%v_%v_q.jpg", p.Farm, p.Server, p.ID, p.Secret)
}

func (p *Photo) Fetch() ([]byte, error) {
	resp, err := http.Get(p.URL())
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return b, nil
}

// date in "YYYY-MM-DD" format
func fetchInteresting(date, apiKey string) ([]Photo, error) {
	base := "https://api.flickr.com/services/rest/?method=flickr.interestingness.getList"
	url := fmt.Sprintf("%v&api_key=%v&date=%v&per_page=500&format=json&nojsoncallback=1", base, apiKey, date)

	resp, err := http.Get(url)
	if err != nil {
		return []Photo{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Photo{}, err
	}
	defer resp.Body.Close()

	var respJSON InterestingResp
	err = json.Unmarshal(data, &respJSON)
	if err != nil {
		return []Photo{}, err
	}

	return respJSON.Photos.Photo, nil
}

func store(p Photo, url string) error {
	b, err := p.Fetch()
	if err != nil {
		return err
	}

	_, err = http.Post(url, "image/jpeg", bytes.NewReader(b))
	return err
}

func datesBetween(start, stop string) ([]string, error) {
	const shortForm = "2006-01-02"
	dates := []string{}

	begin, err := time.Parse(shortForm, start)
	if err != nil {
		return dates, err
	}

	end, err := time.Parse(shortForm, stop)
	if err != nil {
		return dates, err
	}

	for d := begin; d != end; d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format(shortForm))
	}

	return dates, nil
}

func main() {
	apiKey := os.Getenv("FLICKR_API_KEY")
	if apiKey == "" {
		fmt.Println("FLICKR_API_KEY not set")
		return
	}

	//dates, err := datesBetween("2010-01-01", "2017-01-01")
	dates, err := datesBetween("2017-01-01", "2017-02-01")
	if err != nil {
		fmt.Println(err)
	}

	for _, d := range dates {
		fmt.Println("Fetching interesting photos for", d)
		photos, err := fetchInteresting(d, apiKey)
		if err != nil {
			fmt.Println(nil)
		}

		for _, p := range photos {
			err := store(p, "http://localhost:8080/store")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
