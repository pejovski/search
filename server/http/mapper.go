package http

import (
	"github.com/pejovski/search/entity"
)

func mapDomainProductToProduct(dp *entity.Product) *Product {
	return &Product{
		Id:    dp.Id,
		Title: dp.Title,
		Brand: dp.Brand,
		Price: dp.Price,
		Stock: dp.Stock,
	}
}

func mapDomainProductsToProducts(dps []*entity.Product) []*Product {
	ps := []*Product{}
	for _, dp := range dps {
		ps = append(ps, mapDomainProductToProduct(dp))
	}
	return ps
}
