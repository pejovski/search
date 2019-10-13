package http

import (
	"github.com/pejovski/search/model"
)

func mapDomainProductToProduct(dp *model.Product) *Product {
	return &Product{
		Id:    dp.Id,
		Title: dp.Title,
		Brand: dp.Brand,
		Price: dp.Price,
		Stock: dp.Stock,
	}
}

func mapDomainProductsToProducts(dps []*model.Product) []*Product {
	ps := []*Product{}
	for _, dp := range dps {
		ps = append(ps, mapDomainProductToProduct(dp))
	}
	return ps
}
