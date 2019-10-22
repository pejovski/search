package controller

import (
	"github.com/pejovski/search/model"
	"github.com/pejovski/search/pkg/scope"
	"github.com/pejovski/search/repository"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	GetProduct(id string) (*model.Product, error)
	GetProducts(s *scope.Scope) (ps []*model.Product, total int, err error)
	CreateProduct(p *model.Product) (id string, err error)
	DeleteProduct(id string) error
}

type controller struct {
	repository repository.Repository
}

func New(r repository.Repository) Controller {
	return controller{repository: r}
}

func (c controller) GetProduct(id string) (*model.Product, error) {
	p, err := c.repository.Product(id)
	if err != nil {
		logrus.Errorf("Failed to get product %s; Error: %s", id, err)
		return nil, err
	}

	return p, nil
}

func (c controller) GetProducts(s *scope.Scope) ([]*model.Product, int, error) {
	ps, total, err := c.repository.Products(s)
	if err != nil {
		logrus.Errorf("Failed to get products for scope %v; Error: %s", s, err)
		return nil, 0, err
	}

	return ps, total, nil
}

func (c controller) CreateProduct(p *model.Product) (string, error) {
	p.ID = ksuid.New().String()
	id, err := c.repository.Create(p)
	if err != nil {
		logrus.Errorf("Failed to create product; Error: %s", err)
		return "", err
	}

	return id, nil
}

func (c controller) DeleteProduct(id string) error {
	err := c.repository.Delete(id)
	if err != nil {
		logrus.Errorf("Failed to delete product %s; Error: %s", id, err)
		return err
	}

	return nil
}
