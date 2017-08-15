package image_caching

/*пакет кеширования
интерфейс:
функция получает картинку
в отдельной go подпрограмме
	находит хеш файла - понятно
	ищет в базе данных этот хеш
	если находит - возвращает описание картинки
	иначе обращается к другому модулю, который обращается к Google распознаванию картинок
	возвращает
 */

import (
	"crypto/sha256"
	"unsafe" //ради самопальной реализации union
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//Надеюсь, скоро. https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
	"fmt"
)

//глобальная переменная -- ссылка на базу данных? (https://habrahabr.ru/post/332122/)
/*при инициализации и открытии файла надо выполнить создание базы
CREATE TABLE IF NOT EXISTS list_of_memes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hash_id INTEGER,
    meme_name TEXT,
    link_to_mempedia TEXT
);
 */


func get_image_name(data []uint8) (string, string) { //передаётся указатель на область памяти без копирования массива
	hash_sum := sha256.Sum256([]byte("hello world\n")) //найдём хеш сумму файла
	//возьмём только первые 8 байт из 32-х, чтоб влезло в uint64, который совместим с числом в базе данных
	hash_id := *(*uint64)(unsafe.Pointer(&hash_sum))
	fmt.Printf("hash sum = %x, hash id = %x\n", hash_sum, hash_id)
	//запрашиваем из базы данных имя мема и ссылку на мемопедию
	name, link, err := get_from_database(hash_id)
	if err != nil {
		name, _ := get_from_google_api() //не этот модуль
		link, _ := get_from_memopedia() //не этот модуль
		go insert_to_database(name, link)
	}
	return name, link
}

func get_from_database(hash_id uint64) (string, string, error){
	var name, link string
	var err error
	/*SELECT meme_name, link_to_mempedia FROM list_of_memes
    WHERE hash_id = ?;*/
	return name, link, err
}

func insert_to_database(name, link string){
	/*INSERT INTO list_of_memes
	(hash_id, meme_name, link_to_mempedia)
	VALUES(?, ?, ?);*/
}