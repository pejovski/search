package repository

import (
	"github.com/pejovski/search/model"
	"github.com/pejovski/search/pkg/scope"
)

type Repository interface {
	Product(id string) (*model.Product, error)
	Products(s *scope.Scope) (ps []*model.Product, total int, err error)
	Create(p *model.Product) (id string, err error)
	Delete(id string) error
}
