package search

import (
	"regexp" // - для проверки url ли прислали. Это проверяется тут?
	"get_mempedia_url" //перемещено в src - иначе не видит
	"image_caching"
)
//горутина, обрабатывающая приходящие запросы
// Start for now implements echo instead of search

func Start(queriesChan <-chan Query, processedQueriesChan chan<- ProcessedQuery) {
	for q := range queriesChan { //берём запрос(блокирующе) пока канал не закрыт
		//зададим литеральную функцию - отдельную горутину обработчик
		//передаём запрос
		go func(q Query) {
			if q.IsURL { //если ссылка - скачиваем
				meme_name, link_to_mempedia, nil := image_caching.Select_from_cash(q.Query)
				//тут проверить ошибку и вызвать поисковик в случае неудачи
				
				//как найдётся - поискать в мемпедии по заголовкам ссылку на мемпедию
				//вернуть в канал ответ
				//и вставить в таблицу имя мема(встать в очередь на добавление)
			}else{ // если не ссылка - ищем в мемпедии по названию
				mempedia_url, err := get_mempedia_url.Get_mempedia_url(q.Query)
				if err != nil {
					entries := []Entry{{"Не знаю такого мема ☹. Перефразируйте.️", ""}} //QueryEscape escapes the string so it can be safely placed inside a URL query.
					processedQueriesChan <- ProcessedQuery{q, entries}
				}else {
					entries := []Entry{{"Вот что мы нашли по этому мему", mempedia_url}} //QueryEscape escapes the string so it can be safely placed inside a URL query.
					processedQueriesChan <- ProcessedQuery{q, entries}
					//просто возвращаем пользователю ссылочку
				}
			}

		}(q)//строка запроса должна быть неразделяемой, а канал возврата общим для всех обработчиков
	}
	close(processedQueriesChan)
}

func is_link(may_be_url string) bool {
	checking, _ := regexp.Compile(`https?://.*\.(jpe?g|gif|png)`)
	link := checking.FindString(may_be_url)
	if link == ""{
		return false
	} else {
		return true
	}
}