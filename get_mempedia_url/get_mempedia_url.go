package  get_mempedia_url

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"errors"
)

/*func main() {
	link, _:= Get_mempedia_url("нельзя просто так взять и")
	print(link, "\n") //https://memepedia.ru/boromir/
	link, _ = Get_mempedia_url("давай досвидания")
	print(link) //https://memepedia.ru/davaj-do-svidaniya/
}*/


func Get_mempedia_url(search_word string) (string, error) {
	req, err := http.NewRequest("GET", "http://memepedia.ru/search/"+search_word, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linu…) Gecko/20100101 Firefox/55.0")
	req.Header.Set("Accept-Language", "en-US;q=0.5")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	r, _ := regexp.Compile("<div id=\"sf-module-search\" class=\"sf-module sf-module-grid-posts sf-entries sf-entries-grid_posts sf-entries-col-2 sf-entries-vertical\">.*?<h3 class=\"sf-entry-title\"><a href=\"(https://memepedia.ru/.*?/)\" rel=\"bookmark\">.*?</a></h3>")
	link := r.FindStringSubmatch(string(body))[1]
	if link == "" {
		return "", errors.New("нет такого мема")
	}
	return link, nil
}
