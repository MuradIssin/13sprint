package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MuradIssin/go_final_project/utils"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(w http.ResponseWriter, r *http.Request) {

	log.Println("new task")

	// Десериализация JSON в структуру Task
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	} else {
		log.Println("десериализации JSON успешно ")
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	} else {
		log.Println("обязательное поле заголовок задачи указан ")
	}

	// Проверка формата даты и если дата задачи позже текущей даты и нет повтора
	if task.Date == "" || (task.Date < time.Now().Format("20060102") && task.Repeat == "") {
		// Если дата не указана, берем сегодняшнюю дату
		task.Date = time.Now().Format("20060102")
		log.Println("дата не указана, проставляется текущая дата")
	} else {
		// Проверка правильности формата даты
		_, err := time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат даты"}`, http.StatusBadRequest)
			return
		}
		log.Println("правильность формата данных установлена")
	}

	// 	Если дата меньше сегодняшнего числа, есть два варианта:
	// если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число;
	// при указанном правиле повторения вам нужно вычислить и записать в таблицу дату выполнения, которая будет больше сегодняшнего числа. Для этого используйте функцию NextDate(), которую вы уже написали раньше.

	now := time.Now()

	taskDate, _ := time.Parse("20060102", task.Date)

	newDate, err := utils.NextDate(now, taskDate.Format("20060102"), task.Repeat)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}
	log.Println("установлена новая дата", newDate)

	// 		task.Date = nextDate

	// Добавляем задачу в базу данных (здесь пример SQL-запроса)
	// db.Exec - здесь предполагается, что у вас есть подключение к базе данных (db)
	// Пример добавления задачи

	// query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	// res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	// if err != nil {
	// 	http.Error(w, `{"error": "Ошибка добавления задачи в базу данных"}`, http.StatusInternalServerError)
	// 	return
	// }

	// // Получаем ID добавленной записи
	// id, err := res.LastInsertId()
	// if err != nil {
	// 	http.Error(w, `{"error": "Ошибка получения ID задачи"}`, http.StatusInternalServerError)
	// 	return
	// }

	// // Возвращаем успешный результат с ID
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]interface{}{"id": id})

}
