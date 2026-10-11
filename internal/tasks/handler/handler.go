package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"taskTitle/internal/tasks/service"
)

type Handler struct {
	s *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{s: s}
}

func (s *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		http.Error(w, "Пустое тело запроса", http.StatusBadGateway)
		return
	}

	table, err := s.s.GetAll()
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось обработать запрос", http.StatusBadGateway)
		return
	}

	if err := json.NewEncoder(w).Encode(table); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Ошибка при составлении ответа", http.StatusBadGateway)
		return

	}
}

func (s *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		http.Error(w, "Пустое тело запроса", http.StatusBadGateway)
		return
	}

	name := r.FormValue("name")
	text := r.FormValue("text")
	if name == "" {
		http.Error(w, "Пустое имя", http.StatusBadGateway)
		return
	}
	if text == "" {
		http.Error(w, "Пустой текст", http.StatusBadGateway)
		return
	}

	err, id := s.s.CreateNew(name, text)
	if err != nil {
		http.Error(w, "Не удалось обработать запрос", http.StatusBadGateway)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Ошибка при составлении ответа", http.StatusBadGateway)
		return
	}
}

func (s *Handler) AddPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		http.Error(w, "Пустое тело запроса", http.StatusBadGateway)
		return
	}

	name := r.FormValue("name")
	text := r.FormValue("text")
	if name == "" {
		http.Error(w, "Пустое имя", http.StatusBadGateway)
		return
	}
	if text == "" {
		http.Error(w, "Пустой текст", http.StatusBadGateway)
		return
	}

	err, _ := s.s.CreateNew(name, text)
	if err != nil {
		http.Error(w, "Не удалось обработать запрос", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		http.Error(w, "Пустое тело запроса", http.StatusBadGateway)
		return
	}

	index := r.FormValue("index")
	if index == "" {
		http.Error(w, "Пустой индекс", http.StatusBadGateway)
		return
	}

	intIndex, err := strconv.Atoi(index)
	if err != nil {
		http.Error(w, "Пустое тело запроса", http.StatusBadGateway)
		return
	}

	if err := s.s.DeleteById(intIndex); err != nil {
		http.Error(w, "Пустое тело запроса", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Handler) Index(w http.ResponseWriter, r *http.Request) {
	table, err := s.s.GetAll()
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось прочитать", http.StatusBadGateway)
		return
	}

	tmpl.Execute(w, table)
}

func (s *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := s.s.HealthSer()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
	return
}

var (
	mu   sync.Mutex
	tmpl = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Заметки</title>
</head>
<body>
  <h3>Заметки:</h3>
  <ul>
  {{range .}}
    <li>
      {{.Name}} — {{.Text}}
      <form action="/delete" method="post" style="display:inline">
        <input type="hidden" name="index" value="{{.Id}}">
        <button>Удалить</button>
      </form>
    </li>
  {{end}}
  </ul>

  <form action="/add" method="post">
    <input name="name" autocomplete="off" placeholder="Имя">
    <input name="text" autocomplete="off" placeholder="Заметка">
    <button>Добавить</button>
  </form>

  <h3>Регистрация</h3>
  <form action="/register" method="post">
    <input name="name" autocomplete="off" placeholder="Имя">
    <input name="password" type="password" autocomplete="off" placeholder="Пароль">
    <button>Зарегистрироваться</button>
  </form>

  <h3>Вход</h3>
  <form action="/login" method="post">
    <input name="name" autocomplete="off" placeholder="Имя">
    <input name="password" type="password" autocomplete="off" placeholder="Пароль">
    <button>Войти</button>
  </form>
</body>
</html>
`))
)
