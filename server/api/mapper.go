package api

import (
	"github.com/pejovski/search/model"
)

type Mapper interface {
	mapDomainProductToProduct(dp *model.Product) *Product
	mapDomainProductsToProducts(dps []*model.Product) []*Product
}

type mapper struct {
}

func newMapper() Mapper {
	return mapper{}
}

func (m mapper) mapDomainProductToProduct(dp *model.Product) *Product {
	return &Product{
		Id:    dp.Id,
		Title: dp.Title,
		Brand: dp.Brand,
		Price: dp.Price,
		Stock: dp.Stock,
	}
}

func (m mapper) mapDomainProductsToProducts(dps []*model.Product) []*Product {
	var ps []*Product
	for _, dp := range dps {
		ps = append(ps, m.mapDomainProductToProduct(dp))
	}
	return ps
}
