package controller

import (
	"github.com/pejovski/search/domain"
	"github.com/pejovski/search/entity"
	"github.com/pejovski/search/scope"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

type Search struct {
	repository domain.SearchRepository
}

func NewSearch(r domain.SearchRepository) Search {
	return Search{repository: r}
}

func (c Search) GetProduct(id string) (*entity.Product, error) {
	p, err := c.repository.Product(id)
	if err != nil {
		logrus.Errorf("Failed to get product %s; Error: %s", id, err)
		return nil, err
	}

	return p, nil
}

func (c Search) GetProducts(s *scope.Scope) ([]*entity.Product, int, error) {
	ps, total, err := c.repository.Products(s)
	if err != nil {
		logrus.Errorf("Failed to get products for scope %v; Error: %s", s, err)
		return nil, 0, err
	}

	return ps, total, nil
}

func (c Search) CreateProduct(p *entity.Product) (string, error) {
	p.Id = ksuid.New().String()
	id, err := c.repository.Create(p)
	if err != nil {
		logrus.Errorf("Failed to create product; Error: %s", err)
		return "", err
	}

	return id, nil
}

func (c Search) DeleteProduct(id string) error {
	err := c.repository.Delete(id)
	if err != nil {
		logrus.Errorf("Failed to delete product %s; Error: %s", id, err)
		return err
	}

	return nil
}
