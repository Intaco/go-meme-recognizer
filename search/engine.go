package search

import "net/url"

// Start for now implements echo instead of search
func Start(queriesChan <-chan Query, processedQueriesChan chan<- ProcessedQuery) {
	for q := range queriesChan {
		entries := []Entry{{q.Query, "https://www.google.ru/search?q=" + url.QueryEscape(q.Query)}}
		processedQueriesChan <- ProcessedQuery{q, entries}
	}
	close(processedQueriesChan)
}
