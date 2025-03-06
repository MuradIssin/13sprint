package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// https://github.com/Yandex-Practicum/go_final_project/pull/97/files#diff-3766b045965a17661430b847219d1b624e9a517256eaa1e9e9df5b467022045f

func NextDate(now time.Time, date string, repeat string) (string, error) {

	// now — время от которого ищется ближайшая дата;
	// date — исходное время в формате 20060102, от которого начинается отсчёт повторений;
	// repeat — правило повторения в описанном выше формате.

	// 	Если правило не указано, отмеченная выполненной задача будет удаляться из таблицы;
	// d <число> — задача переносится на указанное число дней. Максимально допустимое число равно 400. Примеры:
	// d 1 — каждый день;
	// d 7 — для вычисления следующей даты добавляем семь дней;
	// d 60 — переносим на 60 дней.
	// y — задача выполняется ежегодно. Этот параметр не требует дополнительных уточнений. При выполнении задачи дата перенесётся на год вперёд.

	// NextDate(now, "20240229", "y") возвращает 20250301;
	// NextDate(now, "20240113", "d 7") возвращает 20240127;

	// Парсим исходную дату из формата "20060102" в объект time.Time
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("неверный формат исходной даты: %w", err)
	}

	// Убираем пробелы и приводим строку repeat к нижнему регистру
	repeat = strings.ToLower(strings.TrimSpace(repeat))

	// если в колонке repeat — пустая строка;
	// if repeat == "" {
	// 	return "", fmt.Errorf("правило повторений не содержит данных")
	// }

	// Переменная для хранения следующей даты
	var nextDate time.Time

	// Если правило повторения "d <число>"
	if strings.HasPrefix(repeat, "d") {
		// Извлекаем число после "d"
		var days int
		_, err := fmt.Sscanf(repeat, "d %d", &days)
		if err != nil {
			return "", fmt.Errorf("неверный формат для 'd': %w", err)
		}
		// Проверка на максимально допустимое количество дней (до 400)
		if days < 1 || days > 400 {
			return "", errors.New("число дней должно быть в пределах от 1 до 400")
		}
		// Если правило "d", добавляем дни к стартовой дате
		nextDate = startDate.AddDate(0, 0, days)
	}

	// Если правило повторения "y" (ежегодно)
	if repeat == "y" {
		// Если правило "y", переносим дату на 1 год вперед
		nextDate = startDate.AddDate(1, 0, 0)
	}

	// Если не поддерживаемое правило
	if nextDate.IsZero() {
		return "", errors.New("неизвестное правило повторения")
	}

	// Если полученная дата меньше или равна текущей, повторяем шаг
	for !nextDate.After(now) {
		// Для "d", добавляем дни ещё раз
		if strings.HasPrefix(repeat, "d") {
			nextDate = nextDate.AddDate(0, 0, 1) // Переносим на день дальше
		}
		// Для "y", переносим на 1 год вперед
		if repeat == "y" {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
	}

	// Возвращаем дату в формате "20060102"
	return nextDate.Format("20060102"), nil

}
