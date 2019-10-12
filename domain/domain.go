package domain

import (
	"github.com/pejovski/search/entity"
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
	GetProduct(id string) (*entity.Product, error)
	GetProducts(s *scope.Scope) (ps []*entity.Product, total int, err error)
	CreateProduct(p *entity.Product) (id string, err error)
	DeleteProduct(id string) error
}

type SearchRepository interface {
	Product(id string) (*entity.Product, error)
	Products(s *scope.Scope) (ps []*entity.Product, total int, err error)
	Create(p *entity.Product) (id string, err error)
	Delete(id string) error
}
