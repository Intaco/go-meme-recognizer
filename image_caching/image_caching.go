package image_caching

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"net/http"
)

var db *sql.DB      //глобальная переменная соединения с базой данных
var mu sync.RWMutex //sqlite3 блокируется при записи, и если попытаться записать во время блокировки -
// поток выполнения будет прерван

//вычисление хеша фотографии, неимпортируемая
func get_id(data []uint8) string {
	hash_sum := sha256.Sum256(data)    //найдём хеш сумму файла
	return fmt.Sprintf("%x", hash_sum) //вернём строку в 16 значном формате
}

//создание базы данных,
//принимает путь до файла sqlite
func InitDb(path string) error {
	var err error // если использовать := то создастся локальная переменная db
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//index
	sqlStmt := `
	create table if not exists list_of_memes(
		hash_id text primary key,
		meme_name text,
		link_to_mempedia text);
	create table if not exists name_and_id(
		meme_name text primary key,
		links_to_similar text
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}
	return err
}

//вставляет новую картинку в базу данных по хешу.
//на каждое добавление создаётся новая транзакция
//должно быть потокобезопасно, так как sqlite окружает rw mutex
//чертовски медленно на запись - 6 вызовов/секунда
func Insert_to_cash(data []uint8, name, link string) error {
	mu.Lock()
	defer mu.Unlock()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO list_of_memes (hash_id, meme_name, link_to_mempedia) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(get_id(data), name, link)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return nil
}

//скачиваем картинку и ищем по ней в кеше
func Select_from_cash(link string) (string, string, error) {
	response, _ := http.Get(link)
	defer response.Body.Close()
	var img_as_bytes []byte
	response.Body.Read(img_as_bytes)
	meme_name, link_to_mempedia, err := Select_from_cash_by_image(img_as_bytes)
	return meme_name, link_to_mempedia, err
}

//пытаемся просто вытащить строки
func Select_from_cash_by_image(data []byte) (string, string, error) {
	row := db.QueryRow("SELECT meme_name, link_to_mempedia FROM list_of_memes WHERE hash_id = $1",
		get_id(data))
	var meme_name, link_to_mempedia string
	err := row.Scan(&meme_name, &link_to_mempedia)
	if err != nil {
		return meme_name, link_to_mempedia, errors.New("no in the database")
	}
	return meme_name, link_to_mempedia, nil
}

func Insert_id(name, link string) error {
	mu.Lock()
	defer mu.Unlock()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO name_and_id (meme_name, links_to_similar) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, link)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func Select_id_by_name(name string)(string, error)  {
	row := db.QueryRow("SELECT links_to_similar FROM name_and_id WHERE meme_name = ?",
		name)
	var link_to_similar string
	err := row.Scan(&link_to_similar)
	if err != nil {
		return "", errors.New("no in the database")
	}
	return link_to_similar, nil
}

func ClouseDb() {
	db.Close()
}
