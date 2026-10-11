package service

import (
	"errors"
	"taskTitle/internal/tasks/models"
	"taskTitle/internal/tasks/repository"
)

type Service struct {
	repo repository.Repository
}

type ServiceInter interface {
	GetAll() ([]models.Table, error)
	CreateNew(name, text string) (error, int)
	DeleteById(id int) error
	HealthSer() error
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (r *Service) GetAll() ([]models.Table, error) {
	table, err := r.repo.ReadDB()
	if err != nil {
		return nil, err
	}

	return table, err
}

func (r *Service) CreateNew(name, text string) (error, int) {
	if name == "" {
		return errors.New("Пустое имя"), 0
	}
	if text == "" {
		return errors.New("Текст пустой"), 0
	}

	id, err := r.repo.CreateDB(name, text)
	if err != nil {
		return err, 0
	}

	idInt := int(id)

	return nil, idInt
}

func (r *Service) DeleteById(id int) error {
	if id == 0 {
		return errors.New("Не корректный id")
	}

	err := r.repo.DeleteDB(id)

	return err
}

func (r *Service) HealthSer() error {
	if err := r.repo.HealthDB(); err != nil {
		return err
	}

	return nil
}
