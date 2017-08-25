package image_caching

import "testing"

// TDD - какие запросы нужны

func TestAll(t *testing.T) {
	InitDb("caching.db")
	defer ClouseDb()

	//кеширование
	s := []byte{1, 2, 3}
	Insert_to_cash_by_image(s, "test string", "test link")
	name, link, err := Select_from_cash_by_image(s)
	print("test string - ", name, ", test link - ", link, ", nill - ", err, "\n")
	s = []byte{3, 2, 1}
	name, link, err = Select_from_cash_by_image(s)
	print("empty - ", name, ", empty - ", link, ", not nil - ", err, "\n")

	//сохранение имени мема и id
	Insert_id("имя мема", "id123456")
	id, err := Select_id_by_name("имя мема")
	print("id123456 - ", id, ", nil - ", err, "\n")
	id, err = Select_id_by_name("не имя мема")
	print("empty - ", id, ", not nill - ", err, "\n")

	//кеширование по url
	Insert_to_cash("http://memepedia.ru/wp-content/uploads/2017/08/%D0%BF%D0%BE%D1%82%D1%80%D0%B0%D1%87%D0%B5%D0%BD%D0%BE-8.jpg", "Потрачено", "https://memepedia.ru/potracheno/")
	name, link, err = Select_from_cash("http://memepedia.ru/wp-content/uploads/2017/08/%D0%BF%D0%BE%D1%82%D1%80%D0%B0%D1%87%D0%B5%D0%BD%D0%BE-8.jpg")
	print(name, " ",link, " ", err, "\n")

}
