package  get_mempedia_url

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"errors"
)


func Get_mempedia_url(search_word string) (string, string, error) {
	req, err := http.NewRequest("GET", "http://memepedia.ru/search/"+search_word, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linu…) Gecko/20100101 Firefox/55.0")
	req.Header.Set("Accept-Language", "en-US;q=0.5")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile("<div id=\"sf-module-search\" class=\"sf-module sf-module-grid-posts sf-entries sf-entries-grid_posts sf-entries-col-2 sf-entries-vertical\">.*?<h3 class=\"sf-entry-title\"><a href=\"(https://memepedia.ru/.*?/)\" rel=\"bookmark\">(.*?)</a></h3>")
	var result []string = r.FindStringSubmatch(string(body))
	if len(result) != 3{
		return "", "", errors.New("нет такого мема")
	}
	var link string = result[1]
	var name string = result[2]
	return link, name, nil
}
