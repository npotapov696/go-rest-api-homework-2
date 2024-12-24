package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task описывает задачу.
type Task struct {
	ID           string   `json:"id"`           // ID задачи
	Description  string   `json:"description"`  // Заголовок
	Note         string   `json:"note"`         // Описание задачи
	Applications []string `json:"applications"` // Используемые приложения
}

// task содержит стартовый набор задач типа Task.
var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks возвращает все задачи, содержащиеся в мапе tasks.
// В случае ошибки возвращает статус 500 Iternal Server Error.
func getTasks(w http.ResponseWriter, r *http.Request) {

	// сериализуем данные из мапы tasks,
	// обрабатываем ошибку сериализации
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок ответа записываем тип контента - данные в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, записываем статус OK
	w.WriteHeader(http.StatusOK)

	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// postTask добавляет новую задачу в мапу tasks, или заменяет существующую,
// если ID новой задачи совпадает с ID одной из задач в мапе.
// В случае ошибки возвращает статус 400 Bad Request.
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	// записываем данные из тела запроса в буфер,
	// обрабатываем ошибку чтения данных
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// десериализуем данные из буфера и записываем их в перемнную task,
	// обрабатываем ошибку десериализации
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// добавляем в мапу tasks новый элемент из переменной task
	tasks[task.ID] = task

	// в заголовок ответа записываем тип контента - данные в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, записываем статус успешного создания элемента
	w.WriteHeader(http.StatusCreated)
}

// getTask возвращает задачу по указанному ID в параметрах URL запроса из мапы tasks.
// В случае ошибки или если задача не найдена, возвращает статус 400 Bad Request.
func getTask(w http.ResponseWriter, r *http.Request) {

	// записываем в переменную параметр id из URL запроса
	id := chi.URLParam(r, "id")

	// проверяем наличие запрашиваемого id в мапе tasks, если нет, возвращаем соотвующее сообщение,
	// если есть, записываем задачу в переменную task
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача с указанным id отсутствует", http.StatusBadRequest)
		return
	}

	// сериализуем данные из мапы tasks,
	// обрабатываем ошибку сериализации
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// в заголовок записываем тип контента - данные в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус 200 OK
	w.WriteHeader(http.StatusOK)

	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// deleteTask удаляет задачу по указанному ID в параметрах URL запроса из мапы tasks.
// В случае ошибки или если задача не найдена, возвращает статус 400 Bad Request.
func deleteTask(w http.ResponseWriter, r *http.Request) {

	// записываем в переменную парамет id из URL запроса
	id := chi.URLParam(r, "id")

	// проверяем наличие запрашиваемого id в мапе tasks, если нет, возвращаем соотвующее сообщение и статус
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача с указанным id отсутствует", http.StatusBadRequest)
		return
	}

	// удаляем задачу из мапы tasks с ключом id
	delete(tasks, id)

	// в заголовок записываем тип контента - данные в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// так как все успешно, то статус 200 OK
	w.WriteHeader(http.StatusOK)
}

func main() {

	// создаем роутер
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)

	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTask`
	r.Post("/tasks", postTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)

	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE, для которого используется обработчик `deleteTask`
	r.Delete("/tasks/{id}", deleteTask)

	// запускаем сервер, обрабатываем возможную ошибку
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
