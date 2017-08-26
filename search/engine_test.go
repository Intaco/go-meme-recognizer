package search

import (
	"image_caching"
	"testing"
	"fmt"
)

const queriesChannelSize = 1024
const processedQueriesChannelSize = 1024

var stop int32

func TestStart(t *testing.T) {
	//установим соединение с базой данных
	image_caching.InitDb("database.db")
	defer image_caching.ClouseDb()

	queriesChan := make(chan Query, queriesChannelSize)
	processedQueriesChan := make(chan ProcessedQuery, processedQueriesChannelSize)

	//тестим - запихиваем в канал норм данные
	//queriesChan <- Query{1, "http://memepedia.ru/wp-content/uploads/2017/08/%D1%81%D0%BF%D0%B0%D0%B9%D0%B4%D0%B5%D1%80%D0%BC%D0%B5%D0%BC-%D1%88%D0%B0%D0%B1%D0%BB%D0%BE%D0%BD.jpeg", true}
	//queriesChan <- Query{2, "http://memepedia.ru/wp-content/uploads/2017/08/%D0%B1%D0%B5%D0%BD%D0%B4%D0%B5%D1%80-%D1%84%D1%83%D1%82%D1%83%D1%80%D0%B0%D0%BC%D0%B0.png", true}
	queriesChan <- Query{3, "кот качок", false}
	queriesChan <- Query{4, "С блэкджеком и шлюхами", false}
	// start search engine that:
	// 1. listens queriesChan
	// 2. processes queries
	// 3. put processedQueries to processedQueriesChan
	// 4. when queriesChannel is closed (and there're no more unprocessed queries) return
	go Start(queriesChan, processedQueriesChan)
	//print(<- processedQueriesChan, '\n')
	//print(<- processedQueriesChan, '\n')
	result := <- processedQueriesChan
	fmt.Printf("%t", result)
	fmt.Sprintf("изначальный запрос %s, заголовок %s, URL %s\n", result.Query.Query, result.Results[0].Title, result.Results[0].URL)
	result = <- processedQueriesChan
	fmt.Printf("%t", result)
	fmt.Sprintf("изначальный запрос %s, заголовок %s, URL %s\n", result.Query.Query, result.Results[0].Title, result.Results[0].URL)
	return
}
