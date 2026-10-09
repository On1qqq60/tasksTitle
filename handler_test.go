package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestEmptyTextAndDuplicateDelete(t *testing.T) {
	original, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.WriteFile(filePath, original, 0644); err != nil {
			t.Error(err)
		}
	})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("GET /api/notes", List)
	mux.HandleFunc("POST /api/notes", Add)
	mux.HandleFunc("POST /delete", Delete)
	mux.HandleFunc("POST /add", AddPage)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	const before = "хлеб\n"
	if err := os.WriteFile(filePath, []byte(before), 0644); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/api/notes", "/add"} {
		postForm(t, srv.URL+path, url.Values{"text": {""}})
		got, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != before {
			t.Fatalf("пустой текст на %s дописал файл: %q", path, got)
		}
	}

	if err := os.WriteFile(filePath, []byte("молоко\nмолоко\n"), 0644); err != nil {
		t.Fatal(err)
	}
	postForm(t, srv.URL+"/delete", url.Values{"index": {"0"}})

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	lines := nonEmptyLines(string(got))
	if len(lines) != 1 || lines[0] != "молоко" {
		t.Fatalf("удаление одной из двух одинаковых строк оставило %q", got)
	}
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

func nonEmptyLines(content string) []string {
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
