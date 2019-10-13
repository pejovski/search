package domain

import (
	"github.com/pejovski/search/model"
	"github.com/pejovski/search/scope"
	"net/http"
)

type HttpHandler interface {
	Product() http.HandlerFunc
	Products() http.HandlerFunc
	CreateProduct() http.HandlerFunc
	DeleteProduct() http.HandlerFunc
}

type SearchController interface {
	GetProduct(id string) (*model.Product, error)
	GetProducts(s *scope.Scope) (ps []*model.Product, total int, err error)
	CreateProduct(p *model.Product) (id string, err error)
	DeleteProduct(id string) error
}

type SearchRepository interface {
	Product(id string) (*model.Product, error)
	Products(s *scope.Scope) (ps []*model.Product, total int, err error)
	Create(p *model.Product) (id string, err error)
	Delete(id string) error
}
