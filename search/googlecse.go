package search

import (
	"encoding/json"
	"net/url"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
)

var config CFG = loadCfg()

func GetConfig() CFG {
	return config
}
type configSource struct {
	CFG CFG `json:"GoogleCSE"`
}
type CFG struct {
	Cx  string
	Key string
}

func loadCfg() CFG {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config.json! Error: %s", err)
	}
	configuration := new(configSource)

	err = json.Unmarshal(bytes,configuration)

	if err != nil {
		log.Fatalf("Failed to unmarshall json! Error: %s", err)
	}

	return configuration.CFG
}


type Answers struct {
	Memes [] CSAMemeDesc `json:"items"`
}
type CSAMemeDesc struct {
	Title string
	Link string
}
const GOOGLE_URL = "https://www.googleapis.com/customsearch/v1"

func GetSimilarMemes(q Query, processedQueriesChan chan<- ProcessedQuery) {
	fmt.Println("=======================\n in memes \n")
	if q.IsURL {
		memes, err := searchSimilarByUrl(q.Query)
		if err != nil {
			entries := []Entry{{"Упс, кажется, я не смог найти похожие мемы!", ""}}
			processedQueriesChan <- ProcessedQuery{q, entries}
		} else {
			size := len(memes) + 1 // +1 for comment
			entries := make([]Entry, size)
			entries[0] = Entry{Title: "Вот какие похожие мемы я смог найти:"}
			i := 1
			for _, s := range memes {
				entries[i].Title = s.Title
				entries[i].URL = s.Link
				i++
			}
			processedQueriesChan <- ProcessedQuery{q, entries}
		}
	}
}

func searchSimilarByUrl(imageURL string) ([]CSAMemeDesc, error) {
	iu, err := url.Parse(imageURL)
	if err != nil {
		log.Fatalf("Failed to parse image url! Error: %s", err)
		return nil, err
	}
	cfg := GetConfig()
	u,err  := url.Parse(GOOGLE_URL)
	if err != nil {
		log.Fatalf("Failed to parse CSE base url! Error: %s", err)
		return nil, err
	}
	q := u.Query()
	q.Add("key", cfg.Key)
	q.Add("cx", cfg.Cx)
	q.Add("searchType", "image")
	q.Add("alt", "json")
	q.Add("q", iu.String())
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("Failed to proceed request to google CSE! Error: %s code", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read CSE answer! Error: %s", err)
		return nil, err
	}
	return parseResponseBody(body)
}
func parseResponseBody(body []byte) ([]CSAMemeDesc, error) {
	println(string(body))
	var answers = Answers{}
	err := json.Unmarshal(body, &answers)
	if err != nil {
		log.Fatalf("Failed to unmarshal response body! Error: %s", err)
		return nil, err
	}
	return answers.Memes, nil
}
