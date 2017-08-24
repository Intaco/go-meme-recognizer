package main

import (
	"./googlecse"
)

func main() {

}
func cseSample(){
	url := "http://runt-of-the-web.com/wordpress/wp-content/uploads/2013/11/shibe-meme-wont-go.jpg"
	googlecse.SearchSimilarByUrl(url)
}