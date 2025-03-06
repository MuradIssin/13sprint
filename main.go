package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MuradIssin/go_final_project/handlers"
	"github.com/MuradIssin/go_final_project/operateDB"
	"github.com/MuradIssin/go_final_project/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// для запуска и проверки в командой строке	// $env:TODO_PORT="8080"; go run main.go

	//тестирование NextDate
	now := time.Now()
	now = time.Date(2024, time.January, 26, 0, 0, 0, 0, time.UTC)
	resfun, err := utils.NextDate(now, "20240229", "y")
	fmt.Println(resfun)
	resfun, err = utils.NextDate(now, "20240113", "d 7")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resfun)

	const defaultPort = "7540"
	const webDir = "./web"

	// Получаем значение переменной окружения TODO_PORT
	portIncome, exists := os.LookupEnv("TODO_PORT")
	if !exists || portIncome == "" {
		log.Println("Port:", portIncome)
	}
	fileNameDB, exists := os.LookupEnv("TODO_DBFILE")
	if !exists || fileNameDB == "" {
		log.Println("Database File:", fileNameDB)
	}

	// запрашиваем подключение к БД - проверяем наличия денег
	operateDB.СheckDb(fileNameDB)

	// Подключаемся к БД
	db, err := operateDB.InitDB(fileNameDB)
	defer db.Close()
	if err != nil {
		fmt.Println("Ошибка инициализации БД ", err)
	}

	// Если переменная окружения пустая, устанавливаем порт по умолчанию
	if portIncome == "" {
		portIncome = defaultPort
		log.Println("сработало замена полученного порта")
	}
	fmt.Println("Используем порт:", portIncome)

	// // Настроим обработчик для статических файлов из директории webDir
	// http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Создаём роутер
	r := chi.NewRouter()

	// Запускаем Web интерфейс
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Post("/api/task", handlers.AddTask)

	// Запуск сервера на указанном порту
	if err := http.ListenAndServe(":"+portIncome, r); err != nil {
		log.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
