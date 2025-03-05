package operateDB

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// type DB struct {
// 	db *sql.DB
// }

// go get modernc.org/sqlite - требуется запустить  для начала рабты с бд
func СheckDb(fileName string) {

	// $env:CGO_ENABLED="1"

	if fileName == "" {
		fileName = "scheduler.db"
	}

	// Получаем путь к исполняемому файлу приложения
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("appPath:", appPath)

	// Формируем путь к файлу базы данных
	dbFile := filepath.Join(filepath.Dir(appPath), fileName)
	log.Println("dbFile:", dbFile)

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(fileName)
	var install bool
	if err != nil {
		log.Println("os.Stat(dbFile)", err)
		install = true
		log.Println("требуется создания файла БД ", fileName, install)
	}

	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		err = createTables(fileName)
		if err != nil {
			log.Fatal("Ошибка при создании таблиц:", err)
		}
		log.Println("база данных создана", fileName)
	} else {
		log.Println("база данных существует", fileName)
	}

	// можно убрать т.к. инициализация идет позже
	// db, err := sql.Open("sqlite", dbFile)
	// if err != nil {
	// 	log.Fatal("Ошибка при открытии базы данных:", err)
	// }
	// defer db.Close()
}

func createTables(fileName string) error {
	query := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		date TEXT, 
		title TEXT,
		comment TEXT,
		repeat TEXT CHECK(length(repeat) <= 128)
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		log.Fatal("Ошибка при открытии базы данных:", err)
	}
	defer db.Close()

	_, err = db.Exec(query)
	return err
}

func InitDB(fileName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}
	return db, nil
}
