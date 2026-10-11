package handler

import (
	"encoding/json"
	"net/http"
	"taskTitle/internal/users/service"
)

type Service struct {
	srv service.Service
}

type Handler interface {
	RegistrHandler()
	LoginHandler()
}

func NewHandler(s service.Service) *Service {
	return &Service{srv: s}
}

func (s *Service) RegistrHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	password := r.FormValue("password")

	uuid, err := s.srv.Register(name, password)
	if err != nil {
		http.Error(w, "Ошибка", http.StatusBadRequest)
		return
	}

	idStr := uuid.String()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"uuid": idStr,
	})
}

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.FormValue("name")
	password := r.FormValue("password")

	cor, err := s.srv.Login(name, password)
	if err != nil {
		http.Error(w, "Ошибка", http.StatusBadRequest)
		return
	}

	if !cor {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "!ok",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"name":   name,
	})
}
