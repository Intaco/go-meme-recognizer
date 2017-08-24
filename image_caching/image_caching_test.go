package image_caching

import "testing"

// TDD - какие запросы нужны

func TestAll(t *testing.T) {
	InitDb("caching.db")
	defer ClouseDb()

	//кеширование
	s := []byte{1, 2, 3}
	Insert_to_cash(s, "test string", "test link")
	name, link, err := Select_from_cash(s)
	print("test string - ", name, ", test link - ", link, ", nill - ", err, "\n")
	s = []byte{3, 2, 1}
	name, link, err = Select_from_cash(s)
	print("empty - ", name, ", empty - ", link, ", not nil - ", err, "\n")

	//сохранение имени мема и id
	Insert_id("имя мема", "id123456")
	id, err := Select_id_by_name("имя мема")
	print("id123456 - ", id, ", nil - ", err, "\n")
	id, err = Select_id_by_name("не имя мема")
	print("empty - ", id, ", not nill - ", err, "\n")
}
