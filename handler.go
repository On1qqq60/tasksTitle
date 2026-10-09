package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	filePath = "D:/hz/taskTitle/notes.txt"
)

func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось открыть файл", http.StatusBadGateway)
		return
	}
	defer file.Close()

	mu.Lock()
	defer mu.Unlock()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось прочитать файл", http.StatusBadGateway)
		return
	}

	var items []string
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			items = append(items, line)
		}
	}

	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, "Не удалось создать Json", http.StatusInternalServerError)
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Пустой текст", http.StatusBadGateway)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось открыть файл", http.StatusBadGateway)
		return
	}
	defer file.Close()

	err = saveNuties(text, file, w)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Ошибка записи в файл", http.StatusBadGateway)
		return
	}

	i, err := file.Seek(0, 0)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		http.Error(w, "Не удалось переместить курсор в файле", http.StatusBadGateway)
		return
	}
	if i != 0 {
		print("Похуй мне")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		http.Error(w, "Не удалось прочитать в файл", http.StatusBadGateway)
		return
	}

	var items []string
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			items = append(items, line)
		}
	}

	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		http.Error(w, "Не удалось создать Json", http.StatusInternalServerError)
		return
	}
}

func AddPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Пустой текст", http.StatusBadGateway)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось открыть файл", http.StatusBadGateway)
		return
	}
	defer file.Close()

	err = saveNuties(text, file, w)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Ошибка записи в файл", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	i := r.FormValue("index")
	index, err := strconv.Atoi(i)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		http.Error(w, "Не удалось преоброзовать строку в int", http.StatusBadGateway)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		http.Error(w, "Не удалось открыть файл", http.StatusBadGateway)
		return
	}
	defer file.Close()

	mu.Lock()
	defer mu.Unlock()
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		http.Error(w, "Не удалось прочитать файл", http.StatusBadGateway)
		return
	}

	var items []string
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			items = append(items, line)
		}
	}

	for i, _ := range items {
		if i == index {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	content = []byte(strings.Join(items, "\n"))
	content = append(content, '\n')

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		http.Error(w, "Не удалось записать в файл", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func saveNuties(text string, file *os.File, w http.ResponseWriter) error {
	byteText := []byte(text)
	byteText = append(byteText, '\n')

	mu.Lock()
	defer mu.Unlock()
	n, err := file.Write(byteText)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		http.Error(w, "Не удалось записать в файл", http.StatusBadGateway)
		return err
	}
	if n < len([]byte(byteText)) {
		http.Error(w, "Не всё записалось в файл", http.StatusBadGateway)
		return err
	}

	return nil
}
