package handlers

import (
	"encoding/json"
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

type Response struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// type Storage struct {
// 	db *sql.DB
// }

func AddTask(w http.ResponseWriter, r *http.Request, db *database.Database) {
	// func AddTask(w http.ResponseWriter, r *http.Request) {
	// (d *Database)
	// func AddTask(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// if r.Method != http.MethodPost {
	// 	SendErrorResponse(w, "AddTaskHandler: Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	log.Println("new task")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Десериализация JSON в структуру Task
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		// http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		// return
		response := Response{Error: "Ошибка десериализации JSON"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)

	}
	log.Println("1: десериализации JSON успешно ")

	// Проверка обязательного поля title
	if task.Title == "" {
		// http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		response := Response{Error: "Не указан заголовок задачи"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else {
		log.Println("2: обязательное поле заголовок задачи указан ")
	}

	// если дата не проставлена установим текущую дату
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
		log.Println("проставляется текущая дата", task.Date)
	} else {
		log.Println("3: дата не пустая")
	}

	// Проверка правильности формата даты
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		log.Println("Неверный формат даты")
		// http.Error(w, `{"error": "Неверный формат даты"}`, http.StatusBadRequest)
		response := Response{Error: "Неверный формат даты"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Println("4: правильность формата данных установлена")

	// 	Если дата меньше сегодняшнего числа, есть два варианта:
	// - если правило повторения не указано или равно пустой строке, подставляется сегодняшнее число;
	if task.Date < time.Now().Format("20060102") && task.Repeat == "" {
		task.Date = time.Now().Format("20060102")
		log.Println("Проставляется текущая дата для задачи", task.Date)
	} else {
		log.Println("5: нет повтора и дата выше актуальной ")
	}
	// - при указанном правиле повторения вам нужно вычислить
	// и записать в таблицу дату выполнения, которая будет больше сегодняшнего числа.
	// Для этого используйте функцию NextDate(), которую вы уже написали раньше.

	if task.Repeat != "" {
		now := time.Now()
		taskDate, _ := time.Parse("20060102", task.Date)
		newDate, err := utils.NextDate(now, taskDate.Format("20060102"), task.Repeat)
		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			// json.NewEncoder(w).Encode("response ошибка ")
			// http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
			// return
			response := Response{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			log.Println("вернули ошибку с формированием даты ")
			return
		}
		task.Date = newDate
		log.Println("установлена новая дата для задачи", newDate)
	}

	log.Println("6: записать в БД ", task.Date, task.Title, task.Comment, task.Repeat)

	// Добавляем задачу в базу данных (здесь пример SQL-запроса)
	// db.Exec - здесь предполагается, что у вас есть подключение к базе данных (db)
	// Пример добавления задачи

	// var err error
	// var result []byte
	// var returnData IDType

	// Формируем запрос в базу

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		http.Error(w, `{"error": "Ошибка добавления задачи в базу данных"}`, http.StatusInternalServerError)
		return
	}

	// Получаем ID добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, `{"error": "Ошибка получения ID задачи"}`, http.StatusInternalServerError)
		return
	}

	log.Panicln(id)

	// // Возвращаем успешный результат с ID
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(map[string]interface{}{"id": id})

}
