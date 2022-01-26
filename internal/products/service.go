package products

import "github.com/rodrisouzag/storage-implementation/internal/models"

type Service interface {
	GetByName(name string) models.Product
	Store(name string, group string, count int, price float64) (models.Product, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) GetByName(name string) models.Product {
	return s.repository.GetByName(name)
}

func (s *service) Store(name string, category string, count int, price float64) (models.Product, error) {
	p := models.Product{ID: 0, Name: name, Category: category, Count: count, Price: price}
	return s.repository.Store(p)
}
