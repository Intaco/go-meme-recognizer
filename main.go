package main

import (
	"log"
	"os"

	"./src/search"
	"github.com/abogovski/Go-TelegramBotAPI/tgbot"
	"./image_caching"
)

const queriesChannelSize = 1024
const processedQueriesChannelSize = 1024

var stop int32

func main() {
	//установим соединение с базой данных
	image_caching.InitDb("database.db")
	defer image_caching.ClouseDb()

	queriesChan := make(chan search.Query, queriesChannelSize)
	processedQueriesChan := make(chan search.ProcessedQuery, processedQueriesChannelSize)

	// start queries polling from telegram
	APIURL, err := tgbot.LoadBotAPIURL("tgbot.token")
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	var offset tgbot.Integer
	go PollTgBotQueries(APIURL, offset, queriesChan, &stop)

	// start telegram message dispatcher
	go DispatchTgBotResults(APIURL, processedQueriesChan)

	// start search engine that:
	// 1. listens queriesChan
	// 2. processes queries
	// 3. put processedQueries to processedQueriesChan
	// 4. when queriesChannel is closed (and there're no more unprocessed queries) return
	search.Start(queriesChan, processedQueriesChan)

	return
}