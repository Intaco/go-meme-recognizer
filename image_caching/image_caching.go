package image_caching

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
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
	sqlStmt := `
	create table if not exists list_of_memes(
		id integer not null primary key autoincrement,
		hash_id text,
		meme_name text,
		link_to_mempedia text
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
func Insert_to_database(data []uint8, name, link string) error {
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

//пытаемся просто вытащить строки
func Select_from_database(data []byte) (string, string, error) {
	row := db.QueryRow("SELECT meme_name, link_to_mempedia FROM list_of_memes WHERE hash_id = $1",
		get_id(data))
	var meme_name, link_to_mempedia string
	err := row.Scan(&meme_name, &link_to_mempedia)
	if err != nil {
		return meme_name, link_to_mempedia, errors.New("no in the database")
	}
	return meme_name, link_to_mempedia, nil
}

func ClouseDb() {
	db.Close()
}
