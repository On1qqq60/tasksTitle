package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	mu   sync.Mutex
	tmpl = template.Must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Товары</title>
</head>
<body>
  <h3>Товары:</h3>
  <ul>
  {{range $i, $item := .}}
    <li>
      {{$item}}
      <form action="/delete" method="post" style="display:inline">
        <input type="hidden" name="index" value="{{$i}}">
        <button>Удалить</button>
      </form>
    </li>
  {{end}}
  </ul>

  <form action="/add" method="post">
    <input name="text" autocomplete="off">
    <button>Добавить</button>
  </form>
</body>
</html>
`))
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /api/notes", List)
	mux.HandleFunc("POST /api/notes", Add)
	mux.HandleFunc("POST /delete", Delete)
	mux.HandleFunc("POST /add", AddPage)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка: %s", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	data, err := os.ReadFile("D:/hz/taskTitle/notes.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			items = append(items, line)
		}
	}

	tmpl.Execute(w, items)
}
