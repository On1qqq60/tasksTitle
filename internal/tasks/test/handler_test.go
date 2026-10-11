package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"taskTitle/internal/tasks/handler"
	"taskTitle/internal/tasks/models"
	"taskTitle/internal/tasks/repository"
	"taskTitle/internal/tasks/service"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func TestEmptyTextAndDuplicateDelete(t *testing.T) {
	srv := newTestServer(t)

	for _, path := range []string{"/api/notes", "/add"} {
		postForm(t, srv.URL+path, url.Values{
			"name": {"тест"},
			"text": {""},
		})
		if got := listNotes(t, srv.URL); len(got) != 0 {
			t.Fatalf("пустой текст на %s добавил записи: %+v", path, got)
		}
	}

	id1 := addNote(t, srv.URL, "один", "молоко")
	id2 := addNote(t, srv.URL, "один", "молоко")
	if id1 == id2 {
		t.Fatalf("ожидали разные id, получили %d", id1)
	}

	postForm(t, srv.URL+"/delete", url.Values{"index": {strconv.Itoa(id1)}})

	got := listNotes(t, srv.URL)
	if len(got) != 1 {
		t.Fatalf("удаление одной из двух одинаковых оставило %+v", got)
	}
	if got[0].Id != id2 || got[0].Name != "один" || got[0].Text != "молоко" {
		t.Fatalf("осталась не та запись: %+v", got[0])
	}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	ensureTestDatabase(t)
	t.Setenv("DATABASE_URL", postgresDSN("notes_test"))

	db, err := repository.OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	if err := repository.InitSchemaTask(db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM notes"); err != nil {
		t.Fatal(err)
	}

	h := handler.NewHandler(service.NewService(repository.NewRepository(db)))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Index)
	mux.HandleFunc("GET /api/notes", h.ListHandler)
	mux.HandleFunc("POST /api/notes", h.AddHandler)
	mux.HandleFunc("POST /delete", h.DeleteHandler)
	mux.HandleFunc("POST /add", h.AddPageHandler)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func ensureTestDatabase(t *testing.T) {
	t.Helper()

	cfg, err := pgx.ParseConfig(postgresDSN("notes"))
	if err != nil {
		t.Fatal(err)
	}
	db := stdlib.OpenDB(*cfg)
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	var exists bool
	err = db.QueryRow(`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'notes_test')`).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		return
	}
	if _, err := db.Exec(`CREATE DATABASE notes_test`); err != nil {
		t.Fatal(err)
	}
}

func postgresDSN(dbName string) string {
	port := os.Getenv("PGPORT")
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("postgres://notes:notes@127.0.0.1:%s/%s?sslmode=disable", port, dbName)
}

func addNote(t *testing.T, baseURL, name, text string) int {
	t.Helper()
	resp, err := http.PostForm(baseURL+"/api/notes", url.Values{
		"name": {name},
		"text": {text},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("добавление: статус %d, тело %s", resp.StatusCode, body)
	}
	var id int
	if err := json.Unmarshal(body, &id); err != nil {
		t.Fatalf("id: %v, тело %s", err, body)
	}
	return id
}

func listNotes(t *testing.T, baseURL string) []models.Table {
	t.Helper()
	resp, err := http.Get(baseURL + "/api/notes")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("список: статус %d, тело %s", resp.StatusCode, body)
	}
	var notes []models.Table
	if err := json.Unmarshal(body, &notes); err != nil {
		t.Fatalf("json: %v, тело %s", err, body)
	}
	return notes
}

func postForm(t *testing.T, rawURL string, form url.Values) {
	t.Helper()
	resp, err := http.PostForm(rawURL, form)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}
