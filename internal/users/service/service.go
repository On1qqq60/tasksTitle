package service

import (
	"errors"
	"taskTitle/internal/users/repository"

	"github.com/google/uuid"
)

type Repository struct {
	repo repository.Repository
}

type Service interface {
	Register(name string, password string) (uuid.UUID, error)
	Login(name string, password string) (bool, error)
}

func NewService(r repository.Repository) *Repository {
	return &Repository{repo: r}
}

func (r *Repository) Register(name string, password string) (uuid.UUID, error) {
	if name == "" {
		return uuid.Nil, errors.New("Пустое имя")
	}
	if password == "" {
		return uuid.Nil, errors.New("Пустой пароль")
	}

	id := uuid.New()

	err := r.repo.Add(id, name, password)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *Repository) Login(name string, password string) (bool, error) {
	if name == "" {
		return false, errors.New("Пустое имя")
	}
	if password == "" {
		return false, errors.New("Пустой пароль")
	}

	cor, err := r.repo.IsCorrect(name, password)
	return cor, err
}
