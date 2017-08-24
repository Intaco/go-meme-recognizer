package request_processor

import "regexp"
import "get_mempedia_url"

//горутина, обрабатывающая приходящие запросы

func Run(/*два канала*/){

	//блокирующе ждём данные из канала
	go func(input_sting string, /*канал ответа*/) {
		if is_link(input_sting){
			//идём и загружаем
			//ищем в базе данных
		} else { //если нам передали строку с названием
			memepedia_url, err := get_mempedia_url.Get_mempedia_url(input_sting)
			//блокирующе ждём поиска на сайте и парсинг ответа
			if err != nil {
				return "не нашли такой мем" // в канал записать
			}else {
				return "Больше информации об этом меме вы можете узнать по ссылке\n" + memepedia_url
			}
		}
	}(input_sting)
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