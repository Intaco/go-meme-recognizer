package main

import (
	"fmt"
	"log"
	"sync"

	"./search"
	"github.com/abogovski/Go-TelegramBotAPI/tgbot"
)

const timeout = tgbot.Integer(15) // return from tgbot.getUpdates every 15 secs
const limit = tgbot.Integer(100)  // default max updates returned per getUpdates call

// PollTgBotQueries poll telegram bot for queries
func PollTgBotQueries(botAPIURL string, offset tgbot.Integer, output chan<- search.Query, stop *int32) (tgbot.Integer, int, error) {
	log.Println("Start Polling")
	params := tgbot.Params{
		"offset":          offset,
		"timeout":         timeout,
		"limit":           limit,
		"allowed_updates": []string{"message"}}

	input := make(chan tgbot.Update)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(output chan<- search.Query, input <-chan tgbot.Update) {
		for update := range input {
			if update.Message != nil && update.Message.Text != nil {
				isURL := len(update.Message.Entities) == 1 && update.Message.Entities[0].URL != nil
				output <- search.Query{update.Message.Chat.ID, *update.Message.Text, isURL}
			} else {
				log.Printf("Skipping update that is not a text message: %v\n", update)
			}
		}
		fmt.Println("End polling")
		close(output)
		wg.Done()
		return
	}(output, input)

	offset, status, err := tgbot.PollUpdates(botAPIURL, params, input, stop)
	wg.Wait()
	return offset, status, err
}

// DispatchTgBotResults distpacth processed queries (telegram bot)
func DispatchTgBotResults(botAPIURL string, input <-chan search.ProcessedQuery) {
	for pq := range input {
		_, status, err := tgbot.SendMessage(botAPIURL, tgbot.Params{
			"chat_id":                  pq.Query.ClientID,
			"text":                     pq.ToMarkdown(),
			"parse_mode":               "Markdown",
			"disable_web_page_preview": true})
		if status == 429 {
			log.Println("DispatchTgBotResults: SendMessage RateLimitExceeded")
		} else if err != nil {
			panic(err)
		}
	}
}
